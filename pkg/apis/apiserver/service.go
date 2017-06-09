package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"apiserver/pkg/api/apiserver"
	"apiserver/pkg/configz"
	"apiserver/pkg/resource"
	k8sclient "apiserver/pkg/resource/common"
	"apiserver/pkg/resource/configMap"
	"apiserver/pkg/resource/deployment"
	"apiserver/pkg/resource/service"
	"apiserver/pkg/resource/sync"
	r "apiserver/pkg/router"
	"apiserver/pkg/util/parseUtil"

	// "apiserver/pkg/util/jsonx"
	"apiserver/pkg/util/log"

	"github.com/gorilla/mux"
	res "k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/pkg/api/v1"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

func GetServices(request *http.Request) (string, interface{}) {
	pageCnt := request.FormValue("pageCnt")
	pageNum := request.FormValue("pageNum")
	cnt, _ := strconv.Atoi(pageCnt)
	num, _ := strconv.Atoi(pageNum)
	serviceName := request.FormValue("name")
	appId, _ := strconv.ParseUint(request.FormValue("appId"), 10, 64)
	svcs, total := apiserver.QueryServices(serviceName, cnt, num, uint(appId))
	return r.StatusOK, map[string]interface{}{"services": svcs, "total": total}
}

func GetService(request *http.Request) (string, interface{}) {
	id, _ := strconv.ParseUint(mux.Vars(request)["id"], 10, 64)
	return r.StatusOK, map[string]interface{}{"service": apiserver.QueryServiceById(uint(id))}
}

func CreateService(request *http.Request) (string, interface{}) {
	reqservice, err := validateService(request)
	if err != nil {
		return r.StatusBadRequest, err
	}

	app := apiserver.GetAppOnly(reqservice.AppId)
	app.Items = []*apiserver.Service{reqservice}
	k8ssvc := service.NewService(app)
	if k8sclient.ExsitResource(k8ssvc) {
		return r.StatusForbidden, "the service exist"
	}

	svc, err := k8sclient.CreateService(k8ssvc)
	if err != nil {
		return r.StatusInternalServerError, err
	}

	if app.Items[0].Config.ConfigMap != nil {
		cfgMap := configMap.NewConfigMap(app)
		if err := k8sclient.CreateResource(cfgMap); err != nil {
			return r.StatusInternalServerError, err
		}
	}

	k8sDeploy := deployment.NewDeployment(app)
	if k8sclient.ExsitResource(k8sDeploy) {
		return r.StatusForbidden, "the deployment exist"
	}
	if err = k8sclient.CreateResource(k8sDeploy); err != nil {
		k8sclient.DeleteResource(svc)
		return r.StatusInternalServerError, err
	}
	external := fmt.Sprintf("http://%s:%v", configz.GetString("apiserver", "clusterNodes", "127.0.0.1"), svc.Spec.Ports[0].NodePort)
	app.External = external
	app.Items[0].External = external

	time.Sleep(5 * time.Second)
	podList, err := k8sclient.GetDeploymentPods(svc.Name, app.UserName)
	if err != nil {
		return r.StatusInternalServerError, err
	}

	var cs []*apiserver.Container
	for _, pod := range podList {
		c := &apiserver.Container{Name: pod.ObjectMeta.Name, Image: app.Items[0].Image, Internal: pod.Status.PodIP}
		if pod.Status.Phase == "running" {
			c.Status = resource.AppRunning
		}

		cs = append(cs, c)
	}
	app.Items[0].Items = cs
	apiserver.InsertService(app.Items[0])
	return r.StatusCreated, "ok"
}

func DeleteService(request *http.Request) (string, interface{}) {
	id, _ := strconv.ParseUint(mux.Vars(request)["id"], 10, 64)
	svc := apiserver.QueryServiceById(uint(id))
	// log.Debug(jsonx.ToJson(svc))
	apiserver.DeleteService(svc)
	return r.StatusOK, "ok"
}

func UpdateServiceConfig(request *http.Request) (string, interface{}) {
	id, _ := strconv.ParseUint(mux.Vars(request)["id"], 10, 64)
	svc := apiserver.QueryServiceById(uint(id))
	namespace := mux.Vars(request)["namespace"]
	verb := mux.Vars(request)["verb"] //verb the action of app , scale or expansion or roll

	deploy, exsit := sync.ListDeployment[namespace][svc.Name]
	if !exsit {
		return r.StatusNotFound, "service named " + svc.Name + ` does't exist`
	}
	if verb == "scale" {
		scaleOption, err := validateScaleOption(request)
		log.Debugf("%#v", scaleOption)
		if err != nil {
			return r.StatusBadRequest, err
		}
		deploy.Spec.Replicas = parseUtil.IntToInt32Pointer(scaleOption.ServiceInstanceCnt)
		svc.InstanceCount = scaleOption.ServiceInstanceCnt
		if err := k8sclient.UpdateResouce(&deploy); err != nil {
			return r.StatusInternalServerError, "rolling updte application named " + svc.Name + ` failed`
		}
		for _, container := range svc.Items {
			apiserver.DeleteContainer(container)
		}
		ChangeServiceStatus(svc, namespace)
	}

	if verb == "expansion" {
		expansionOption, err := validateExpansionOption(request)
		if err != nil {
			return r.StatusBadRequest, err
		}

		deploy.Spec.Template.Spec.Containers[0].Resources = v1.ResourceRequirements{
			Limits: v1.ResourceList{
				v1.ResourceCPU:    res.MustParse(expansionOption.Cpu),    //TODO 根据前端传入的值做资源限制
				v1.ResourceMemory: res.MustParse(expansionOption.Memory), //TODO 根据前端传入的值做资源限制
			},
			Requests: v1.ResourceList{
				v1.ResourceCPU:    res.MustParse(expansionOption.Cpu),
				v1.ResourceMemory: res.MustParse(expansionOption.Memory),
			},
		}
		svc.Config.BaseConfig.Cpu = expansionOption.Cpu
		svc.Config.BaseConfig.Memory = expansionOption.Memory

		svc.Items[0].Config.BaseConfig.Cpu = expansionOption.Cpu
		svc.Items[0].Config.BaseConfig.Memory = expansionOption.Memory
		if err := k8sclient.UpdateResouce(&deploy); err != nil {
			return r.StatusInternalServerError, "rolling updte application named " + svc.Name + ` failed`
		}
	}

	if verb == "roll" {
		rollOption, err := validateRollOption(request)
		if err != nil {
			return r.StatusBadRequest, err
		}
		deploy.Spec.Template.Spec.Containers[0].Image = rollOption.Image
		deploy.Spec.Strategy = extensions.DeploymentStrategy{
			Type: extensions.RollingUpdateDeploymentStrategyType,
			RollingUpdate: &extensions.RollingUpdateDeployment{
				MaxUnavailable: &intstr.IntOrString{Type: intstr.String, StrVal: "20%"},
				MaxSurge:       &intstr.IntOrString{Type: intstr.String, StrVal: "120%"},
			},
		}

		svc.Image = rollOption.Image
		svc.Items[0].Image = rollOption.Image
		if err := k8sclient.UpdateResouce(&deploy); err != nil {
			return r.StatusInternalServerError, "rolling updte application named " + svc.Name + ` failed`
		}
	}

	return r.StatusOK, "ok"
}

func StopOrStartOrRedployService(request *http.Request) (string, interface{}) {
	id, _ := strconv.ParseUint(mux.Vars(request)["id"], 10, 64)
	svc := apiserver.QueryServiceById(uint(id))
	namespace := mux.Vars(request)["namespace"]
	verb := mux.Vars(request)["verb"] //verb the action of app , start or stop or redeploy
	deploy, exsit := sync.ListDeployment[namespace][svc.Name]
	if !exsit {
		return r.StatusNotFound, "service named " + svc.Name + ` does't exist`
	}
	if verb == "stop" {
		deploy.Spec.Replicas = parseUtil.IntToInt32Pointer(0)
		if err := k8sclient.UpdateResouce(&deploy); err != nil {
			return r.StatusInternalServerError, err
		}

		svc.Status = resource.AppStop
		for _, container := range svc.Items {
			apiserver.DeleteContainer(container)
		}
	}
	if verb == "start" {
		deploy.Spec.Replicas = parseUtil.IntToInt32Pointer(svc.InstanceCount)

		if err := k8sclient.UpdateResouce(&deploy); err != nil {
			return r.StatusInternalServerError, err
		}

		svc.Status = resource.AppRunning
		if err := ChangeServiceStatus(svc, namespace); err != nil {
			return r.StatusInternalServerError, err
		}

	}
	if verb == "redeploy" {
		pods, err := k8sclient.GetPods(namespace, svc.Name)
		if err != nil {
			return r.StatusInternalServerError, err
		}
		for _, pod := range pods {
			k8sclient.DeleteResource(&pod)
		}

		for _, container := range svc.Items {
			apiserver.DeleteContainer(container)
		}

		svc.Status = resource.AppRunning
		if err := ChangeServiceStatus(svc, namespace); err != nil {
			return r.StatusInternalServerError, err
		}
	}
	return r.StatusOK, "ok"
}

func validateService(request *http.Request) (*apiserver.Service, error) {
	svc := &apiserver.Service{}
	if err := json.NewDecoder(request.Body).Decode(svc); err != nil {
		return nil, err
	}
	return svc, nil
}

func validateScaleOption(request *http.Request) (*apiserver.ScaleOption, error) {
	option := &apiserver.ScaleOption{}
	if err := json.NewDecoder(request.Body).Decode(option); err != nil {
		return nil, err
	}
	return option, nil
}

func validateExpansionOption(request *http.Request) (*apiserver.ExpansionOption, error) {
	option := &apiserver.ExpansionOption{}
	if err := json.NewDecoder(request.Body).Decode(option); err != nil {
		return nil, err
	}
	return option, nil
}

func validateRollOption(request *http.Request) (*apiserver.RollOption, error) {
	option := &apiserver.RollOption{}
	if err := json.NewDecoder(request.Body).Decode(option); err != nil {
		return nil, err
	}
	return option, nil
}

func ChangeServiceStatus(svc *apiserver.Service, namespace string) error {
	time.Sleep(5 * time.Second)
	podList, err := k8sclient.GetPods(namespace, svc.Name)
	if err != nil {
		return err
	}

	var containers []*apiserver.Container
	for _, pod := range podList {
		container := &apiserver.Container{Name: pod.ObjectMeta.Name, Image: svc.Image, Internal: pod.Status.PodIP, ServiceId: svc.ID}
		// if pod.Status.ContainerStatuses[0].Ready == true {
		// 	container.Status = resource.AppRunning
		// }
		containers = append(containers, container)

	}
	svc.Items = containers
	apiserver.UpdateService(svc)
	return nil
}

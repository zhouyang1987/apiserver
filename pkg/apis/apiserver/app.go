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

	"github.com/gorilla/mux"
	res "k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/pkg/api/v1"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

func GetApps(request *http.Request) (string, interface{}) {
	namespace := mux.Vars(request)["namespace"]
	pageCnt, _ := strconv.Atoi(request.FormValue("pageCnt"))
	pageNum, _ := strconv.Atoi(request.FormValue("pageNum"))
	appName := request.FormValue("name")
	apps, total := apiserver.QueryApps(namespace, appName, pageCnt, pageNum)
	return r.StatusOK, map[string]interface{}{"apps": apps, "total": total}
}

func GetApp(request *http.Request) (string, interface{}) {
	id, _ := strconv.ParseUint(mux.Vars(request)["id"], 10, 64)
	return r.StatusOK, apiserver.QueryAppById(uint(id))

}

func CreateApp(request *http.Request) (string, interface{}) {
	app, err := validateApp(request)
	if err != nil {
		return r.StatusBadRequest, err
	}

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
	apiserver.InsertApp(app)
	return r.StatusCreated, "ok"
}

func DeleteApp(request *http.Request) (string, interface{}) {
	id, _ := strconv.ParseUint(mux.Vars(request)["id"], 10, 64)
	app := apiserver.QueryAppById(uint(id))
	appName := app.Items[0].Name
	namespace := app.UserName
	rc, exsit := sync.ListDeployment[namespace][appName]
	if !exsit {
		return r.StatusNotFound, "application named " + appName + ` does't exist`
	}

	if err := k8sclient.DeleteResource(&rc); err != nil {
		return r.StatusInternalServerError, "delete application err: " + err.Error()
	}
	svc := sync.ListService[namespace][appName]
	if &svc == nil {
		return r.StatusNotFound, "application named " + appName + `does't exist`
	}
	if err := k8sclient.DeleteResource(&svc); err != nil {
		return r.StatusInternalServerError, "delete application err: " + err.Error()
	}

	apiserver.DeleteApp(app)
	return r.StatusNoContent, "ok"
}

func UpdateAppConfig(request *http.Request) (string, interface{}) {
	app, err := validateApp(request)
	if err != nil {
		return r.StatusBadRequest, err
	}
	// id, _ := strconv.ParseUint(mux.Vars(request)["id"], 10, 64)
	// app := apiserver.QueryById(uint(id))
	appName := app.Items[0].Name
	namespace := mux.Vars(request)["namespace"]
	verb := mux.Vars(request)["verb"] //verb the action of app , scale or expansion or roll

	deploy, exsit := sync.ListDeployment[namespace][appName]
	if !exsit {
		return r.StatusNotFound, "application named " + appName + ` does't exist`
	}
	if verb == "scale" {
		deploy.Spec.Replicas = parseUtil.IntToInt32Pointer(app.Items[0].InstanceCount)
	}

	if verb == "expansion" {
		deploy.Spec.Template.Spec.Containers[0].Resources = v1.ResourceRequirements{
			Limits: v1.ResourceList{
				v1.ResourceCPU:    res.MustParse(app.Items[0].Config.BaseConfig.Cpu),    //TODO 根据前端传入的值做资源限制
				v1.ResourceMemory: res.MustParse(app.Items[0].Config.BaseConfig.Memory), //TODO 根据前端传入的值做资源限制
			},
			Requests: v1.ResourceList{
				v1.ResourceCPU:    res.MustParse(app.Items[0].Config.BaseConfig.Cpu),
				v1.ResourceMemory: res.MustParse(app.Items[0].Config.BaseConfig.Memory),
			},
		}
	}

	if verb == "roll" {
		deploy.Spec.Template.Spec.Containers[0].Image = app.Items[0].Image
		deploy.Spec.Strategy = extensions.DeploymentStrategy{
			Type: extensions.RollingUpdateDeploymentStrategyType,
			RollingUpdate: &extensions.RollingUpdateDeployment{
				MaxUnavailable: &intstr.IntOrString{Type: intstr.String, StrVal: "20%"},
				MaxSurge:       &intstr.IntOrString{Type: intstr.String, StrVal: "120%"},
			},
		}
	}

	if err := k8sclient.UpdateResouce(&deploy); err != nil {
		return r.StatusInternalServerError, "rolling updte application named " + appName + ` failed`
	}

	return r.StatusCreated, "ok"
}

func StopOrStartOrRedeployApp(request *http.Request) (string, interface{}) {
	id, _ := strconv.ParseUint(mux.Vars(request)["id"], 10, 64)
	app := apiserver.QueryAppById(uint(id))
	for _, svc := range app.Items {
		appName := svc.Name
		namespace := mux.Vars(request)["namespace"]
		verb := mux.Vars(request)["verb"] //verb the action of app , start or stop
		deploy, exsit := sync.ListDeployment[namespace][appName]
		if !exsit {
			return r.StatusNotFound, "service named " + appName + ` does't exist`
		}
		if verb == "stop" {
			deploy.Spec.Replicas = parseUtil.IntToInt32Pointer(0)
			if err := k8sclient.UpdateResouce(&deploy); err != nil {
				return r.StatusInternalServerError, err
			}

			app.AppStatus = resource.AppStop
			svc.Status = resource.AppStop
			apiserver.UpdateApp(app)
			for _, container := range svc.Items {
				apiserver.DeleteContainer(container)
			}
		}
		if verb == "start" {
			deploy.Spec.Replicas = parseUtil.IntToInt32Pointer(app.Items[0].InstanceCount)

			if err := k8sclient.UpdateResouce(&deploy); err != nil {
				return r.StatusInternalServerError, err
			}

			app.AppStatus = resource.AppRunning
			svc.Status = resource.AppRunning
			if err := ChangeAppStatus(svc, app); err != nil {
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

			app.AppStatus = resource.AppRunning
			svc.Status = resource.AppRunning
			if err := ChangeAppStatus(svc, app); err != nil {
				return r.StatusInternalServerError, err
			}
		}
	}
	return r.StatusCreated, "ok"
}

func validateApp(request *http.Request) (*apiserver.App, error) {
	app := &apiserver.App{}
	if err := json.NewDecoder(request.Body).Decode(app); err != nil {
		return nil, err
	}
	return app, nil
}

func validateConfig(request *http.Request) (*apiserver.ServiceConfig, error) {
	cfg := &apiserver.ServiceConfig{}
	if err := json.NewDecoder(request.Body).Decode(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func ChangeAppStatus(svc *apiserver.Service, app *apiserver.App) error {
	time.Sleep(5 * time.Second)
	podList, err := k8sclient.GetPods(app.UserName, svc.Name)
	if err != nil {
		return err
	}

	var containers []*apiserver.Container
	for _, pod := range podList {
		container := &apiserver.Container{Name: pod.ObjectMeta.Name, Image: app.Items[0].Image, Internal: pod.Status.PodIP, ServiceId: svc.ID}
		// if pod.Status.ContainerStatuses[0].Ready == true {
		// 	container.Status = resource.AppRunning
		// }
		containers = append(containers, container)

	}
	svc.Items = containers
	apiserver.UpdateApp(app)
	return nil
}

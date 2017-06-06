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

func Register(router *mux.Router) {
	r.RegisterHttpHandler(router, "/{namespace}/apps", "GET", GetApps)
	r.RegisterHttpHandler(router, "/{namespace}/apps/{id}", "GET", GetApp)
	r.RegisterHttpHandler(router, "/{namespace}/apps", "POST", CreateApp)
	r.RegisterHttpHandler(router, "/{namespace}/apps/{id}", "DELETE", DeleteApp)
	r.RegisterHttpHandler(router, "/{namespace}/apps/{id}/{verb}", "PUT", UpdateAppConfig)
	r.RegisterHttpHandler(router, "/{namespace}/apps/{id}/{verb}", "PATCH", StopOrStartApp)
	r.RegisterHttpHandler(router, "/{namespace}/apps", "OPTIONS", Option)

}

func Option(request *http.Request) (string, interface{}) {
	return r.StatusOK, "ok"
}

func GetApps(request *http.Request) (string, interface{}) {
	namespace := mux.Vars(request)["namespace"]
	list := apiserver.QueryAll(namespace)
	return r.StatusOK, list
}

func GetApp(request *http.Request) (string, interface{}) {
	id, _ := strconv.ParseUint(mux.Vars(request)["id"], 10, 64)
	return r.StatusOK, apiserver.QueryById(uint(id))
}

func StopOrStartApp(request *http.Request) (string, interface{}) {
	id, _ := strconv.ParseUint(mux.Vars(request)["id"], 10, 64)
	app := apiserver.QueryById(uint(id))
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
			app.AppStatus = resource.AppStop
		} else {
			deploy.Spec.Replicas = parseUtil.IntToInt32Pointer(app.Items[0].InstanceCount)
			app.AppStatus = resource.AppRunning
		}
		if err := k8sclient.UpdateResouce(&deploy); err != nil {
			return r.StatusInternalServerError, err
		}
	}
	apiserver.Update(app)
	return r.StatusOK, "ok"
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

	return r.StatusOK, "ok"
}

func DeleteApp(request *http.Request) (string, interface{}) {
	id, _ := strconv.ParseUint(mux.Vars(request)["id"], 10, 64)
	app := apiserver.QueryById(uint(id))
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

	apiserver.Delete(app)
	return r.StatusNoContent, "ok"
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

	cfgMap := configMap.NewConfigMap(app)
	if err := k8sclient.CreateResource(cfgMap); err != nil {
		return r.StatusInternalServerError, err
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
	apiserver.Insert(app)
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

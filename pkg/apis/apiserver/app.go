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
	"apiserver/pkg/resource/deployment"
	"apiserver/pkg/resource/service"
	r "apiserver/pkg/router"
	"apiserver/pkg/storage/cache"
	"apiserver/pkg/util/parseUtil"

	"github.com/gorilla/mux"
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

	if cache.ExsitResource(app.UserName, app.Items[0].Name, resource.ResourceKindService) {
		return r.StatusForbidden, "the service exist"
	}

	k8ssvc := service.NewService(app)
	svc, err := k8sclient.CreateService(k8ssvc)
	if err != nil {
		return r.StatusInternalServerError, err
	}

	k8sDeploy := deployment.NewDeployment(app)
	if err = k8sclient.CreateResource(k8sDeploy); err != nil {
		if err = k8sclient.DeleteResource(*svc); err != nil {
			return r.StatusInternalServerError, err
		}
		return r.StatusInternalServerError, err
	}
	external := fmt.Sprintf("http://%s:%v", configz.GetString("apiserver", "clusterNodes", "127.0.0.1"), svc.Spec.Ports[0].NodePort)
	app.External = external
	app.Items[0].External = external
	app.AppStatus = resource.AppBuilding

	apiserver.InsertApp(app)

	return r.StatusCreated, "ok"
}

func DeleteApp(request *http.Request) (string, interface{}) {
	id, _ := strconv.ParseUint(mux.Vars(request)["id"], 10, 64)
	namespace := mux.Vars(request)["namespace"]
	app := apiserver.QueryAppById(uint(id))

	for _, service := range app.Items {
		appName := service.Name

		if !cache.ExsitResource(namespace, appName, resource.ResourceKindDeployment) {
			return r.StatusNotFound, "application named " + appName + ` does't exist`
		}
		if err := k8sclient.DeleteResource(cache.Store.DeploymentCache.List[namespace][appName]); err != nil {
			return r.StatusInternalServerError, "delete application err: " + err.Error()
		}

		if !cache.ExsitResource(namespace, appName, resource.ResourceKindService) {
			return r.StatusNotFound, "application named " + appName + ` does't exist`
		}
		if err := k8sclient.DeleteResource(cache.Store.ServiceCache.List[namespace][appName]); err != nil {
			return r.StatusInternalServerError, "delete application err: " + err.Error()
		}
	}

	for _, svc := range app.Items {
		delete(cache.Store.ServiceCache.List[namespace], svc.Name)
		delete(cache.Store.DeploymentCache.List[namespace], svc.Name)
		for _, c := range svc.Items {
			delete(cache.Store.PodCache.List[namespace], c.Name)
		}
	}

	apiserver.DeleteApp(app)
	return r.StatusNoContent, "ok"
}

func StopOrStartOrRedeployApp(request *http.Request) (string, interface{}) {
	id, _ := strconv.ParseUint(mux.Vars(request)["id"], 10, 64)
	app := apiserver.QueryAppById(uint(id))

	for _, svc := range app.Items {
		appName := svc.Name
		namespace := mux.Vars(request)["namespace"]
		verb := mux.Vars(request)["verb"] //verb the action of app , start or stop

		if !cache.ExsitResource(namespace, appName, resource.ResourceKindDeployment) {
			return r.StatusNotFound, "service named " + appName + ` does't exist`
		}
		deploy := cache.Store.DeploymentCache.List[namespace][appName]
		if verb == "stop" {
			deploy.Spec.Replicas = parseUtil.IntToInt32Pointer(0)
			if err := k8sclient.UpdateResouce(&deploy); err != nil {
				return r.StatusInternalServerError, err
			}

			app.AppStatus = resource.AppStop
			svc.Status = resource.AppStop
			apiserver.UpdateApp(app)
			for _, container := range svc.Items {
				delete(cache.Store.PodCache.List[namespace], container.Name)
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
			apiserver.UpdateAppOnly(app)
			apiserver.UpdateServiceOnly(svc)

		}
		if verb == "redeploy" {
			pods, err := k8sclient.GetPods(namespace, svc.Name)
			if err != nil {
				return r.StatusInternalServerError, err
			}

			for _, pod := range pods {
				if err = k8sclient.DeleteResource(pod); err != nil {
					return r.StatusInternalServerError, err
				}
			}
			app.AppStatus = resource.AppBuilding
			svc.Status = resource.AppBuilding
			apiserver.UpdateAppOnly(app)
			apiserver.UpdateServiceOnly(svc)

			for _, container := range svc.Items {
				delete(cache.Store.PodCache.List[namespace], container.Name)
				apiserver.DeleteContainer(container)
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

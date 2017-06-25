package apiserver

import (
	"net/http"
	"strconv"
	"time"

	"apiserver/pkg/api/apiserver"
	"apiserver/pkg/resource"

	k8sclient "apiserver/pkg/resource/common"
	r "apiserver/pkg/router"
	"apiserver/pkg/storage/cache"

	"github.com/gorilla/mux"
	"k8s.io/client-go/pkg/api/v1"
)

// maximum number of lines loaded from the apiserver
var lineReadLimit int64 = 5000

// maximum number of bytes loaded from the apiserver
var byteReadLimit int64 = 500000

func GetContainers(request *http.Request) (string, interface{}) {
	// namespace := mux.Vars(request)["namespace"]
	pageCnt, _ := strconv.Atoi(request.FormValue("pageCnt"))
	pageNum, _ := strconv.Atoi(request.FormValue("pageNum"))
	containerName := request.FormValue("name")
	id, _ := strconv.ParseUint(mux.Vars(request)["id"], 10, 64)
	containers, total := apiserver.QueryContainers(containerName, pageCnt, pageNum, uint(id))
	return r.StatusOK, map[string]interface{}{"containers": containers, "total": total}
}

func RedeployContainer(request *http.Request) (string, interface{}) {
	id, _ := strconv.ParseUint(mux.Vars(request)["id"], 10, 64)
	container := apiserver.QueryContainerById(uint(id))
	svc := apiserver.QueryServiceById(container.ServiceId)
	namespace := mux.Vars(request)["namespace"]

	if !cache.ExsitResource(namespace, svc.Name, resource.ResourceKindDeployment) {
		return r.StatusNotFound, "service named " + svc.Name + ` does't exist`
	}

	pods, err := k8sclient.GetPods(namespace, svc.Name)
	if err != nil {
		return r.StatusInternalServerError, err
	}
	for _, pod := range pods {
		if pod.Name == container.Name {
			k8sclient.DeleteResource(&pod)
		}
	}

	for _, c := range svc.Items {
		if container.Name == c.Name {
			apiserver.DeleteContainer(c)
		}
	}
	svc = apiserver.QueryServiceById(container.ServiceId)
	svc.Status = resource.AppRunning
	apiserver.UpdateService(svc)
	// if err := ChangeContainerStatus(svc, namespace); err != nil {
	// 	return r.StatusInternalServerError, err
	// }
	return r.StatusCreated, "ok"
}

func ChangeContainerStatus(svc *apiserver.Service, namespace string) error {
	time.Sleep(5 * time.Second)
	podList, err := k8sclient.GetPods(namespace, svc.Name)
	if err != nil {
		return err
	}

	var containers []*apiserver.Container
	for _, pod := range podList {
		for _, c := range svc.Items {
			if c.Name != pod.Name {
				container := &apiserver.Container{Name: pod.ObjectMeta.Name, Image: svc.Image, Internal: pod.Status.PodIP, ServiceId: svc.ID}
				// if pod.Status.ContainerStatuses[0].Ready == true {
				// 	container.Status = resource.AppRunning
				// }
				containers = append(containers, container)
			}
		}
	}
	svc.Items = containers
	apiserver.UpdateService(svc)
	return nil
}

func GetContainerEvents(request *http.Request) (string, interface{}) {
	namespace := mux.Vars(request)["namespace"]
	containerName := mux.Vars(request)["name"]
	list, err := k8sclient.GetEventsForContainer(namespace, containerName)
	if err != nil {
		return r.StatusInternalServerError, err
	}
	return r.StatusOK, map[string]interface{}{"events": list}
}

func GetContainerLog(request *http.Request) (string, interface{}) {
	namespace := mux.Vars(request)["namespace"]
	podName := mux.Vars(request)["name"]
	containerName := cache.Store.PodCache.List[namespace][podName].Spec.Containers[0].Name
	logOptions := &v1.PodLogOptions{
		Container:  containerName,
		Follow:     false,
		Previous:   false,
		Timestamps: true,
		LimitBytes: &byteReadLimit,
		TailLines:  &lineReadLimit,
	}

	result, err := k8sclient.GetLogForContainer(namespace, podName, logOptions)
	if err != nil {
		return r.StatusInternalServerError, err
	}
	return r.StatusOK, map[string]interface{}{"logs": result}
}

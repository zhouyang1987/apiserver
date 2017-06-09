package apiserver

import (
	"net/http"
	"strconv"
	"time"

	"apiserver/pkg/api/apiserver"
	"apiserver/pkg/resource"
	k8sclient "apiserver/pkg/resource/common"
	"apiserver/pkg/resource/sync"
	r "apiserver/pkg/router"

	// "apiserver/pkg/util/log"

	"github.com/gorilla/mux"
)

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

	if _, exsit := sync.ListDeployment[namespace][svc.Name]; !exsit {
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
	if err := ChangeContainerStatus(svc, namespace); err != nil {
		return r.StatusInternalServerError, err
	}
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

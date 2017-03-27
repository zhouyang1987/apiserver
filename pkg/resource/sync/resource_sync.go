//this file is use to sync k8s's resource to memory ,then storage to the map[id]resources

package sync

import (
	"time"

	"apiserver/pkg/client"
	"apiserver/pkg/util/log"
	"apiserver/pkg/util/parseUtil"

	// "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/pkg/api/v1"
)

const (
	MAX_SIZE          = 1024
	LIST_WATCH_PERIOD = 30
)

var (
	ListNameSpace             = make(map[string]v1.Namespace, MAX_SIZE)
	ListService               = make(map[string]v1.Service, MAX_SIZE)
	ListReplicationController = make(map[string]v1.ReplicationController, MAX_SIZE)
)

//Sync ervery 30 Second to list k8s's resource to memory
func Sync() {
	ticker := time.NewTicker(time.Second * 30)
	for {
		select {
		case <-ticker.C:
			go ListResource()
		}
	}
}

//watch and list k8s's resource (namespace,service,replicationController) to resource memory
func ListResource() {
	nsList, err := client.K8sClient.
		CoreV1().
		Namespaces().
		List(v1.ListOptions{Watch: true, TimeoutSeconds: parseUtil.IntToInt64Pointer(LIST_WATCH_PERIOD)})
	if err != nil {
		log.Errorf("list and watch k8s's namespace err: %v", err)
	} else {
		loop(nsList)
	}

	for _, v := range ListNameSpace {
		svcList, err := client.K8sClient.
			CoreV1().
			Services(v.Name).
			List(v1.ListOptions{Watch: true, TimeoutSeconds: parseUtil.IntToInt64Pointer(LIST_WATCH_PERIOD)})
		if err != nil {
			log.Errorf("list and watch k8s's service of namespace [%v] err: %v", v.Name, err)
		} else {
			loop(svcList)
		}
	}

	for _, v := range ListNameSpace {
		rcList, err := client.K8sClient.
			CoreV1().
			ReplicationControllers(v.Name).
			List(v1.ListOptions{Watch: true, TimeoutSeconds: parseUtil.IntToInt64Pointer(LIST_WATCH_PERIOD)})
		if err != nil {
			log.Errorf("list and watch k8s's service of namespace [%v] err: %v", v.Name, err)
		} else {
			loop(rcList)
		}
	}
}

//loop add the k8s's resource (namespace,service,replicationController) to resource map
func loop(param interface{}) {
	switch param.(type) {
	case *v1.NamespaceList:
		for _, ns := range param.(*v1.NamespaceList).Items {
			ListNameSpace[""] = ns
		}
	case *v1.ServiceList:
		for _, svc := range param.(*v1.ServiceList).Items {
			ListService[""] = svc
		}
	case *v1.ReplicationControllerList:
		for _, rc := range param.(*v1.ReplicationControllerList).Items {
			ListReplicationController[""] = rc
		}
	}
}

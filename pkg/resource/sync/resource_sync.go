//this file is use to sync k8s's resource to memory ,then storage to the map[id]resources

package sync

import (
	"time"

	"apiserver/pkg/client"
	// "apiserver/pkg/util/jsonx"
	"apiserver/pkg/util/log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

const (
	MAX_SIZE          = 1024
	LIST_WATCH_PERIOD = 30
)

var (
	ListNameSpace  = make(map[string]map[string]v1.Namespace, MAX_SIZE)
	ListService    = make(map[string]map[string]v1.Service, MAX_SIZE)
	ListDeployment = make(map[string]map[string]extensions.Deployment, MAX_SIZE)
)

//Sync ervery 30 Second to list k8s's resource to memory
func Sync() {
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:
			/*	log.Debugf("ListNameSpace ==== %v", jsonx.ToJson(ListNameSpace))
				log.Debugf("ListService ==== %v", jsonx.ToJson(ListService))*/
			// log.Debugf("ListDeployment ==== %v", jsonx.ToJson(ListDeployment["huangjia"]))
			go ListResource()
		}
	}
}

//watch and list k8s's resource (namespace,service,replicationController) to resource memory
func ListResource() {
	nsList, err := client.K8sClient.
		CoreV1().
		Namespaces().
		List(metav1.ListOptions{})
	if err != nil {
		log.Errorf("list and watch k8s's namespace err: %v", err)
		return
	} else {
		if len(nsList.Items) > 0 {
			loop(nsList, "")
		}
	}
	for k, v := range ListNameSpace {
		svcList, err := client.K8sClient.
			CoreV1().
			Services(v[k].ObjectMeta.Name).
			List(metav1.ListOptions{})
		if err != nil {
			log.Errorf("list and watch k8s's service of namespace [%v] err: %v", v[k].Name, err)
		} else {
			loop(svcList, v[k].ObjectMeta.Name)
		}
	}
	for k, v := range ListNameSpace {
		dpList, err := client.K8sClient.
			ExtensionsV1beta1().
			Deployments(v[k].ObjectMeta.Name).
			List(metav1.ListOptions{})
		if err != nil {
			log.Errorf("list and watch k8s's service of namespace [%v] err: %v", v[k].Name, err)
		} else {
			loop(dpList, v[k].ObjectMeta.Name)
		}
	}
}

//loop add the k8s's resource (namespace,service,replicationController) to resource map
func loop(param interface{}, nsname string) {
	switch param.(type) {
	case *v1.NamespaceList:
		for _, ns := range param.(*v1.NamespaceList).Items {
			nsmap := make(map[string]v1.Namespace, MAX_SIZE)
			nsmap[ns.ObjectMeta.Name] = ns
			ListNameSpace[ns.ObjectMeta.Name] = nsmap
		}
	case *v1.ServiceList:
		items := param.(*v1.ServiceList).Items
		if len(items) == 0 {
			ListService[nsname] = make(map[string]v1.Service, MAX_SIZE)
		} else {
			svcmap := make(map[string]v1.Service, MAX_SIZE)
			for _, svc := range items {
				svcmap[svc.ObjectMeta.Name] = svc
			}
			ListService[nsname] = svcmap
		}

	case *extensions.DeploymentList:
		items := param.(*extensions.DeploymentList).Items
		if len(items) == 0 {
			ListDeployment[nsname] = make(map[string]extensions.Deployment, MAX_SIZE)
		} else {
			dpmap := make(map[string]extensions.Deployment, MAX_SIZE)
			for _, deploy := range items {
				dpmap[deploy.ObjectMeta.Name] = deploy
			}
			ListDeployment[nsname] = dpmap
		}
	}
}

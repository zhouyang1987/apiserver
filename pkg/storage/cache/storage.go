package cache

import (
	"sync"
	"time"

	"apiserver/pkg/client"
	"apiserver/pkg/resource"
	"apiserver/pkg/util/log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

const (
	MAX_SIZE          = 1024
	LIST_WATCH_PERIOD = 30
)

var Store *Cache

type Cache struct {
	*NamespaceCache
	*ServiceCache
	*DeploymentCache
	*ConfigMapCache
}

type NamespaceCache struct {
	sync.RWMutex
	List map[string]map[string]v1.Namespace
}

type ServiceCache struct {
	sync.RWMutex
	List map[string]map[string]v1.Service
}

type DeploymentCache struct {
	sync.RWMutex
	List map[string]map[string]extensions.Deployment
}

type ConfigMapCache struct {
	sync.RWMutex
	List map[string]map[string]v1.ConfigMap
}

func init() {
	Store = &Cache{
		&NamespaceCache{
			List: make(map[string]map[string]v1.Namespace, MAX_SIZE),
		},
		&ServiceCache{
			List: make(map[string]map[string]v1.Service, MAX_SIZE),
		},
		&DeploymentCache{
			List: make(map[string]map[string]extensions.Deployment, MAX_SIZE),
		},
		&ConfigMapCache{
			List: make(map[string]map[string]v1.ConfigMap, MAX_SIZE),
		},
	}
}

//Sync ervery 30 Second to list k8s's resource to memory
func List() {
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:
			go listResource()
		}
	}
}

//watch and list k8s's resource (namespace,service,replicationController) to resource memory
func listResource() {
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
	for k, v := range Store.NamespaceCache.List {
		svcList, err := client.K8sClient.
			CoreV1().
			Services(v[k].ObjectMeta.Name).
			List(metav1.ListOptions{})
		if err != nil {
			log.Errorf("list and watch k8s's service of namespace [%v] err: %v", v[k].Name, err)
		} else {
			loop(svcList, v[k].ObjectMeta.Name)
		}

		dpList, err := client.K8sClient.
			ExtensionsV1beta1().
			Deployments(v[k].ObjectMeta.Name).
			List(metav1.ListOptions{})
		if err != nil {
			log.Errorf("list and watch k8s's deployment of namespace [%v] err: %v", v[k].Name, err)
		} else {
			loop(dpList, v[k].ObjectMeta.Name)
		}

		cfgMapList, err := client.K8sClient.
			CoreV1().
			ConfigMaps(v[k].ObjectMeta.Name).
			List(metav1.ListOptions{})
		if err != nil {
			log.Errorf("list and watch k8s's configMap of namespace [%v] err: %v", v[k].Name, err)
		} else {
			loop(cfgMapList, v[k].ObjectMeta.Name)
		}
	}
}

//loop add the k8s's resource (namespace,service,replicationController) to resource map
func loop(param interface{}, nsname string) {
	switch param.(type) {
	case *v1.NamespaceList:
		for _, ns := range param.(*v1.NamespaceList).Items {
			Store.NamespaceCache.Lock()
			nsmap := make(map[string]v1.Namespace)
			nsmap[ns.ObjectMeta.Name] = ns
			Store.NamespaceCache.List[ns.ObjectMeta.Name] = nsmap
			Store.NamespaceCache.Unlock()
		}
	case *v1.ServiceList:
		items := param.(*v1.ServiceList).Items
		if len(items) == 0 {
			Store.ServiceCache.List[nsname] = make(map[string]v1.Service)
		} else {
			Store.ServiceCache.Lock()
			svcmap := make(map[string]v1.Service)
			for _, svc := range items {
				svcmap[svc.ObjectMeta.Name] = svc
			}
			Store.ServiceCache.List[nsname] = svcmap
			Store.ServiceCache.Unlock()
		}

	case *extensions.DeploymentList:
		items := param.(*extensions.DeploymentList).Items
		if len(items) == 0 {
			Store.DeploymentCache.List[nsname] = make(map[string]extensions.Deployment)
		} else {
			Store.DeploymentCache.Lock()
			dpmap := make(map[string]extensions.Deployment)
			for _, deploy := range items {
				dpmap[deploy.ObjectMeta.Name] = deploy
			}
			Store.DeploymentCache.List[nsname] = dpmap
			Store.DeploymentCache.Unlock()
		}
	case *v1.ConfigMapList:
		items := param.(*v1.ConfigMapList).Items
		if len(items) == 0 {
			Store.ConfigMapCache.List[nsname] = make(map[string]v1.ConfigMap)
		} else {
			Store.ConfigMapCache.Lock()
			cfgmap := make(map[string]v1.ConfigMap, MAX_SIZE)
			for _, configMap := range items {
				cfgmap[configMap.ObjectMeta.Name] = configMap
			}
			Store.ConfigMapCache.List[nsname] = cfgmap
			Store.ConfigMapCache.Unlock()
		}
	}
}

func ExsitResource(namespace, resourceName, resourceKind string) bool {
	switch resourceKind {
	case resource.ResourceKindNamespace:
		_, exist := Store.NamespaceCache.List[resourceName]
		return exist
	case resource.ResourceKindService:
		_, exist := Store.ServiceCache.List[namespace][resourceName]
		return exist
	case resource.ResourceKindDeployment:
		_, exist := Store.DeploymentCache.List[namespace][resourceName]
		return exist
	case resource.ResourceKindConfigMap:
		_, exist := Store.ConfigMapCache.List[namespace][resourceName]
		return exist
	}
	return false
}

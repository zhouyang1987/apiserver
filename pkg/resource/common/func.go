package common

import (
	"strings"

	"apiserver/pkg/api/apiserver"
	"apiserver/pkg/client"
	"apiserver/pkg/resource"
	"apiserver/pkg/resource/event"
	"apiserver/pkg/util/log"
	"apiserver/pkg/util/parseUtil"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

//CreateResource create namespace,service,replicationController
func CreateResource(param interface{}) error {
	switch param.(type) {
	case *v1.Namespace:
		ns := param.(*v1.Namespace)
		_, err := client.K8sClient.
			CoreV1().
			Namespaces().
			Create(ns)
		if err != nil {
			log.Errorf("create namespace [%v] err:%v", ns.Name, err)
			return err
		}
		log.Noticef("namespace [%v] is created]", ns.Name)
		return nil
	case *v1.Service:
		svc := param.(*v1.Service)
		_, err := client.K8sClient.
			CoreV1().
			Services(svc.Namespace).
			Create(svc)
		if err != nil {
			log.Errorf("create service [%v] err:%v", svc.Name, err)
			return err
		}
		log.Noticef("service [%v] is created]", svc.Name)
		return nil
	case *v1.ConfigMap:
		cfgMap := param.(*v1.ConfigMap)
		_, err := client.K8sClient.
			CoreV1().
			ConfigMaps(cfgMap.Namespace).
			Create(cfgMap)
		if err != nil {
			log.Errorf("create configMap [%v] err:%v", cfgMap.Name, err)
			return err
		}
		log.Noticef("configMap [%v] is created]", cfgMap.Name)
		return nil
	case *v1.ReplicationController:
		rc := param.(*v1.ReplicationController)
		_, err := client.K8sClient.
			CoreV1().
			ReplicationControllers(rc.Namespace).
			Create(rc)
		if err != nil {
			log.Errorf("create replicationControllers [%v] err:%v", rc.Name, err)
			return err
		}
		log.Noticef("replication [%v] is created]", rc.Name)
		return nil
	case *extensions.Deployment:
		deploy := param.(*extensions.Deployment)
		_, err := client.K8sClient.ExtensionsV1beta1().Deployments(deploy.Namespace).Create(deploy)
		if err != nil {
			log.Errorf("create deployment [%v] err:%v", deploy.Name, err)
			return err
		}
		log.Noticef("deployment [%v] is created]", deploy.Name)
		return nil
	}

	return nil
}

//ExsitResource decide namesapce,service,replicationController exsit or not by name;false is not exsit,true exsit
func ExsitResource(param interface{}) bool {
	switch param.(type) {
	case v1.Namespace:
		_, err := client.K8sClient.
			CoreV1().
			Namespaces().
			Get(param.(v1.Namespace).Name, metav1.GetOptions{})
		if err != nil {
			return false
		}
		return true
	case v1.Service:
		svc := param.(v1.Service)
		_, err := client.K8sClient.
			CoreV1().
			Services(svc.Namespace).
			Get(svc.Name, metav1.GetOptions{})
		if err != nil {
			return false
		}
		return true

	case v1.ConfigMap:
		cfgMap := param.(v1.Service)
		_, err := client.K8sClient.
			CoreV1().
			ConfigMaps(cfgMap.Namespace).
			Get(cfgMap.Name, metav1.GetOptions{})
		if err != nil {
			return false
		}
		return true
	case v1.ReplicationController:
		rc := param.(v1.ReplicationController)
		_, err := client.K8sClient.
			CoreV1().
			ReplicationControllers(rc.Namespace).
			Get(rc.Name, metav1.GetOptions{})
		if err != nil {
			return false
		}
		return true
	case extensions.Deployment:
		deploy := param.(*extensions.Deployment)
		_, err := client.K8sClient.
			ExtensionsV1beta1().
			Deployments(deploy.Namespace).
			Get(deploy.Name, metav1.GetOptions{})
		if err != nil {
			return false
		}
		return true
	}
	return false
}

//DeleteResource delete namespace,service,replicationController
func DeleteResource(param interface{}) error {
	switch param.(type) {
	case v1.Namespace:
		ns := param.(v1.Namespace)
		err := client.K8sClient.
			CoreV1().
			Namespaces().
			Delete(ns.Name, &metav1.DeleteOptions{TypeMeta: resource.NewTypeMeta(resource.ResourceKindNamespace, "v1"), GracePeriodSeconds: parseUtil.IntToInt64Pointer(30)})
		if err != nil {
			log.Errorf("delete namespace [%v] err:%v", ns.Name, err)
			return err
		}
		log.Noticef("namespace [%v] was deleted", ns.Name)
		return nil
	case v1.Service:
		svc := param.(v1.Service)
		err := client.K8sClient.
			CoreV1().
			Services(svc.Namespace).
			Delete(svc.Name, &metav1.DeleteOptions{})
		if err != nil {
			log.Errorf("delete service [%v] err:%v", svc.Name, err)
			return err
		}
		log.Noticef("service [%v] was deleted]", svc.Name)
		return nil

	case v1.ConfigMap:
		cfgMap := param.(v1.ConfigMap)
		err := client.K8sClient.
			CoreV1().
			ConfigMaps(cfgMap.Namespace).
			Delete(cfgMap.Name, &metav1.DeleteOptions{})
		if err != nil {
			log.Errorf("delete configMap [%v] err:%v", cfgMap.Name, err)
			return err
		}
		log.Noticef("configMap [%v] was deleted]", cfgMap.Name)
		return nil
	case v1.ReplicationController:
		rc := param.(v1.ReplicationController)
		err := client.K8sClient.
			CoreV1().
			ReplicationControllers(rc.Namespace).
			Delete(rc.Name, &metav1.DeleteOptions{TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "ReplicationController"}, OrphanDependents: parseUtil.BoolToPointer(false)})
		if err != nil {
			log.Errorf("delete replicationControllers [%v] err:%v", rc.Name, err)
			return err
		}
		log.Noticef("replication [%v] is delete]", rc.Name)
		return nil
	case extensions.Deployment:
		backend := new(metav1.DeletionPropagation)
		*backend = metav1.DeletePropagationForeground
		deploy := param.(extensions.Deployment)
		err := client.K8sClient.
			ExtensionsV1beta1().
			Deployments(deploy.Namespace).
			Delete(deploy.Name, &metav1.DeleteOptions{PropagationPolicy: backend})
		if err != nil {
			log.Errorf("delete deployment [%v] err:%v", deploy.Name, err)
			return err
		}
		log.Noticef("deployment [%v] was deleted]", deploy.Name)
		return nil
	case v1.Pod:
		backend := new(metav1.DeletionPropagation)
		*backend = metav1.DeletePropagationForeground
		pod := param.(v1.Pod)
		err := client.K8sClient.
			CoreV1().
			Pods(pod.Namespace).
			Delete(pod.Name, &metav1.DeleteOptions{PropagationPolicy: backend})
		if err != nil {
			log.Errorf("delete pod [%v] err:%v", pod.Name, err)
			return err
		}
		log.Noticef("pod [%v] was deleted]", pod.Name)
		return nil
	}
	return nil
}

func UpdateResouce(param interface{}) error {
	switch param.(type) {
	case *v1.Namespace:
		ns := param.(*v1.Namespace)
		_, err := client.K8sClient.
			CoreV1().
			Namespaces().
			Update(ns)
		if err != nil {
			log.Errorf("update namespace [%v] err:%v", ns.Name, err)
			return err
		}
		log.Noticef("namespace [%v] was updated", ns.Name)
		return nil
	case *v1.Service:
		svc := param.(*v1.Service)
		_, err := client.K8sClient.
			CoreV1().
			Services(svc.Namespace).
			Update(svc)
		if err != nil {
			log.Errorf("update service [%v] err:%v", svc.Name, err)
			return err
		}
		log.Noticef("service [%v] was updated]", svc.Name)
		return nil
	case *v1.ReplicationController:
		rc := param.(*v1.ReplicationController)
		_, err := client.K8sClient.
			CoreV1().
			ReplicationControllers(rc.Namespace).
			Update(rc)
		if err != nil {
			log.Errorf("update replicationControllers [%v] err:%v", rc.Name, err)
			return err
		}
		log.Noticef("replication [%v] is updated]", rc.Name)
		return nil

	case *extensions.Deployment:
		deploy := param.(*extensions.Deployment)
		_, err := client.K8sClient.
			ExtensionsV1beta1().
			Deployments(deploy.Namespace).Update(deploy)
		if err != nil {
			log.Errorf("update replicationControllers [%v] err:%v", deploy.Name, err)
			return err
		}
		log.Noticef("replication [%v] is updated]", deploy.Name)
		return nil
	case *v1.ConfigMap:
		cfgMap := param.(*v1.ConfigMap)
		_, err := client.K8sClient.CoreV1().ConfigMaps(cfgMap.Namespace).Update(cfgMap)
		if err != nil {
			log.Errorf("update configMap [%v] err:%v", cfgMap.Name, err)
			return err
		}
		log.Noticef("configMap [%v] is updated]", cfgMap.Name)
		return nil
	}
	return nil
}

func GetPods(namespace, deployName string) ([]v1.Pod, error) {
	list, err := client.K8sClient.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: "name=" + deployName})
	if err != nil {
		return []v1.Pod{}, err
	}
	return list.Items, nil
}

func CreateService(svc *v1.Service) (*v1.Service, error) {
	service, err := client.K8sClient.
		CoreV1().
		Services(svc.Namespace).
		Create(svc)
	if err != nil {
		log.Errorf("create service [%v] err:%v", svc.Name, err)
		return nil, err
	}
	log.Noticef("service [%v] was created]", svc.Name)
	return service, nil
}

func GetDeploymentPods(appName, namespace string) ([]v1.Pod, error) {
	list, err := client.K8sClient.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: "name=" + appName})
	if err != nil {
		log.Errorf("get deployment [%v]'s pods err:%v", appName, err)
		return nil, err
	}
	return list.Items, nil
}

func GetEventsForContainer(namespace, containerName string) (list []event.Event, err error) {
	listEvent, err := client.K8sClient.CoreV1().Events(namespace).List(resource.ListEverything)
	if err != nil {
		log.Errorf("get pod [%v]'s event err:%v", containerName, err)
		return
	}
	for _, ev := range listEvent.Items {
		if strings.Contains(ev.Name, containerName) {
			list = append(
				list,
				event.Event{
					Reason:        ev.Reason,
					Type:          ev.Type,
					LastTimestamp: ev.LastTimestamp,
					Message:       ev.Message,
				},
			)
		}
	}
	return
}

func GetEventsForService(namespace, serviceName string) (list []event.Event, err error) {
	listEvent, err := client.K8sClient.CoreV1().Events(namespace).List(resource.ListEverything)
	if err != nil {
		log.Errorf("get service [%v]'s event err:%v", serviceName, err)
		return
	}

	svcs, _ := apiserver.QueryServices(serviceName, 100000, 0, 0)
	if len(svcs[0].Items) > 0 {
		for _, ev := range listEvent.Items {
			if strings.Contains(ev.Name, svcs[0].Items[0].Name) {
				list = append(
					list,
					event.Event{
						Reason:        ev.Reason,
						Type:          ev.Type,
						LastTimestamp: ev.LastTimestamp,
						Message:       ev.Message,
					},
				)
			}

			if strings.Contains(ev.Name, serviceName) {
				list = append(
					list,
					event.Event{
						Reason:        ev.Reason,
						Type:          ev.Type,
						LastTimestamp: ev.LastTimestamp,
						Message:       ev.Message,
					},
				)
			}
		}

	}
	return
}

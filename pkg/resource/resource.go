// Copyright © 2017 huang jia <449264675@qq.com>
//
// Licensed under the Apache License, Vereion 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resouce

import (
	"apiserver/pkg/api/application"
	"apiserver/pkg/client"
	"apiserver/pkg/util/log"
	"apiserver/pkg/util/parseUtil"

	"k8s.io/client-go/pkg/api/resource"
	"k8s.io/client-go/pkg/api/v1"
	metav1 "k8s.io/client-go/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/util/intstr"
)

func newTypeMeta(kind, vereion string) metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       kind,
		APIVersion: vereion,
	}
}

func newOjectMeta(app *application.App) v1.ObjectMeta {
	return v1.ObjectMeta{
		Name:      app.UserName,
		Namespace: app.UserName,
		Labels:    map[string]string{"name": app.Name},
	}
}

func NewPodSpec(app *application.App) v1.PodSpec {
	var containerPorts []v1.ContainerPort
	if app.Ports != nil {
		for _, port := range app.Ports {
			containerPorts = append(containerPorts, v1.ContainerPort{
				HostPort:      int32(port.TargetPort),
				ContainerPort: int32(port.TargetPort),
				Protocol:      v1.Protocol(port.Schame),
			})
		}
	}

	return v1.PodSpec{
		RestartPolicy: v1.RestartPolicyAlways,
		Containers: []v1.Container{
			v1.Container{
				Name:            app.Name,
				Image:           app.Image,
				Command:         app.Command,
				Ports:           containerPorts,
				ImagePullPolicy: v1.PullIfNotPresent,
				Resources: v1.ResourceRequirements{
					Limits: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse(app.Cpu),    //TODO 根据前端传入的值做资源限制
						v1.ResourceMemory: resource.MustParse(app.Memory), //TODO 根据前端传入的值做资源限制
					},
					Requests: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse(app.Cpu),
						v1.ResourceMemory: resource.MustParse(app.Memory),
					},
				},
				VolumeMounts: []v1.VolumeMount{
				/*	v1.VolumeMount{
					Name:      app.Mount.Name,
					MountPath: app.Mount.MountPath,
					SubPath:   app.Mount.SubPath,
					ReadOnly:  app.Mount.ReadOnly,
				},*/
				},
			},
		},
	}
}

func newPodTemplateSpec(app *application.App) *v1.PodTemplateSpec {
	return &v1.PodTemplateSpec{
		ObjectMeta: newOjectMeta(app),
		Spec:       NewPodSpec(app),
	}
}

func newReplicationControllerepec(app *application.App) v1.ReplicationControllerSpec {
	return v1.ReplicationControllerSpec{
		Replicas: parseUtil.Int32ToPointer(int32(app.InstanceCount)),
		Selector: map[string]string{"name": app.Name},
		Template: newPodTemplateSpec(app),
	}
}

func newServiceSpec(app *application.App) v1.ServiceSpec {
	/*	var svcPorts []v1.ServicePort
		for _, port := range app.Ports {
			svcPorts = append(svcPorts, v1.ServicePort{
				Name:       app.Name,
				Port:       int32(port.ServicePort),
				TargetPort: intstr.FromInt(port.TargetPort),
				Protocol:   v1.Protocol(port.Schame),
			})
		}*/
	return v1.ServiceSpec{
		Selector: map[string]string{"name": app.Name},
		Ports: []v1.ServicePort{
			v1.ServicePort{
				Name:       app.Name,
				Port:       int32(800),
				TargetPort: intstr.FromInt(8080),
				Protocol:   v1.ProtocolTCP,
			},
		},
	}
}

func newNamespaceSpec(app *application.App) v1.NamespaceSpec {
	return v1.NamespaceSpec{
		Finalizers: []v1.FinalizerName{v1.FinalizerKubernetes},
	}
}
func NewSVC(app *application.App) *v1.Service {
	return &v1.Service{
		TypeMeta:   newTypeMeta("Service", "v1"),
		ObjectMeta: newOjectMeta(app),
		Spec:       newServiceSpec(app),
	}
}

func NewRC(app *application.App) *v1.ReplicationController {
	return &v1.ReplicationController{
		TypeMeta:   newTypeMeta("ReplicationController", "v1"),
		ObjectMeta: newOjectMeta(app),
		Spec:       newReplicationControllerepec(app),
	}
}

func NewNS(app *application.App) *v1.Namespace {
	return &v1.Namespace{
		TypeMeta:   newTypeMeta("Namespace", "v1"),
		ObjectMeta: newOjectMeta(app),
		Spec:       newNamespaceSpec(app),
	}
}

//CreateResource create namespace,service,replicationController
func CreateResource(param interface{}) error {
	switch param.(type) {
	case *v1.Namespace:
		ns := param.(*v1.Namespace)
		_, err := client.K8sClient.CoreV1().Namespaces().Create(ns)
		if err != nil {
			log.Errorf("create namespace [%v] err:%v", ns.Name, err)
			return err
		}
		log.Noticef("namespace [%v] is created]", ns.Name)
		return nil

	case *v1.Service:
		svc := param.(*v1.Service)
		_, err := client.K8sClient.CoreV1().Services(svc.Namespace).Create(svc)
		if err != nil {
			log.Errorf("create service [%v] err:%v", svc.Name, err)
			return err
		}
		log.Noticef("service [%v] is created]", svc.Name)
		return nil

	case *v1.ReplicationController:
		rc := param.(*v1.ReplicationController)
		_, err := client.K8sClient.CoreV1().ReplicationControllers(rc.Namespace).Create(rc)
		if err != nil {
			log.Errorf("create replicationControllers [%v] err:%v", rc.Name, err)
			return err
		}
		log.Noticef("replication [%v] is created]", rc.Name)
		return nil
	}

	return nil
}

//ExsitResource decide namesapce,service,replicationController exsit or not by name
func ExsitResource(param interface{}) bool {
	switch param.(type) {
	case *v1.Namespace:
		_, err := client.K8sClient.CoreV1().Namespaces().Get(param.(*v1.Namespace).Name, metav1.GetOptions{})
		if err != nil {
			return false
		}
		return true

	case *v1.Service:
		svc := param.(*v1.Service)
		_, err := client.K8sClient.CoreV1().Services(svc.Namespace).Get(svc.Name, metav1.GetOptions{})
		if err != nil {
			return false
		}
		return true

	case *v1.ReplicationController:
		rc := param.(*v1.ReplicationController)
		_, err := client.K8sClient.CoreV1().ReplicationControllers(rc.Namespace).Get(rc.Name, metav1.GetOptions{})
		if err != nil {
			return false
		}
		return true
	}

	return false
}

package service

import (
	"apiserver/pkg/api/apiserver"
	"apiserver/pkg/resource"

	"k8s.io/client-go/pkg/api/v1"
)

func NewServiceSpec(svc *apiserver.Service) v1.ServiceSpec {
	var ports [][]v1.ServicePort
	if svc.Config.SuperConfig != nil {
		for _, port := range svc.Config.SuperConfig.Ports {
			ports = append(ports, v1.ServicePort{port.ServicePort, port.ContainerPort, port.Protocol})
		}
	}
	return v1.ServiceSpec{
		Ports:          ports,
		Selector:       map[string]string{"name": svc.Name},
		Type:           v1.ServiceTypeLoadBalancer,
		LoadBalancerIP: svc.LoadbalanceIp,
	}
}

func NewService(svc *apiserver.Service) *v1.Service {
	return &v1.Service{
		TypeMeta:   resource.NewTypeMeta(resource.ResourceKindDeployment, "extensions/v1beta1"),
		ObjectMeta: resource.NewObjectMeta(svc),
		Spec:       NewServiceSpec(svc),
	}
}

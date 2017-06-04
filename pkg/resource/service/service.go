package service

import (
	"apiserver/pkg/api/apiserver"
	"apiserver/pkg/resource"

	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/pkg/api/v1"
)

func NewServiceSpec(app *apiserver.App) v1.ServiceSpec {
	var ports []v1.ServicePort
	svc := app.Items[0]
	if svc.Config.SuperConfig != nil {
		for _, port := range svc.Config.SuperConfig.Ports {
			ports = append(ports, v1.ServicePort{Port: int32(port.ServicePort), TargetPort: intstr.FromInt(port.ContainerPort), Protocol: v1.Protocol(port.Protocol)})
		}
	}
	return v1.ServiceSpec{
		Ports:    ports,
		Selector: map[string]string{"name": svc.Name},
		Type:     v1.ServiceTypeNodePort,
		// LoadBalancerIP: svc.LoadbalanceIp,
	}
}

func NewService(app *apiserver.App) *v1.Service {
	return &v1.Service{
		TypeMeta:   resource.NewTypeMeta(resource.ResourceKindDeployment, "extensions/v1beta1"),
		ObjectMeta: resource.NewObjectMeta(app),
		Spec:       NewServiceSpec(app),
	}
}

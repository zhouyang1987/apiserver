package deployment

import (
	"apiserver/pkg/api/apiserver"
	"apiserver/pkg/resource"
	"apiserver/pkg/resource/pod"
	"apiserver/pkg/util/parseutil"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

func NewDeploymentSpce(svc *apiserver.Service) extensions.DeploymentSpec {
	return extensions.DeploymentSpec{
		Replicas: parseUtil.IntToInt32Pointer(svc.InstanceCount),
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"name": svc.Name,
			},
		},
		Template: v1.PodTemplateSpec{
			ObjectMeta: resource.NewObjectMeta(svc),
			Spec:       pod.NewPodSpec(svc),
		},
		Strategy: extensions.DeploymentStrategy{
			Type: extensions.RecreateDeploymentStrategyType,
			// RollingUpdate: extensions.RollingUpdateDeployment{},
		},
	}
}

func NewDeployment(svc *apiserver.Service) *extensions.Deployment {
	return &extensions.Deployment{
		TypeMeta:   resource.NewTypeMeta(resource.ResourceKindDeployment, "extensions/v1beta1"),
		ObjectMeta: resource.NewObjectMeta(svc),
		Spec:       NewDeploymentSpce(svc),
	}
}

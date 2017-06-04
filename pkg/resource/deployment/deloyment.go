package deployment

import (
	"apiserver/pkg/api/apiserver"
	"apiserver/pkg/resource"
	"apiserver/pkg/resource/pod"
	"apiserver/pkg/util/parseUtil"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

func NewDeploymentSpce(app *apiserver.App) extensions.DeploymentSpec {
	return extensions.DeploymentSpec{
		Replicas: parseUtil.IntToInt32Pointer(app.Items[0].InstanceCount),
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"name": app.Items[0].Name,
			},
		},
		Template: v1.PodTemplateSpec{
			ObjectMeta: resource.NewObjectMeta(app),
			Spec:       pod.NewPodSpec(app),
		},
		Strategy: extensions.DeploymentStrategy{
			Type: extensions.RecreateDeploymentStrategyType,
			// RollingUpdate: extensions.RollingUpdateDeployment{},
		},
	}
}

func NewDeployment(app *apiserver.App) *extensions.Deployment {
	return &extensions.Deployment{
		TypeMeta:   resource.NewTypeMeta(resource.ResourceKindDeployment, "extensions/v1beta1"),
		ObjectMeta: resource.NewObjectMeta(app),
		Spec:       NewDeploymentSpce(app),
	}
}

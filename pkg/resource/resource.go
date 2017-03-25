package resouce

import (
	"apiserver/pkg/api/application"

	"k8s.io/client-go/pkg/api/v1"
	metav1 "k8s.io/client-go/pkg/apis/meta/v1"
)

func newTypeMeta(kind, version string) *metav1.TypeMeta {
	return &metav1.TypeMeta{
		Kind:       kind,
		APIVersion: version,
	}
}

func newOjectMeta(app *application.App) *v1.ObjectMeta {
	return &v1.ObjectMeta{
		Name:      app.Name,
		Namespace: app.UserName,
	}
}

func newPodTemplateSpec(app *application.App) *v1.PodTemplateSpec {
	return &v1.PodTemplateSpec{}
}

func newReplicationControllerSpec(app *application.App) *v1.ReplicationControllerSpec {
	return &v1.ReplicationControllerSpec{}
}

func newServiceSpec(app *application.App) *v1.ServiceSpec {
	return &v1.ServiceSpec{}
}

func newNamespaceSpec(app *application.App) *v1.NamespaceSpec {
	return &v1.NamespaceSpec{}
}
func NewSVC(app *application.App) *v1.Service {
	return &v1.Service{}
}

func NewRC(app *application.App) *v1.ReplicationController {
	return &v1.ReplicationController{}
}

func NewNS(app *application.App) *v1.Namespace {
	return &v1.Namespace{}
}

package namespace

import (
	"apiserver/pkg/api/apiserver"
	"apiserver/pkg/resource"

	"k8s.io/client-go/pkg/api/v1"
)

func NewNamespaceSpec(svc *apiserver.Service) v1.NamespaceSpec {
	return v1.NamespaceSpec{
		Finalizers: []v1.FinalizerName{v1.FinalizerKubernetes},
	}
}

func NewNamespace(svc *apiserver.Service) *v1.Namespace {
	return &v1.Namespace{
		TypeMeta:   resource.NewTypeMeta(resource.ResourceKindNamespace, "v1"),
		ObjectMeta: resource.NewObjectMeta(svc),
		Spec:       NewNamespaceSpec(svc),
	}
}

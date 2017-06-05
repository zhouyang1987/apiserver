package configMap

import (
	"apiserver/pkg/api/apiserver"
	"apiserver/pkg/resource"

	"k8s.io/client-go/pkg/api/v1"
)

func NewConfigMap(app *apiserver.App) *v1.ConfigMap {
	return &v1.ConfigMap{
		TypeMeta:   resource.NewTypeMeta(resource.ResourceKindConfigMap, "v1"),
		ObjectMeta: resource.NewObjectMeta(app),
		Data:       map[string]string{app.Items[0].Config.ConfigMap.Name: app.Items[0].Config.ConfigMap.Content},
	}
}

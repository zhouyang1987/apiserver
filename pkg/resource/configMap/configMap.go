package configMap

import (
	"apiserver/pkg/api/apiserver"
	"apiserver/pkg/resource"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

func NewConfigMap(app *apiserver.App) *v1.ConfigMap {

	data := map[string]string{}

	if app.Items[0].Config.ConfigGroup != nil {
		for _, configMap := range app.Items[0].Config.ConfigGroup.ConfigMaps {
			data[configMap.Name] = configMap.Content
		}
	}

	return &v1.ConfigMap{
		TypeMeta:   resource.NewTypeMeta(resource.ResourceKindConfigMap, "v1"),
		ObjectMeta: resource.NewObjectMeta(app),
		Data:       data,
	}
}

func NewConfigMapByService(svc *apiserver.Service, namespace string) *v1.ConfigMap {
	data := map[string]string{}

	if svc.Config.ConfigGroup != nil {
		for _, configMap := range svc.Config.ConfigGroup.ConfigMaps {
			data[configMap.Name] = configMap.Content
		}
	}

	return &v1.ConfigMap{
		TypeMeta: resource.NewTypeMeta(resource.ResourceKindConfigMap, "v1"),
		ObjectMeta: metav1.ObjectMeta{
			Name:        svc.Name,
			Namespace:   namespace,
			Labels:      map[string]string{"name": svc.Name},
			Annotations: map[string]string{"name": svc.Name},
		},
		Data: data,
	}
}

func NewConfigMapByConfig(c *apiserver.ConfigGroup) v1.ConfigMap {
	datas := map[string]string{}
	for _, cfg := range c.ConfigMaps {
		datas[cfg.Name] = cfg.Content
	}

	return v1.ConfigMap{
		TypeMeta: resource.NewTypeMeta(resource.ResourceKindConfigMap, "v1"),
		ObjectMeta: metav1.ObjectMeta{
			Name:        c.Name,
			Namespace:   c.Namespace,
			Labels:      map[string]string{"name": c.Name},
			Annotations: map[string]string{"name": c.Name},
		},
		Data: datas,
	}
}

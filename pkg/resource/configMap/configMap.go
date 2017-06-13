package configMap

import (
	"apiserver/pkg/api/apiserver"
	"apiserver/pkg/resource"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

func NewConfigMap(app *apiserver.App) *v1.ConfigMap {
	return &v1.ConfigMap{
		TypeMeta:   resource.NewTypeMeta(resource.ResourceKindConfigMap, "v1"),
		ObjectMeta: resource.NewObjectMeta(app),
		Data:       map[string]string{app.Items[0].Config.ConfigMap.Name: app.Items[0].Config.ConfigMap.Content},
	}
}

func NewConfigMapByService(svc *apiserver.Service, namespace string) *v1.ConfigMap {
	return &v1.ConfigMap{
		TypeMeta: resource.NewTypeMeta(resource.ResourceKindConfigMap, "v1"),
		ObjectMeta: metav1.ObjectMeta{
			Name:        svc.Name,
			Namespace:   namespace,
			Labels:      map[string]string{"name": svc.Name},
			Annotations: map[string]string{"name": svc.Name},
		},
		Data: map[string]string{svc.Config.ConfigMap.Name: svc.Config.ConfigMap.Content},
	}
}

func NewConfigMapByConfig(c *apiserver.Config) v1.ConfigMap {
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

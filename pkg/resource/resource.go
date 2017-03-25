package resouce

import (
	"apiserver/pkg/api/application"
	"apiserver/pkg/util/parseUtil"

	"k8s.io/client-go/pkg/api/resource"
	"k8s.io/client-go/pkg/api/v1"
	metav1 "k8s.io/client-go/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/util/intstr"
)

func newTypeMeta(kind, version string) metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       kind,
		APIVersion: version,
	}
}

func newOjectMeta(app *application.App) v1.ObjectMeta {
	return v1.ObjectMeta{
		Name:      app.Name,
		Namespace: app.UserName,
		Labels:    map[string]string{"name": app.Name},
	}
}

func NewPodSpec(app *application.App) v1.PodSpec {
	var containerPorts []v1.ContainerPort
	for _, port := range app.Ports {
		containerPorts = append(containerPorts, v1.ContainerPort{
			HostPort:      int32(port.TargetPort),
			ContainerPort: int32(port.TargetPort),
			Protocol:      v1.Protocol(port.Schame),
		})
	}
	return v1.PodSpec{
		RestartPolicy: v1.RestartPolicyOnFailure,
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
					v1.VolumeMount{
						Name:      app.Mount.Name,
						MountPath: app.Mount.MountPath,
						SubPath:   app.Mount.SubPath,
						ReadOnly:  app.Mount.ReadOnly,
					},
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

func newReplicationControllerSpec(app *application.App) v1.ReplicationControllerSpec {
	return v1.ReplicationControllerSpec{
		Replicas: parseUtil.Int32ToPointer(int32(app.InstanceCount)),
		Selector: map[string]string{"name": app.Name},
		Template: newPodTemplateSpec(app),
	}
}

func newServiceSpec(app *application.App) v1.ServiceSpec {
	var svcPorts []v1.ServicePort
	for _, port := range app.Ports {
		svcPorts = append(svcPorts, v1.ServicePort{
			Name:       app.Name,
			Port:       int32(port.ServicePort),
			TargetPort: intstr.FromInt(port.TargetPort),
			Protocol:   v1.Protocol(port.Schame),
		})
	}
	return v1.ServiceSpec{
		Selector: map[string]string{"name": app.Name},
		Ports:    svcPorts,
	}
}

func newNamespaceSpec(app *application.App) v1.NamespaceSpec {
	return v1.NamespaceSpec{
		Finalizers: []v1.FinalizerName{v1.FinalizerKubernetes},
	}
}
func NewSVC(app *application.App) v1.Service {
	return v1.Service{
		TypeMeta:   newTypeMeta("Service", "v1"),
		ObjectMeta: newOjectMeta(app),
		Spec:       newServiceSpec(app),
	}
}

func NewRC(app *application.App) v1.ReplicationController {
	return v1.ReplicationController{
		TypeMeta:   newTypeMeta("ReplicationController", "v1"),
		ObjectMeta: newOjectMeta(app),
		Spec:       newReplicationControllerSpec(app),
	}
}

func NewNS(app *application.App) v1.Namespace {
	return v1.Namespace{
		TypeMeta:   newTypeMeta("Namespace", "v1"),
		ObjectMeta: newOjectMeta(app),
		Spec:       newNamespaceSpec(app),
	}
}

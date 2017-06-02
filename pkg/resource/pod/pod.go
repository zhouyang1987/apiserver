package pod

import (
	"apiserver/pkg/api/apiserver"
	res "apiserver/pkg/resource"

	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/pkg/api/v1"
)

//newPodSpec create k8s's PodSpec
func NewPodSpec(svc *apiserver.Service) v1.PodSpec {
	var (
		containerPorts       []v1.ContainerPort
		envs                 []v1.EnvVar
		volumeMs             []v1.VolumeMount
		resourceRequirements v1.ResourceRequirements
	)
	if svc.Config != nil {
		ports := svc.Config.SuperConfig.Ports
		if len(ports) != 0 {
			for _, port := range ports {
				containerPorts = append(containerPorts, v1.ContainerPort{
					HostPort:      int32(port.ServicePort),
					ContainerPort: int32(port.ContainerPort),
					Protocol:      v1.Protocol(port.Protocol),
				})
			}
		}

		if len(svc.Config.Envs) != 0 {
			for _, env := range svc.Config.Envs {
				envs = append(envs, v1.EnvVar{Name: env.Key, Value: env.Val})
			}
		}

		if svc.Config.MapConfig != nil {
			volumeMs = append(volumeMs, v1.VolumeMount{Name: svc.Name, MountPath: svc.Config.MapConfig.ContainerPath, ReadOnly: false})
		}

		if svc.Config.BaseConfig != nil {
			resourceRequirements = v1.ResourceRequirements{
				Limits: v1.ResourceList{
					v1.ResourceCPU:    resource.MustParse(svc.Config.BaseConfig.Cpu),    //TODO 根据前端传入的值做资源限制
					v1.ResourceMemory: resource.MustParse(svc.Config.BaseConfig.Memory), //TODO 根据前端传入的值做资源限制
				},
				Requests: v1.ResourceList{
					v1.ResourceCPU:    resource.MustParse(svc.Config.BaseConfig.Cpu),
					v1.ResourceMemory: resource.MustParse(svc.Config.BaseConfig.Memory),
				},
			}
		}
	}
	return v1.PodSpec{
		Volumes:       []v1.Volume{},
		RestartPolicy: v1.RestartPolicyAlways,
		Containers: []v1.Container{
			v1.Container{
				Name:            svc.Name,
				Image:           svc.Image,
				ImagePullPolicy: v1.PullIfNotPresent,
				Resources:       resourceRequirements,
				Ports:           containerPorts,
				Env:             envs,
				VolumeMounts:    volumeMs,
			},
		},
	}
}

func NewPod(svc *apiserver.Service) *v1.Pod {
	return &v1.Pod{
		TypeMeta:   res.NewTypeMeta(res.ResourceKindPod, "v1"),
		ObjectMeta: res.NewObjectMeta(svc),
		Spec:       NewPodSpec(svc),
	}
}

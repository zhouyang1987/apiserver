package apiserver

import (
	"time"
)

type App struct {
	ID            uint       `json:"id"`
	CreatedAt     time.Time  `json:"createAt"`
	Name          string     `json:"nmae,omitempty"`
	UserName      string     `json:"nameSpace,omitempty"`
	Description   string     `json:"description,omitempty"`
	AppStatus     int        `json:"appStatus,omitempty"`
	ServiceCount  int        `json:"serviceCount,omitempty"`
	InstanceCount int        `json:"intanceCount,omitempty"`
	External      string     `json:"external,omitempty"`
	Items         []*Service `json:"services,omitempty"`
}

type Service struct {
	ID            uint           `json:"id"`
	CreatedAt     time.Time      `json:"createAt"`
	Name          string         `json:"name,omitempty"`
	Image         string         `json:"image,omitempty"`
	InstanceCount int            `json:"instanceCount,omitempty" `
	Status        int            `json:"status,omitempty"`
	External      string         `json:"external,omitempty"`
	LoadbalanceIp string         `json:"loadbalanceIp,omitempty"`
	Config        *ServiceConfig `json:"config,omitempty"`
	Items         []*Container   `json:"containers,omitempty"`
	AppId         uint           `json:"appId,omitempty"`
}

type ServiceConfig struct {
	ID          uint         `json:"id"`
	CreatedAt   time.Time    `json:"createAt"`
	BaseConfig  *BaseConfig  `json:"base,omitempty"`
	ConfigMap   *ConfigMap   `json:"config,omitempty"`
	SuperConfig *SuperConfig `json:"super,omitempty"`
	ServiceId   uint
}

type Container struct {
	ID        uint             `json:"id"`
	CreatedAt time.Time        `json:"createAt"`
	Name      string           `json:"name,omitempty"`
	Image     string           `json:"image,omitempty"`
	Status    int              `json:"status,omitempty"`
	Internal  string           `json:"internal,omitempty"`
	Config    *ContainerConfig `json:"config,omitempty"`
	ServiceId uint
}

type ContainerConfig struct {
	ID          uint         `json:"id"`
	CreatedAt   time.Time    `json:"createAt"`
	BaseConfig  *BaseConfig  `json:"base,omitempty"`
	ConfigMap   *ConfigMap   `json:"config,omitempty"`
	SuperConfig *SuperConfig `json:"super,omitempty"`
	ContainerId uint
}

type BaseConfig struct {
	ID              uint      `json:"id"`
	CreatedAt       time.Time `json:"createAt"`
	Cpu             string    `json:"cpu,omitempty"`
	Memory          string    `json:"memory,omitempty"`
	Type            int       `json:"type,omitempty"` //0 stateless 1 stateful
	Volumes         []*Volume `json:"volumes,omitempty"`
	ServiceConfigId uint
}

type Volume struct {
	ID           uint      `json:"id"`
	CreatedAt    time.Time `json:"createAt"`
	TargetPath   string    `json:"targetPath,omitempty"`
	Storage      string    `json:"storage,omitempty"`
	BaseConfigId uint
}

type ConfigMap struct {
	ID              uint      `json:"id"`
	CreatedAt       time.Time `json:"createAt"`
	Name            string    `json:"name,omitempty" `
	Content         string    `json:"content,omitempty"`
	ContainerPath   string    `json:"containerPath,omitempty"`
	ServiceConfigId uint
}

type SuperConfig struct {
	ID              uint      `json:"id"`
	CreatedAt       time.Time `json:"createAt"`
	Envs            []*Env    `json:"envs,omitempty"`
	Ports           []*Port   `json:"ports,omitempty"`
	ServiceConfigId uint
}

type Env struct {
	ID            uint      `json:"id"`
	CreatedAt     time.Time `json:"createAt"`
	Key           string    `json:"key,omitempty"`
	Val           string    `json:"val,omitempty"`
	SuperConfigId uint
}

type Port struct {
	ID            uint      `json:"id"`
	CreatedAt     time.Time `json:"createAt"`
	ContainerPort int       `json:"containerPort,omitempty"`
	ServicePort   int       `json:"servicePort,omitempty"`
	Protocol      string    `json:"protocol,omitempty"`
	SuperConfigId uint
}

type Logs struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createAt"`
	UserName  string    `json:",omitempty"`
	AppName   string    `json:",omitempty"`
	EventType string    `json:",omitempty"`
}

type ScaleOption struct {
	ServiceInstanceCnt int `json:"serviceInstanceCnt"`
}

type ExpansionOption struct {
	Cpu    string `json:"cpu"`
	Memory string `json:"memory"`
}

type RollOption struct {
	Image  string     `json:"image"`
	Conifg *ConfigMap `json:"config"`
}
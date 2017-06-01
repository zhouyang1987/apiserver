package app

import (
	"time"
)

type App struct {
	Name          string    `json:"nmae,omitempty" xorm:"pk index(app_index)"`
	UserName      string    `json:"nameSpace,omitempty" xorm:"pk index(app_index)"`
	CreateTime    time.Time `json:"create_time,omitempty" xorm:"created"`
	Description   string    `json:"description,omitempty" xorm:"varchar(255)"`
	AppStatus     int       `json:"appStatus,omitempty" xorm:"int(1) default 0"`
	ServiceCount  int       `json:"serviceCount,omitempty" xorm:"int not null"`
	InstanceCount int       `json:"intanceCount,omitempty" xorm:"int "`
	External      string    `json:"external,omitempty" xorm:"varchar(255)"`
	Items         []Service `json:"items,omitempty" xorm:"_"`
}

type Service struct {
	Name          string        `json:"name,omitempty"`
	AppName       string        `json:"appName,omitempty"`
	Image         string        `json:"image,omitempty"`
	Belong        string        `json:"belong,omitempty"`
	InstanceCount int           `json:"instanceCount,omitempty"`
	Status        int           `json:"status,omitempty"`
	CreateTime    time.Time     `json:"createTime,omitempty"`
	Items         []Container   `json:"items,omitempty"`
	Config        ServiceConfig `json:"serviceConfig,omitempty"`
}

type ServiceConfig struct {
	BaseConfig  `json:"base,omitempty"`
	MapConfig   `json:"config,omitempty"`
	SuperConfig `json:"super,omitempty"`
}

type Container struct {
	Name     string        `json:"name,omityempty"`
	Image    string        `json:"image,omityempty"`
	Belong   string        `json:"belong,omityempty"`
	Status   int           `json:"status,omityempty"`
	internal string        `json:"internal,omityempty"`
	Config   ServiceConfig `json:"conifg,omitempty"`
}

type BaseConfig struct {
	Cpu    string `json:"cpu,omityempty"`
	Memory string `json:"memory,omityempty"`
	//0 stateless 1 stateful
	Type    int      `json:"type,omityempty"`
	Volumes []Volume `json:"volumes,omityempty"`
}

type Volume struct {
	TargetPath string `json:"targetPath,omitempty"`
	Storage    string `json:"storage,omitempty"`
}

type MapConfig struct {
	ContainerPath string    `json:"containerPath,omtiempty"`
	ConfigMap     ConfigMap `json:"configMap,omtiempty"`
}

type ConfigMap struct {
	Name    string `json:"name,omtiempty"`
	Content string `json:"content,omtiempty"`
}

type SuperConfig struct {
	Envs  []Env  `json:"envs,omitempty"`
	Ports []Port `json:"ports,omitempty"`
}

type Env struct {
	Key string `json:"key,omitempty"`
	Val string `json:"val,omitempty"`
}

type Port struct {
	ContainerPort int    `json:"containerPort,omitempty"`
	ServicePort   int    `json:"servicePort,omitempty"`
	Protocol      string `json:"protocol,omitempty"`
}

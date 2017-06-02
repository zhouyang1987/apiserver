package apiserver

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
	Items         []Service `json:"services,omitempty" xorm:"_"`
}

type Service struct {
	Name          string         `json:"name,omitempty" xorm:"pk index(service_index)"`
	AppName       string         `json:"appName,omitempty" xorm:"pk index(service_index)"`
	UserName      string         `json:"nameSpace,omitempty" xorm:"pk index(app_index)"`
	Image         string         `json:"image,omitempty" xorm:"varchar(255) not null"`
	InstanceCount int            `json:"instanceCount,omitempty" xorm:"int"`
	Status        int            `json:"status,omitempty" xorm:"int(1)" default 0`
	CreateTime    time.Time      `json:"createTime,omitempty" xorm:"created"`
	Items         []Container    `json:"cantainers,omitempty" xorm:"_"`
	Config        *ServiceConfig `json:"config,omitempty" xorm:"_"`
	LoadbalanceIp string         `json:"loadbalanceIp,omitempty" xorm:"varchar(255)"`
}

type Container struct {
	Name        string         `json:"name,omitempty" xorm:"pk index(container_index)"`
	Image       string         `json:"image,omitempty" xorm:"varchar(255)"`
	ServiceName string         `json:"serviceName,omitempty" xorm:"pk index(container_index)"`
	Status      int            `json:"status,omitempty" xorm:"int(1) default 0"`
	internal    string         `json:"internal,omitempty" xorm:"varchar(255)"`
	Config      *ServiceConfig `json:"conifg,omitempty" xorm:"_"`
}

type ServiceConfig struct {
	*BaseConfig  `json:"base,omitempty"`
	*MapConfig   `json:"config,omitempty"`
	*SuperConfig `json:"super,omitempty"`
}

type BaseConfig struct {
	Cpu    string `json:"cpu,omitempty" xorm:""`
	Memory string `json:"memory,omitempty" xorm:""`
	//0 stateless 1 stateful
	Type    int      `json:"type,omitempty" xorm:""`
	Volumes []Volume `json:"volumes,omitempty" xorm:"_"`
}

type Volume struct {
	Id          int    `json:"id,omitempty" xorm:"pk autoincr"`
	ServiceName string `json:"serviceName,omitempty" xorm:"varchar(255)"`
	TargetPath  string `json:"targetPath,omitempty" xorm:"varchar(255)"`
	Storage     string `json:"storage,omitempty" xorm:"varchar(255)"`
}

type MapConfig struct {
	Id            int        `json:"id,omitempty" xorm:"pk autoincr"`
	ServiceName   string     `json:"serviceName,omitempty" xorm:"varchar(255)"`
	ContainerPath string     `json:"containerPath,omitempty" xorm:""`
	ConfigMap     *ConfigMap `json:"configMap,omitempty" xorm:"_"`
}

type ConfigMap struct {
	Id          int    `json:"id,omitempty" xorm:"pk autoincr"`
	ServiceName string `json:"serviceName,omitempty" xorm:""`
	Name        string `json:"name,omitempty" xorm:"pk"`
	Content     string `json:"content,omitempty" xorm:"varchar(2048)"`
}

type SuperConfig struct {
	Id          int    `json:"id,omitempty" xorm:"pk autoincr"`
	ServiceName string `json:"serviceName,omitempty" xorm:"varchar(255)"`
	Envs        []Env  `json:"envs,omitempty" xorm:"_"`
	Ports       []Port `json:"ports,omitempty" xorm:"_"`
}

type Env struct {
	Id          int    `json:"id,omitempty" xorm:"pk autoincr"`
	ServiceName string `json:"serviceName,omitempty" xorm:"varchar(255)"`
	Key         string `json:"key,omitempty" xorm:"varchar(255)"`
	Val         string `json:"val,omitempty" xorm:"varchar(1024)"`
}

type Port struct {
	Id            int    `json:"id,omitempty" xorm:"pk autoincr"`
	ServiceName   string `json:"serviceName,omitempty" xorm:"varchar(255)"`
	ContainerPort int    `json:"containerPort,omitempty" xorm:"int"`
	ServicePort   int    `json:"servicePort,omitempty" xorm:"int"`
	Protocol      string `json:"protocol,omitempty" xorm:"varchar(255)"`
}

type Logs struct {
	Id         int       `json:"id,omitempty" xorm:"pk autoincr"`
	UserName   string    `json:",omitempty" xorm:""`
	CreateTime time.Time `json:",omitempty" xorm:""`
	AppName    string    `json:",omitempty" xorm:""`
	EventType  string    `json:",omitempty" xorm:""`
}

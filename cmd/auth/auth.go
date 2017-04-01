package main

import (
	"encoding/json"
	"fmt"
)

type App struct {
	Name          string            `json:"name" xorm:"pk not null varchar(256)"`
	Region        string            `json:"region" xorm:"varchar(256)"`
	Memory        string            `json:"memory" xorm:"varchar(11)"`
	Cpu           string            `json:"cpu" xorm:"varchar(11)"`
	InstanceCount int               `json:"instanceCount" xorm:"int(11)"`
	Envs          map[string]string `json:"envs" xorm:"varchar(1024)"`
	// Ports         []Port            `json:"ports" xorm:"varchar(1024)"`
	Image   string   `json:"image" xorm:"varchar(1024)"`
	Command []string `json:"command" xorm:"varchar(1024)"`
	// Status        AppStatus         `json:"status" xorm:"int(1) default(0)"` //构建中 0 成功 1 失败 2 运行中 3 停止 4 删除 5
	UserName string `json:"userName" xorm:"varchar(256)"`
	Remark   string `json:"remark" xorm:"varchar(1024)"`
	// Mount         VolumeMount       `json:"mount" xorm:"varchar(1024)"`
	// Volume        []string          `json:"volume" xorm:"varchar(1024)"`
}

func main() {
	a := &App{}
	b, _ := json.Marshal(a)
	fmt.Println(string(b))
}

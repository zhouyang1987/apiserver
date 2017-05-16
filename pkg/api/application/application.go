// Copyright © 2017 huang jia <449264675@qq.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package application

import (
	"errors"
	"time"

	"apiserver/pkg/storage/mysqld"
	"apiserver/pkg/util/jsonx"
	"apiserver/pkg/util/log"
)

type AppStatus int32
type UpdateStatus int32

const (
	AppBuilding  AppStatus = 0
	AppSuccessed AppStatus = 1
	AppFailed    AppStatus = 2
	AppRunning   AppStatus = 3
	AppStop      AppStatus = 4
	AppDelete    AppStatus = 5
	AppUnknow    AppStatus = 6

	StartFailed    UpdateStatus = 10
	StartSuccessed UpdateStatus = 11

	StopFailed    UpdateStatus = 20
	StopSuccessed UpdateStatus = 21

	ScaleFailed    UpdateStatus = 30
	ScaleSuccessed UpdateStatus = 31

	UpdateConfigFailed    UpdateStatus = 40
	UpdateConfigSuccessed UpdateStatus = 41

	RedeploymentFailed    UpdateStatus = 50
	RedeploymentSuccessed UpdateStatus = 51
)

//App is struct of application
type App struct {
	Name          string            `json:"name" xorm:"pk not null varchar(255)"`
	Region        string            `json:"region" xorm:"varchar(255)"`
	Memory        string            `json:"memory" xorm:"varchar(11)"`
	Cpu           string            `json:"cpu" xorm:"varchar(11)"`
	InstanceCount int               `json:"instanceCount" xorm:"int(11)"`
	Envs          map[string]string `json:"envs" xorm:"varchar(255)"`
	Ports         []Port            `json:"ports" xorm:"varchar(255)"`
	Image         string            `json:"image" xorm:"varchar(255)"`
	Command       []string          `json:"command" xorm:"varchar(255)"`
	Status        AppStatus         `json:"status" xorm:"int(1) default(0)"` //构建中 0 成功 1 失败 2 运行中 3 停止 4 删除 5
	UserName      string            `json:"userName" xorm:"varchar(256)"`
	Remark        string            `json:"remark" xorm:"varchar(1024)"`
	CreateAt      time.Time         `json:"create_at" xorm:"created"`
	ReviseAt      time.Time         `json:"revise_at" xorm:"updated"`
	Url           string            `json:"url" xorm:"varchar(1024)"`
	// Mount         VolumeMount       `json:"mount" xorm:"varchar(1024)"`
	// Volume        []string          `json:"volume" xorm:"varchar(1024)"`
}

type VolumeMount struct {
	// This must match the Name of a Volume.
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// Mounted read-only if true, read-write otherwise (false or unspecified).
	// Defaults to false.
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,2,opt,name=readOnly"`
	// Path within the container at which the volume should be mounted.  Must
	// not contain ':'.
	MountPath string `json:"mountPath" protobuf:"bytes,3,opt,name=mountPath"`
	// Path within the volume from which the container's volume should be mounted.
	// Defaults to "" (volume's root).
	// +optional
	SubPath string `json:"subPath,omitempty" protobuf:"bytes,4,opt,name=subPath"`
}

type Port struct {
	Schame      string
	ServicePort int
	TargetPort  int
}

var (
	engine = mysqld.GetEngine()
	Status = map[AppStatus]string{
		AppBuilding:  "AppBuilding",
		AppSuccessed: "AppSuccessed",
		AppFailed:    "AppFailed",
		AppRunning:   "AppRunning",
		AppStop:      "AppStop",
		AppDelete:    "AppDelete",
		AppUnknow:    "AppUnknow",
	}
)

func init() {
	if err := engine.Sync(new(App)); err != nil {
		log.Fatalf("Sync fail :%s", err.Error())
	}
}

func (app *App) String() string {
	appStr := jsonx.ToJson(app)
	return appStr
}

func (app *App) Insert() error {
	_, err := engine.Insert(app)
	if err != nil {
		return err
	}
	return nil
}

func (app *App) Delete() error {
	_, err := engine.Id(app.Name).Delete(app)
	if err != nil {
		return err
	}
	return nil
}

func (app *App) Update() error {
	_, err := engine.Id(app.Name).Update(app)
	if err != nil {
		return err
	}
	return nil
}

func (app *App) QueryOne() (*App, error) {
	has, err := engine.Id(app.Name).Get(app)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("current app not exsit")
	}
	return app, nil
}

func (app *App) QuerySet() ([]*App, error) {
	appSet := []*App{}
	err := engine.Where("1 and 1 order by name desc").Find(&appSet)
	if err != nil {
		return nil, err
	}
	return appSet, nil
}

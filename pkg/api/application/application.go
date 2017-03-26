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
	"net/http"

	"apiserver/pkg/storage/mysqld"
	"apiserver/pkg/util/jsonx"
	"apiserver/pkg/util/log"

	"github.com/emicklei/go-restful"
)

func Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/app").
		Doc("manage application").
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.POST("").To(new(App).Insert).
		Doc("create application").
		Produces("CreateApplication").
		Reads(App{}))

	ws.Route(ws.GET("/{app-id}").To(new(App).QueryOne).
		// docs
		Doc("get a app").
		Operation("findUser").
		Param(ws.PathParameter("app-id", "identifier of the app").DataType("int")).
		Writes(App{})) // on the response

	container.Add(ws)
}

type AppStatus int32
type UpdateStatus int32

const (
	AppBuilding  AppStatus = 0
	AppSuccessed AppStatus = 1
	AppFailed    AppStatus = 2
	AppRunning   AppStatus = 3
	AppStop      AppStatus = 4

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
	Id            int               `json:"id" xorm:"pk not null autoincr int(11)"`
	Name          string            `json:"name" xorm:"varchar(256)"`
	Region        string            `json:"region" xorm:"varchar(256)"`
	Memory        string            `json:"memory" xorm:"varchar(11)"`
	Cpu           string            `json:"cpu" xorm:"varchar(11)"`
	InstanceCount int32             `json:"instanceCount" xorm:"int(11)"`
	Envs          map[string]string `json:"envs" xorm:"varchar(1024)"`
	Ports         []Port            `json:"ports" xorm:"varchar(1024)"`
	Image         string            `json:"image" xorm:""`
	Command       []string          `json:"command" xorm:""`
	Status        AppStatus         `json:"status" xorm:"int(1)"` //构建中 0 成功 1 失败 2 运行中 3 停止 4
	UserName      string            `json:"userName" xorm:"varchar(256)"`
	Remark        string            `json:"remark" xorm:"varchar(1024)"`
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
)

func init() {
	engine.ShowSQL(true)
	if err := engine.Sync(new(App)); err != nil {
		log.Fatalf("Sync fail :%s", err.Error())
	}
}

func (app *App) String() string {
	appStr, err := jsonx.ToJson(app)
	if err != nil {
		log.Errorf("node to string err :%s", err.Error())
		return ""
	}
	return appStr
}

func (app *App) Insert(request *restful.Request, response *restful.Response) {
	aps := new(App)
	err := request.ReadEntity(aps)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	_, err = engine.Insert(aps)
	if err != nil {
		return
	}
	response.WriteHeaderAndEntity(http.StatusCreated, aps)
}

func (app *App) Delete() error {
	_, err := engine.Id(app.Id).Delete(app)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) Update() error {
	_, err := engine.Id(app.Id).Update(app)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) QueryOne(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("app-id")

	log.Debug(id)

	aps := &App{}
	// engine.Id(id).Get(&aps)
	_, err := engine.Id(1).Get(aps)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, err.Error())
		return
	}
	// if !has {
	// 	response.AddHeader("Content-Type", "text/plain")
	// 	response.WriteErrorString(http.StatusNotFound, "404: User could not be found.")
	// 	return
	// }

	log.Debugf("%#v", aps)
	response.WriteEntity(aps)
}

func (app *App) QuerySet() ([]*App, error) {
	appSet := []*App{}
	err := engine.Where("1 and 1 order by id desc").Find(&appSet)
	if err != nil {
		return nil, err
	}
	return appSet, nil
}

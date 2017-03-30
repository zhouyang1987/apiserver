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

package app

import (
	"encoding/json"
	"net/http"
	// "strconv"

	"apiserver/pkg/api/application"
	"apiserver/pkg/resource"
	"apiserver/pkg/resource/sync"
	r "apiserver/pkg/router"
	"apiserver/pkg/util/log"
	// "apiserver/pkg/util/parseUtil"

	// res "k8s.io/client-go/pkg/api/resource"
	// "k8s.io/client-go/pkg/api/v1"

	"github.com/gorilla/mux"
)

func Register(rout *mux.Router) {
	r.RegisterHttpHandler(rout, "/app", "POST", CreateApplication)
	r.RegisterHttpHandler(rout, "/app", "DELETE", DeleteApplication)
}

func CreateApplication(request *http.Request) (string, interface{}) {
	//get the request's body ,then marsh to app struct
	decoder := json.NewDecoder(request.Body)
	app := &application.App{}
	err := decoder.Decode(app)
	if err != nil {
		log.Errorf("decode the request body err:%v", err)
		return r.StatusBadRequest, "json format error"
	}

	//create namespace, first query the ns is exsit or not, if not exsit, create it
	ns := resource.NewNS(app)
	if !resource.ExsitResource(ns) {
		err := resource.CreateResource(ns)
		if err != nil {
			return r.StatusInternalServerError, err
		}
	}

	//create service, first query the svc is exsit or not, if not exsit, create it
	svc := resource.NewSVC(app)
	if !resource.ExsitResource(svc) { //if service not exsit,then create service
		if err = app.Insert(); err != nil {
			return r.StatusInternalServerError, "access database err:" + err.Error()
		}
		err := resource.CreateResource(svc)
		if err != nil {
			return r.StatusInternalServerError, err
		}
	} else {
		//if service exsited, the application already exsit,so return and tell the app already exsit
		return r.StatusForbidden, "the application of named " + app.Name + " is already exsit"
	}

	//create replicationControllers, first query the svc is exsit or not, if not exsit, create it
	rc := resource.NewRC(app)
	if !resource.ExsitResource(rc) {
		err := resource.CreateResource(rc)
		if err != nil {
			return r.StatusInternalServerError, err
		}
	} else {
		return r.StatusForbidden, "the application of named " + app.Name + " is already exsit"
	}
	go resource.WatchPodStatus(app)
	return r.StatusCreated, "create app successed"
}

func DeleteApplication(request *http.Request) (string, interface{}) {
	appName := request.FormValue("app_name")
	namespace := request.FormValue("ns")
	ns, exsit := sync.ListReplicationController[namespace][appName]
	if !exsit {
		return r.StatusNotFound, "application named " + appName + ` does't exist`
	}
	if err := resource.DeleteResource(&ns); err != nil {
		return r.StatusInternalServerError, "delete application err: " + err.Error()
	}
	svc := sync.ListService[namespace][appName]
	log.Debug(svc)
	if &svc == nil {
		return r.StatusNotFound, "application named " + appName + `does't exist`
	}
	if err := resource.DeleteResource(&svc); err != nil {
		return r.StatusInternalServerError, "delete application err: " + err.Error()
	}
	app := &application.App{Name: appName}
	if err := app.Delete(); err != nil {
		return r.StatusInternalServerError, "the application is delete,but update the application's status err: " + err.Error()
	}
	return r.StatusNoContent, "delete app successed"
}

// func StopApplication(request *http.Request) (string, interface{}) {
// 	appName := request.FormValue("app_name")
// 	rc := sync.ListReplicationController[appName]
// 	if &rc == nil {
// 		return r.StatusNotFound, "application named " + appName + `does't exist`
// 	}
// 	rc.Spec.Replicas = parseUtil.IntToInt32Pointer(0)
// 	if err := resource.UpdateResouce(&rc); err != nil {
// 		return r.StatusInternalServerError, "stop application named " + appName + " failed"
// 	}
// 	app := &application.App{Name: appName, Status: application.AppStop}
// 	if err := app.Update(); err != nil {
// 		return r.StatusInternalServerError, "stop application named " + appName + " failed"
// 	}
// 	return r.StatusCreated, "stop application named " + appName + " successed"
// }

// func StartApplication(request *http.Request) (string, interface{}) {
// 	appName := request.FormValue("app_name")
// 	rc := sync.ListReplicationController[appName]
// 	if &rc == nil {
// 		return r.StatusNotFound, "application named " + appName + `does't exist`
// 	}
// 	app := &application.App{Name: appName}
// 	temApp, err := app.QueryOne()
// 	if err != nil {
// 		return r.StatusInternalServerError, "get application named " + appName + ` ` + err.Error()
// 	}
// 	rc.Spec.Replicas = parseUtil.IntToInt32Pointer(temApp.InstanceCount)
// 	if err := resource.UpdateResouce(&rc); err != nil {
// 		return r.StatusInternalServerError, "start application named " + appName + ` failed`
// 	}
// 	go resource.WatchPodStatus(app)
// 	return r.StatusCreated, "start application named " + appName + " successed"
// }

// func ScaleApplication(request *http.Request) (string, interface{}) {
// 	appName := request.FormValue("app_name")
// 	app_cnt := request.FormValue("app_cnt")
// 	rc := sync.ListReplicationController[appName]
// 	if &rc == nil {
// 		return r.StatusNotFound, "application named " + appName + `does't exist`
// 	}
// 	cnt, _ := strconv.Atoi(app_cnt)
// 	rc.Spec.Replicas = parseUtil.IntToInt32Pointer(cnt)
// 	if err := resource.UpdateResouce(&rc); err != nil {
// 		return r.StatusInternalServerError, "scale application named " + appName + ` failed`
// 	}
// 	app := &application.App{Name: appName, InstanceCount: cnt}
// 	if err := app.Update(); err != nil {
// 		return r.StatusInternalServerError, "update application named " + appName + ` failed`
// 	}
// 	return r.StatusCreated, "scale application named " + appName + ` successed`
// }

// /*func RollingUpdateApplication(request *http.Request) (string, interface{}) {
// 	appName := request.FormValue("app_name")
// 	image := request.FormValue("image")
// 	// period := request.FormValue("period")
// 	rc := sync.ListReplicationController[appName]
// 	if &rc == nil {
// 		return r.StatusNotFound, "application named " + appName + `does't exist`
// 	}
// 	return r.StatusCreated, "rolling update application named " + appName + ` successed`
// }*/

// func ReDeployApplication(request *http.Request) (string, interface{}) {
// 	appName := request.FormValue("app_name")
// 	rc := sync.ListReplicationController[appName]
// 	if &rc == nil {
// 		return r.StatusNotFound, "application named " + appName + `does't exist`
// 	}
// 	if err := resource.DeleteResource(&rc); err != nil {
// 		return r.StatusInternalServerError, "redploy application named " + appName + " failed"
// 	}
// 	return r.StatusCreated, "redeploy application named " + appName + " successed"
// }

// func ExpansionApplication(request *http.Request) (string, interface{}) {
// 	appName := request.FormValue("app_name")
// 	cpu := request.FormValue("cpu")
// 	memory := request.FormValue("memory")
// 	rc := sync.ListReplicationController[appName]
// 	if &rc == nil {
// 		return r.StatusNotFound, "application named " + appName + `does't exist`
// 	}
// 	rc.Spec.Template.Spec.Containers[0].Resources = v1.ResourceRequirements{
// 		Limits: v1.ResourceList{
// 			v1.ResourceCPU:    res.MustParse(cpu),    //TODO 根据前端传入的值做资源限制
// 			v1.ResourceMemory: res.MustParse(memory), //TODO 根据前端传入的值做资源限制
// 		},
// 		Requests: v1.ResourceList{
// 			v1.ResourceCPU:    res.MustParse(cpu),
// 			v1.ResourceMemory: res.MustParse(memory),
// 		},
// 	}
// 	if err := resource.UpdateResouce(&rc); err != nil {
// 		return r.StatusInternalServerError, "redeploy application named " + appName + " failed"
// 	}
// 	app := &application.App{Name: appName, Cpu: cpu, Memory: memory}
// 	if err := app.Update(); err != nil {
// 		return r.StatusInternalServerError, "redeploy application named " + appName + " failed:" + err.Error()
// 	}
// 	return r.StatusCreated, "redeploy application named " + appName + " successed"
// }

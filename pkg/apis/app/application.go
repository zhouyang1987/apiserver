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

	"apiserver/pkg/api/application"
	"apiserver/pkg/resource"
	"apiserver/pkg/resource/sync"
	r "apiserver/pkg/router"
	"apiserver/pkg/util/log"

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
	if err = app.Insert(); err != nil {
		return r.StatusInternalServerError, "access database err:" + err.Error()
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
	//TODO 掉用k8s的pkg下的方法去获取svc ns rc的状态
	//当ns，svc，rc都创建成功后，进行本地数据库的数据插入操作
	return r.StatusCreated, "create app successed"
}

func DeleteApplication(request *http.Request) (string, interface{}) {
	appName := request.FormValue("app_name")
	log.Debugf("appname=%s", appName)
	svc := sync.ListService[appName]
	if err := resource.DeleteResource(&svc); err != nil {
		return r.StatusInternalServerError, "delete application err: " + err.Error()
	}
	rc := sync.ListReplicationController[appName]
	if err := resource.DeleteResource(&rc); err != nil {
		return r.StatusInternalServerError, "delete application err: " + err.Error()
	}
	app := &application.App{Name: appName, Status: 5}
	if err := app.Update(); err != nil {
		return r.StatusInternalServerError, "the application is delete,but update the application's status err: " + err.Error()
	}
	go resource.WatchPodStatus(app)
	return r.StatusNoContent, "delete app successed"
}

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
	r "apiserver/pkg/router"
	"apiserver/pkg/util/log"

	"github.com/gorilla/mux"
)

func Register(rout *mux.Router) {
	r.RegisterHttpHandler(rout, "/app", "POST", CreateApplication)
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
	ns := resouce.NewNS(app)
	if !resouce.ExsitResource(ns) {
		err := resouce.CreateResource(ns)
		if err != nil {
			return r.StatusInternalServerError, err
		}
	}

	//create service, first query the svc is exsit or not, if not exsit, create it
	svc := resouce.NewSVC(app)
	if !resouce.ExsitResource(svc) { //if service not exsit,then create service
		err := resouce.CreateResource(svc)
		if err != nil {
			return r.StatusInternalServerError, err
		}
	} else {
		//if service exsited, the application already exsit,so return and tell the app already exsit
		return r.StatusForbidden, "the application of named " + app.Name + " is already exsit"
	}

	//create replicationControllers, first query the svc is exsit or not, if not exsit, create it
	rc := resouce.NewRC(app)
	if !resouce.ExsitResource(rc) {
		err := resouce.CreateResource(rc)
		if err != nil {
			return r.StatusInternalServerError, err
		}
	} else {
		return r.StatusForbidden, "the application of named " + app.Name + " is already exsit"
	}

	//TODO 掉用k8s的pkg下的方法去获取svc ns rc的状态
	//当ns，svc，rc都创建成功后，进行本地数据库的数据插入操作
	return r.StatusCreated, "create app successed"
}

func DeleteApplication(request *http.Request) (string, interface{}) {
	// request.FormValue("id")
	decoder := json.NewDecoder(request.Body)
	app := &application.App{}
	err := decoder.Decode(app)
	if err != nil {
		log.Errorf("decode the request body err:%v", err)
		return r.StatusBadRequest, "json format error"
	}
}

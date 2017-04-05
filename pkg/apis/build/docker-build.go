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

package build

import (
	"encoding/json"
	"net/http"

	"apiserver/pkg/api/build"
	r "apiserver/pkg/router"
	"apiserver/pkg/util/log"

	"github.com/gorilla/mux"
)

func Register(rout *mux.Router) {
	r.RegisterHttpHandler(rout, "/build", "POST", OnlineBuild)
	r.RegisterHttpHandler(rout, "/build", "PUT", OfflineBuild)
}

//OnlineBuild build application online
func OnlineBuild(request *http.Request) (string, interface{}) {
	decoder := json.NewDecoder(request.Body)
	builder := &build.Build{}
	err := decoder.Decode(builder)
	if err != nil {
		log.Errorf("decode the request body err:%v", err)
		return r.StatusBadRequest, "json format error"
	}

	return r.StatusCreated, nil
}

//OfflineBuild build application offline
func OfflineBuild(request *http.Request) (string, interface{}) {
	decoder := json.NewDecoder(request.Body)
	builder := &build.Build{}
	err := decoder.Decode(builder)
	if err != nil {
		log.Errorf("decode the request body err:%v", err)
		return r.StatusBadRequest, "json format error"
	}

	return r.StatusCreated, nil
}

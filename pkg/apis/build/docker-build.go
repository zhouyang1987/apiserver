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
	"fmt"
	"net/http"
	"os"

	"apiserver/pkg/api/build"
	"apiserver/pkg/client"
	r "apiserver/pkg/router"
	"apiserver/pkg/util/log"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

const (
	DEFAULT_REGISTRY = `hub.mini-paas.com`
	TARBALL_ROOT_DIR = `/tmp`
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

	//build image's step:
	//1. find the upload tarball file of user uploaded in TARBALL_ROOT_DIR/userid/ dir
	//2. find the Dockerfile file of the user upload in TARBALL_ROOT_DIR/userid/ dir
	//3. if the Dockerfile exsit,we will use the tarball and the Dockerfile to create a tar file's stream in order to build image.
	//4. if the Dokcerfile didn't exsit,will be use the Dockerfile template to create Dockerfile by the project's type. for exampler:if
	//the project's language is golang, we will use the golang's Dockerfile template,and then we will use the tarball and the Dockerfile
	//to create a tar file's stream in order to build image.

	Dockerfile := fmt.Sprintf("%s/%s/%s", TARBALL_ROOT_DIR, builder.UserId, "Dockerfile")
	fileInfo, err := os.Stat(Dockerfile)
	if err != nil {
		log.Errorf("Dockerfile doesn't exsit: %v", err)
		return r.StatusInternalServerError, "build image fail"
	}

	tarball := fmt.Sprintf("%s/%s/%s", TARBALL_ROOT_DIR, builder.UserId, builder.Tarball)
	buildContext, err := os.Open(tarball)
	if err != nil {
		log.Fatal(err)
	}

	image_repo := fmt.Sprintf("%s/%s:%s", DEFAULT_REGISTRY, builder.AppName, builder.Version)
	options := types.ImageBuildOptions{
		Tags:       []string{image_repo},
		Dockerfile: "Dockerfile",
	}
	buildResponse, err := cli.ImageBuild(context.Background(), buildContext, options)
	if err != nil {

	}
	res, err := ioutil.ReadAll(buildResponse.Body)
	if err != nil {
		log.Errorf("read the build image response err: %v", err)
		return r.StatusInternalServerError, err.Error()
	}

	return r.StatusCreated, string(res)
}

func BuildProject(url string) {

}

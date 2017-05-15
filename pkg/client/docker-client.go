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

package client

import (
	"net/http"

	"apiserver/pkg/configz"
	"apiserver/pkg/util/log"

	"github.com/docker/docker/client"
)

var (
	DockerClient *client.Client
)

func init() {
	host := configz.GetString("build", "endpoint", "127.0.0.1:2375")
	version := configz.GetString("build", "version", "1.24")
	cl := &http.Client{
		Transport: new(http.Transport),
	}
	DockerClient, err = client.NewClient(host, version, cl, nil)
	if err != nil {
		log.Fatalf("init docker client err: %v", err)
	}
}

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

package config

/*
import (
	"encoding/json"
	"io/ioutil"

	"apiserver/pkg/util/logger"
)

type Config struct {
	Driver     string `json:"driver"`
	Dsn        string `json:"dsn"`
	Server     string `json:"server"`
	K8sServer  string `json:"k8sserver"`
	Kubeconfig string `json:"kubeconfig"`
}

var (
	log         = logger.New("")
	GloabConfig = &Config{}
)

func Parse(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("read config file fail, the reason is %v ", err)
	}
	err = json.Unmarshal(data, GloabConfig)
	if err != nil {
		log.Fatalf("unmarshal config data to config struct fail, the reason is %v ", err)
	}
}
*/

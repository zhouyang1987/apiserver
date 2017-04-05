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
	"apiserver/pkg/configz"
	"apiserver/pkg/util/log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	K8sClient *kubernetes.Clientset
	err       error
)

//init create client of k8s's apiserver
func init() {
	// uses the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", configz.GetString("apiserver", "k8s-config", "./config"))
	if err != nil {
		log.Fatalf("init k8s config err: %v", err)
	}
	// creates the clientset
	K8sClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("init k8s client err: %v", err)
	}

}

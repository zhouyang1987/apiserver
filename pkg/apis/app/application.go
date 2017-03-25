package app

import (
	"encoding/json"
	"net/http"

	"apiserver/pkg/api/application"
	"apiserver/pkg/client"
	"apiserver/pkg/configz"
	"apiserver/pkg/resource"
	r "apiserver/pkg/router"
	"apiserver/pkg/util/log"

	"github.com/gorilla/mux"
)

var (
	k8sClient = client.NewK8sClient(configz.GetString("apiserver", "k8s-config", "./config"))
)

func Register(router *mux.Route) {
	r.RegisterHttpHandler(router, "/app", "POST", CreateApplication)
}

func CreateApplication(request *http.Request) {
	//获取请求数据，解析成app对象
	decoder := json.NewDecoder(req.Body)
	app := &application.App{}
	err := decoder.Decode(app)
	if err != nil {
		//TODO
	}
	ns := resouce.NewNS(app)
	svc := resouce.NewSVC(app)
	rc := resouce.NewRC(app)

	namespace, err := k8sClient.CoreV1().Namespaces().Create(ns)
	if err != nil {
		//TODO
	}
	service, err := k8sClient.CoreV1().Services(ns.Name).Create(svc)
	if err != nil {
		//TODO
	}
	replication, err := k8sClient.CoreV1().ReplicationControllers(ns.Name).Create(rc)
	if err != nil {
		//TODO
	}

	//当ns，svc，rc都创建成功后，进行本地数据库的数据插入操作
}

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
	k8sClient, _ = client.NewK8sClient(configz.GetString("apiserver", "k8s-config", "./config"))
)

func Register(rout *mux.Router) {
	r.RegisterHttpHandler(rout, "/app", "POST", CreateApplication)
}

func CreateApplication(request *http.Request) (string, interface{}) {
	//获取请求数据，解析成app对象
	decoder := json.NewDecoder(request.Body)
	app := &application.App{}
	err := decoder.Decode(app)
	if err != nil {
		//TODO
		log.Error(err)
	}
	ns := resouce.NewNS(app)
	svc := resouce.NewSVC(app)
	rc := resouce.NewRC(app)

	namespace, err := k8sClient.CoreV1().Namespaces().Create(ns)
	if err != nil {
		//TODO
		log.Error(err)
	}
	log.Info(namespace)
	service, err := k8sClient.CoreV1().Services(ns.Name).Create(svc)
	if err != nil {
		//TODO
		log.Error(err)
	}
	log.Info(service)

	replication, err := k8sClient.CoreV1().ReplicationControllers(ns.Name).Create(rc)
	if err != nil {
		//TODO
		log.Error(err)
	}
	log.Info(replication)

	//TODO 掉用k8s的pkg下的方法去获取svc ns rc的状态
	//当ns，svc，rc都创建成功后，进行本地数据库的数据插入操作
	return "", nil
}

package app

import (
	"fmt"
	"net/http"

	a "apiserver/pkg/apis/app"
	"apiserver/pkg/componentconfig"
	"apiserver/pkg/configz"
	"apiserver/pkg/util/log"

	"github.com/gorilla/mux"
	// "github.com/emicklei/go-restful"
)

type Apiserver struct {
	*componentconfig.ApiserverConfig
}

func NewApiServer() *Apiserver {
	return &Apiserver{
		ApiserverConfig: &componentconfig.ApiserverConfig{
			HttpAddr: configz.GetString("apiserver", "httpAddr", "0.0.0.0"),
			HttpPort: configz.MustInt("apiserver", "httpPort", 9090),
			RpcAddr:  configz.GetString("apiserver", "rpcAddr", "0.0.0.0"),
			RpcPort:  configz.MustInt("apiserver", "rpcPort", 7070),
		},
	}
}

/*func Run(server *Apiserver) error {
	wsContainer := restful.NewContainer()
	application.Register(wsContainer)
	log.Infof("starting apiserver and listen on : %v", fmt.Sprintf("%v:%v", server.HttpAddr, server.HttpPort))
	return http.ListenAndServe(fmt.Sprintf("%v:%v", server.HttpAddr, server.HttpPort), wsContainer)
}*/

func Run(server *Apiserver) error {
	root := mux.NewRouter()
	api := root.PathPrefix("/api/v1").Subrouter()
	installApiGroup(api)
	http.Handle("/", root)
	log.Infof("starting apiserver and listen on : %v", fmt.Sprintf("%v:%v", server.HttpAddr, server.HttpPort))
	return http.ListenAndServe(fmt.Sprintf("%v:%v", server.HttpAddr, server.HttpPort), nil)
}

func installApiGroup(router *mux.Router) {
	a.Register(router)
}

// func installNodeApi(router *mux.Router) {

// }

// func installAppApi(router *mux.Router) {

// }

// func installContainerApi(router *mux.Router) {

// }

// func installDeploymentApi(router *mux.Router) {

// }

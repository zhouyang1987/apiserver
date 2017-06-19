package apiserver

import (
	"net/http"

	r "apiserver/pkg/router"

	"github.com/gorilla/mux"
)

func InstallApi(router *mux.Router) {
	//install app's api handle
	r.RegisterHttpHandler(router, "/{namespace}/apps", "GET", GetApps)
	r.RegisterHttpHandler(router, "/{namespace}/apps/{id}", "GET", GetApp)
	r.RegisterHttpHandler(router, "/{namespace}/apps", "POST", CreateApp)
	r.RegisterHttpHandler(router, "/{namespace}/apps/{id}", "DELETE", DeleteApp)
	r.RegisterHttpHandler(router, "/{namespace}/apps/{id}/{verb}", "PATCH", StopOrStartOrRedeployApp)
	r.RegisterHttpHandler(router, "/{namespace}/apps/{id}/{verb}", "OPTIONS", Option)
	r.RegisterHttpHandler(router, "/{namespace}/apps", "OPTIONS", Option)
	r.RegisterHttpHandler(router, "/{namespace}/apps/", "OPTIONS", Option)
	r.RegisterHttpHandler(router, "/{namespace}/apps/{id}/", "OPTIONS", Option)
	r.RegisterHttpHandler(router, "/{namespace}/apps/{id}", "OPTIONS", Option)

	//install service's api handle
	r.RegisterHttpHandler(router, "/{namespace}/services", "GET", GetServices)
	r.RegisterHttpHandler(router, "/{namespace}/services/{id}", "GET", GetService)
	r.RegisterHttpHandler(router, "/{namespace}/services/{name}/events", "GET", GetServiceEvents)
	r.RegisterHttpHandler(router, "/{namespace}/services", "POST", CreateService)
	r.RegisterHttpHandler(router, "/{namespace}/services/{id}", "DELETE", DeleteService)
	r.RegisterHttpHandler(router, "/{namespace}/services/{id}/{verb}", "PUT", UpdateServiceConfig)
	r.RegisterHttpHandler(router, "/{namespace}/services/{id}/{verb}", "PATCH", StopOrStartOrRedployService)
	r.RegisterHttpHandler(router, "/{namespace}/services", "OPTIONS", Option)
	r.RegisterHttpHandler(router, "/{namespace}/services/{id}", "OPTIONS", Option)
	r.RegisterHttpHandler(router, "/{namespace}/services/{id}/", "OPTIONS", Option)
	r.RegisterHttpHandler(router, "/{namespace}/services/{id}/{verb}", "OPTIONS", Option)

	//install container's api handle
	r.RegisterHttpHandler(router, "/{namespace}/containers", "GET", GetContainers)
	r.RegisterHttpHandler(router, "/{namespace}/containers/{name}/events", "GET", GetContainerEvents)
	r.RegisterHttpHandler(router, "/{namespace}/containers/{id}", "PATCH", RedeployContainer)
	r.RegisterHttpHandler(router, "/{namespace}/containers", "OPTIONS", Option)
	r.RegisterHttpHandler(router, "/{namespace}/containers/", "OPTIONS", Option)
	r.RegisterHttpHandler(router, "/{namespace}/containers/{name}/events", "OPTIONS", Option)

	//install metrics's api handle
	r.RegisterHttpHandler(router, "/{namespace}/metrics/{name}/{metric}/{type}", "OPTIONS", Option)
	r.RegisterHttpHandler(router, "/{namespace}/metrics/{name}/{metric}/{type}", "GET", GetMetrics)
	// r.RegisterHttpHandler(router, "/{namespace}/pod/", "OPTIONS", Option)

	//install configMap's api handle
	r.RegisterHttpHandler(router, "/{namespace}/configs", "GET", GetConfigs)
	r.RegisterHttpHandler(router, "/{namespace}/configs", "POST", CreateConfig)
	r.RegisterHttpHandler(router, "/{namespace}/configs/{id}", "DELETE", DeleteConfig)
	r.RegisterHttpHandler(router, "/{namespace}/configs/{id}/items", "POST", CreateConfigItem)
	r.RegisterHttpHandler(router, "/{namespace}/configs/{id}/items/{itemId}", "DELETE", DeleteConfigItem)
	r.RegisterHttpHandler(router, "/{namespace}/configs", "OPTIONS", Option)
	r.RegisterHttpHandler(router, "/{namespace}/configs/{id}", "OPTIONS", Option)
	r.RegisterHttpHandler(router, "/{namespace}/configs/{id}/items", "OPTIONS", Option)
	r.RegisterHttpHandler(router, "/{namespace}/configs/{id}/items/{itemId}", "OPTIONS", Option)

	//install deploy's api handle
	r.RegisterHttpHandler(router, "/{namespace}/deploys", "POST", CreatDeploy)
	r.RegisterHttpHandler(router, "/{namespace}/deploys", "OPTIONS", Option)
}

func Option(request *http.Request) (string, interface{}) {
	return r.StatusOK, "ok"
}

package apiserver

import (
	"net/http"

	"apiserver/pkg/configz"
	r "apiserver/pkg/router"
)

func Health(request *http.Request) (string, interface{}) {
	return r.StatusOK, "apiserver is healthy"
}

func GetApiserverInfo(request *http.Request) (string, interface{}) {
	return r.StatusOK, map[string]interface{}{"componentName": configz.GetString("apiserver", "componentName", "apiserver"), "version": configz.GetString("apiserver", "version", "v1.0")}
}

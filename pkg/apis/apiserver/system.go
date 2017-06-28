package apiserver

import (
	"net/http"

	"apiserver/pkg/api/apiserver"
	"apiserver/pkg/configz"
	r "apiserver/pkg/router"
)

func Health(request *http.Request) (string, interface{}) {
	return r.StatusOK, "apiserver is healthy"
}

func GetApiserverInfo(request *http.Request) (string, interface{}) {
	return r.StatusOK, map[string]interface{}{"componentName": configz.GetString("registry", "componentName", "registry"), "version": configz.GetString("registry", "version", "v1.0")}
}

func GetAppCount(request *http.Request) (string, interface{}) {
	result, err := apiserver.CountApp()
	if err != nil {
		return r.StatusInternalServerError, err
	}
	return r.StatusOK, result
}

func GetServiceCount(request *http.Request) (string, interface{}) {
	result, err := apiserver.CountService()
	if err != nil {
		return r.StatusInternalServerError, err
	}
	return r.StatusOK, result
}

func GetContainerCount(request *http.Request) (string, interface{}) {
	result, err := apiserver.CountContainer()
	if err != nil {
		return r.StatusInternalServerError, err
	}
	return r.StatusOK, result
}

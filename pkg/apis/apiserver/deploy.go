package apiserver

import (
	"encoding/json"
	"net/http"

	"apiserver/pkg/api/apiserver"
	r "apiserver/pkg/router"
)

func CreatDeploy(request *http.Request) (string, interface{}) {
	deploys, err := validateDeploy(request)
	if err != nil {
		return r.StatusInternalServerError, err
	}

	for _, deploy := range deploys {
		if err = apiserver.InsertDeploy(deploy); err != nil {
			return r.StatusInternalServerError, err
		}
	}
	return r.StatusCreated, "ok"
}

func validateDeploy(request *http.Request) ([]*apiserver.Deploy, error) {
	deploys := []*apiserver.Deploy{}
	if err := json.NewDecoder(request.Body).Decode(&deploys); err != nil {
		return nil, err
	}
	return deploys, nil
}

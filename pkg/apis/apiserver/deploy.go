package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"apiserver/pkg/api/apiserver"
	"apiserver/pkg/configz"
	r "apiserver/pkg/router"
	httpUtil "apiserver/pkg/util/registry"
)

func CreatDeploy(request *http.Request) (string, interface{}) {
	deploys, err := validateDeploy(request)
	if err != nil {
		return r.StatusInternalServerError, err
	}

	tranport := httpUtil.GetHttpTransport(false)
	url := configz.GetString("apiserver", "getConfigUrl", "http://localhost:8080/projects/%s/configs")
	client := &http.Client{Transport: tranport}

	for _, deploy := range deploys {
		if err = apiserver.InsertDeploy(deploy); err != nil {
			return r.StatusInternalServerError, err
		}
		for _, item := range deploy.Items {
			projectConfigs := []*apiserver.ProjectConfig{}
			url = fmt.Sprintf(url, item.ProjectId)
			res, err := client.Get(url)
			if err != nil {
				return r.StatusInternalServerError, fmt.Sprintf("get project [%s] config err:%v", item.ProjectName, err.Error())
			}
			if err = json.NewDecoder(res.Body).Decode(&projectConfigs); err != nil {
				return r.StatusInternalServerError, err
			}
			for _, projectConfig := range projectConfigs {
				if err = apiserver.InsertProjectConfig(projectConfig); err != nil {
					return r.StatusInternalServerError, err
				}
			}
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

//反馈接口待定

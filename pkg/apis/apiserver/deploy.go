package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
			projectConfigOptions := []*apiserver.ProjectConfigOption{}
			url = fmt.Sprintf(url, item.ProjectId)
			res, err := client.Get(url)
			if err != nil {
				return r.StatusInternalServerError, fmt.Sprintf("get project [%s] config err:%v", item.ProjectName, err.Error())
			}
			if err = json.NewDecoder(res.Body).Decode(&projectConfigOptions); err != nil {
				return r.StatusInternalServerError, err
			}
			for _, projectConfigOption := range projectConfigOptions {
				createAt, _ := time.Parse("2006-01-02 15:04:05", projectConfigOption.CreateAt)
				updateAt, _ := time.Parse("2006-01-02 15:04:05", projectConfigOption.UpdateAt)
				projectConfig := &apiserver.ProjectConfig{
					ProjectId: projectConfigOption.ProjectId,
					Key:       projectConfigOption.Key,
					Val:       projectConfigOption.Val,
					Type:      projectConfigOption.Type,
					CreateAt:  createAt,
					UpdateAt:  updateAt,
					Operator:  projectConfigOption.Operator,
				}
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

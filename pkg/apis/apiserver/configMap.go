package apiserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"apiserver/pkg/api/apiserver"
	k8sclient "apiserver/pkg/resource/common"
	"apiserver/pkg/resource/configMap"
	r "apiserver/pkg/router"

	"github.com/gorilla/mux"
)

func GetConfigs(request *http.Request) (string, interface{}) {
	pageCnt, _ := strconv.Atoi(request.FormValue("pageCnt"))
	pageNum, _ := strconv.Atoi(request.FormValue("pageNum"))
	configName := request.FormValue("name")
	configs, total := apiserver.QueryConfigs(configName, pageCnt, pageNum)
	return r.StatusOK, map[string]interface{}{"configs": configs, "total": total}
}

func CreateConfig(request *http.Request) (string, interface{}) {
	config, err := validate(request)
	if err != nil {
		return r.StatusInternalServerError, err
	}
	namespace := mux.Vars(request)["namespace"]
	config.Namespace = namespace
	cfgMap := configMap.NewConfigMapByConfig(config)
	if err = k8sclient.CreateResource(&cfgMap); err != nil {
		return r.StatusInternalServerError, err
	}

	apiserver.InsertConfig(config)

	return r.StatusCreated, "ok"
}

func DeleteConfig(request *http.Request) (string, interface{}) {
	id, _ := strconv.ParseUint(mux.Vars(request)["id"], 10, 64)
	cfg := apiserver.QueryConfigById(uint(id))
	cfgMap := configMap.NewConfigMapByConfig(cfg)
	if err := k8sclient.DeleteResource(cfgMap); err != nil {
		return r.StatusInternalServerError, err
	}
	apiserver.DeleteConfig(uint(id))
	return r.StatusOK, "ok"
}

func DeleteConfigItem(request *http.Request) (string, interface{}) {
	itemid, _ := strconv.ParseUint(mux.Vars(request)["itemId"], 10, 64)
	apiserver.DeleteConfigItem(uint(itemid))

	id, _ := strconv.ParseUint(mux.Vars(request)["id"], 10, 64)
	cfg := apiserver.QueryConfigById(uint(id))
	cfgMap := configMap.NewConfigMapByConfig(cfg)
	if err := k8sclient.UpdateResouce(&cfgMap); err != nil {
		return r.StatusInternalServerError, err
	}
	return r.StatusOK, "ok"
}

func CreateConfigItem(request *http.Request) (string, interface{}) {
	id, _ := strconv.ParseUint(mux.Vars(request)["id"], 10, 64)
	cfgMap, err := validateConfigItem(request)
	if err != nil {
		return r.StatusInternalServerError, err
	}
	cfgMap.ConfigGroupId = uint(id)
	apiserver.InsertConfigItem(cfgMap)

	itemid, _ := strconv.ParseUint(mux.Vars(request)["itemId"], 10, 64)
	apiserver.DeleteConfigItem(uint(itemid))

	cfg := apiserver.QueryConfigById(uint(id))
	cfgk8s := configMap.NewConfigMapByConfig(cfg)
	if err := k8sclient.UpdateResouce(&cfgk8s); err != nil {
		return r.StatusInternalServerError, err
	}

	return r.StatusCreated, "ok"
}

func validate(request *http.Request) (*apiserver.ConfigGroup, error) {
	config := &apiserver.ConfigGroup{}
	if err := json.NewDecoder(request.Body).Decode(config); err != nil {
		return nil, err
	}
	return config, nil
}

func validateConfigItem(request *http.Request) (*apiserver.ConfigMap, error) {
	config := &apiserver.ConfigMap{}
	if err := json.NewDecoder(request.Body).Decode(config); err != nil {
		return nil, err
	}
	return config, nil
}

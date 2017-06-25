package autorization

import (
	"encoding/json"
	"net/http"
	"strconv"

	author "apiserver/pkg/api/authorization"
	r "apiserver/pkg/router"

	"github.com/gorilla/mux"
)

func CreatePermission(request *http.Request) (string, interface{}) {
	permission, err := validatePermission(request)
	if err != nil {
		return r.StatusBadRequest, err
	}

	if err = author.InsertPermission(permission); err != nil {
		return r.StatusInternalServerError, err
	}
	return r.StatusCreated, "ok"
}

func GetPermissions(request *http.Request) (string, interface{}) {
	pageCnt, _ := strconv.Atoi(request.FormValue("pageCnt"))
	pageNum, _ := strconv.Atoi(request.FormValue("pageNum"))
	name := request.FormValue("name")
	permissions, total, err := author.QueryPermissions(name, pageCnt, pageNum)
	if err != nil {
		return r.StatusInternalServerError, err
	}
	return r.StatusOK, map[string]interface{}{"permissions": permissions, "total": total}
}

func DeletePermission(request *http.Request) (string, interface{}) {
	id, _ := strconv.ParseUint(mux.Vars(request)["id"], 10, 64)
	if err := author.DeletePermission(uint(id)); err != nil {
		return r.StatusInternalServerError, err
	}
	return r.StatusNoContent, "ok"
}

func validatePermission(request *http.Request) (*author.Permission, error) {
	permission := &author.Permission{}
	if err := json.NewDecoder(request.Body).Decode(permission); err != nil {
		return nil, err
	}
	return permission, nil
}

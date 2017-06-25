package autorization

import (
	"encoding/json"
	"net/http"
	"strconv"

	author "apiserver/pkg/api/authorization"
	r "apiserver/pkg/router"

	"github.com/gorilla/mux"
)

func CreateUser(request *http.Request) (string, interface{}) {
	user, err := validateUser(request)
	if err != nil {
		return r.StatusBadRequest, err
	}

	if err = author.InsertUser(user); err != nil {
		return r.StatusInternalServerError, err
	}
	return r.StatusCreated, "ok"
}

func GetUsers(request *http.Request) (string, interface{}) {
	pageCnt, _ := strconv.Atoi(request.FormValue("pageCnt"))
	pageNum, _ := strconv.Atoi(request.FormValue("pageNum"))
	name := request.FormValue("name")
	users, total, err := author.QueryUsers(name, pageCnt, pageNum)
	if err != nil {
		return r.StatusInternalServerError, err
	}
	return r.StatusOK, map[string]interface{}{"users": users, "total": total}
}

func DeleteUser(request *http.Request) (string, interface{}) {
	id, _ := strconv.ParseUint(mux.Vars(request)["id"], 10, 64)
	if err := author.DeleteUser(uint(id)); err != nil {
		return r.StatusInternalServerError, err
	}
	return r.StatusNoContent, "ok"
}

func validateUser(request *http.Request) (*author.User, error) {
	user := &author.User{}
	if err := json.NewDecoder(request.Body).Decode(user); err != nil {
		return nil, err
	}
	return user, nil
}

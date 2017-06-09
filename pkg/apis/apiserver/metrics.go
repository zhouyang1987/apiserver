package apiserver

import (
	"net/http"

	"github.com/gorilla/mux"

	m "apiserver/pkg/resource/metrics"
	"apiserver/pkg/router"
)

func GetMetrics(request *http.Request) (string, interface{}) {
	namespace := mux.Vars(request)["namespace"]
	podName := mux.Vars(request)["name"]
	metricsName := mux.Vars(request)["metric"] + "/" + mux.Vars(request)["type"]
	metrics, err := m.GetPodMetrics(namespace, podName, metricsName)
	if err != nil {
		return router.StatusInternalServerError, err
	}
	return router.StatusOK, map[string]interface{}{"metrics": metrics}
}

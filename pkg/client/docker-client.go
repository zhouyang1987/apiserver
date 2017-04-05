package client

import (
	"net/http"

	"apiserver/pkg/configz"
	"apiserver/pkg/util/log"

	"github.com/docker/docker/client"
)

var (
	DockerClient *client.Client
)

func init() {
	host := configz.GetString("build", "endpoint", "127.0.0.1：2375")
	version := configz.GetString("build", "version", "12.4")
	cl := &http.Client{}
	DockerClient, err = client.NewClient(host, version, cl, nil)
	if err != nil {
		log.Fatalf("init docker client err: %v", err)
	}
}

package client

import (
	"net/http"
)

var (
	Heapsterclient *http.Client
)

func init() {
	Heapsterclient = http.DefaultClient
}

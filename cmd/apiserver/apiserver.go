package main

import (
	"apiserver/cmd/apiserver/app"
	"apiserver/pkg/util/log"
)

func main() {
	s := app.NewApiServer()
	if err := app.Run(s); err != nil {
		log.Fatalf("start apiserver err: %v", err)
	}
}

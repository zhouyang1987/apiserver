package main

import (
	"apiserver/cmd/docker-build/app"
	"apiserver/pkg/util/log"
)

func main() {
	s := app.NewBuildServer()
	if err := app.Run(s); err != nil {
		log.Fatalf("start apiserver err: %v", err)
	}
}

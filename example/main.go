package main

import (
	"log"

	"github.com/mikkeloscar/gin-swagger/example/restapi"
)

func main() {
	var apiConfig restapi.Config

	err := apiConfig.WithDefaultFlags().Parse()
	if err != nil {
		log.Fatal(err)
	}

	svc := &ExampleService{Health: false}

	api := restapi.NewServer(svc, &apiConfig)

	err = api.RunWithSigHandler()
	if err != nil {
		log.Fatal(err)
	}
}

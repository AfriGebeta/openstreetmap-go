package main

import (
	"fmt"
	"log"
	"net/http"
	"openstreetmap-go/src/config"
	"openstreetmap-go/src/modules"
)

func main() {
	config.SetupEnv()
	fmt.Println("the db username is ", config.ServerConfiguration.DBUserName)
	var httpServer = http.Server{
		Addr:    ":3003",
		Handler: modules.SetupApiRoute(),
	}
	fmt.Println("listening on :3003")
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"net/http"
	apiMiddlewares "openstreetmap-go/src/api/middlewares"
	apiUtils "openstreetmap-go/src/api/utils"
	cacheClient "openstreetmap-go/src/cache/client"
	"openstreetmap-go/src/config"
	dbClient "openstreetmap-go/src/db/client"
	"openstreetmap-go/src/modules"
	"openstreetmap-go/src/utils"
)

func main() {
	config.Load(&config.ServerConfigObject)

	dbClient.Load(&config.ServerConfigObject, &dbClient.DatabaseClients)
	defer dbClient.DatabaseClients.CloseDbConnections()
	cacheClient.Load(&config.ServerConfigObject, &cacheClient.RedisClient)

	defer cacheClient.RedisClient.CloseCacheConnection()

	var middlewareStack = apiUtils.CreateMiddlewareStack(
		apiMiddlewares.LogAccess,
		apiMiddlewares.Cors,
	)

	var httpServer = http.Server{
		Addr:    ":" + config.ServerConfigObject.Port,
		Handler: middlewareStack(modules.SetupApiRoutes()),
	}

	utils.LogInfo("server", "listening on port "+config.ServerConfigObject.Port)

	var err = httpServer.ListenAndServe()

	if err != nil {
		utils.LogFatal("server", "failed to start the server - %s", err.Error())
	}
}

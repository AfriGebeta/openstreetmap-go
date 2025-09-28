package v1

import (
	"net/http"
	"openstreetmap-go/src/api/middlewares"
)

func Setup() *http.ServeMux {
	route := http.NewServeMux()
	route.HandleFunc("PUT /create", middlewares.HandleError(create))
	route.HandleFunc("POST /{id}/upload", middlewares.HandleError(upload))
	route.HandleFunc("PUT /{id}/close", middlewares.HandleError(closeChangeset))
	return route
}

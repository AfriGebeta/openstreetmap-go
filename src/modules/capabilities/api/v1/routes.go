package v1

import (
	"net/http"
	"openstreetmap-go/src/api/middlewares"
)

func Setup() *http.ServeMux {
	var router = http.NewServeMux()

	router.HandleFunc("/capability", middlewares.HandleError(capability))

	return router
}

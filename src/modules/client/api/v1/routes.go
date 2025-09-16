package v1

import (
	"net/http"
	apiMiddlewares "openstreetmap-go/src/api/middlewares"
)

func Setup() *http.ServeMux {
	var router = http.NewServeMux()

	// TODO
	router.HandleFunc("/application", apiMiddlewares.HandleError(createClient))

	return router
}

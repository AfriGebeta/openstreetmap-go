package v1

import (
	"net/http"
	apiMiddlewares "openstreetmap-go/src/api/middlewares"
)

func Setup() *http.ServeMux {
	var router = http.NewServeMux()

	router.HandleFunc("/base", apiMiddlewares.HandleError(confirmHealth))

	router.HandleFunc("/json", apiMiddlewares.HandleError(confirmJsonHealth))

	router.HandleFunc("/no-content", apiMiddlewares.HandleError(confirmNoContentHealth))

	return router
}

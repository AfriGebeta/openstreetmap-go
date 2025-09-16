package v1

import (
	"net/http"
	"openstreetmap-go/src/api/middlewares"
)

func Setup() *http.ServeMux {
	var router = http.NewServeMux()

	router.HandleFunc("/authorize", middlewares.HandleError(authorize))
	router.HandleFunc("/authenticate", middlewares.HandleError(authenticate))
	router.HandleFunc("/token", middlewares.HandleError(getToken))

	return router
}

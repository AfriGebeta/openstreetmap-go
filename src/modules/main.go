package modules

import (
	"net/http"
	capabilities "openstreetmap-go/src/modules/capabilities/api/v1"
)

func SetupApiRoute() *http.ServeMux {
	var router = http.NewServeMux()
	router.Handle("/api/0.6/capabilities", http.StripPrefix("/api/0.6", SetupV1()))
	router.Handle("/api/capabilities", http.StripPrefix("/api", SetupV1()))
	return router
}

func SetupV1() *http.ServeMux {
	var router = http.NewServeMux()
	router.Handle("/capabilities", capabilities.Setup())
	return router
}

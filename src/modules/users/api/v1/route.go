package v1

import "net/http"

func SetupUserV1Route() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {

	})
	return router

}

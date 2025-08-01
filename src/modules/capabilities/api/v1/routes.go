package v1

import (
	"net/http"
)

func Setup() *http.ServeMux {
	var router = http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := getCapabilites(w, r)
		if err != nil {
			// send the error here if no error the response is already returned by the function
		}
	})

	return router
}

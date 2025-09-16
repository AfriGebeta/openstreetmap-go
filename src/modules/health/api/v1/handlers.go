package v1

import (
	"net/http"
	apiValues "openstreetmap-go/src/api/values"
)

func confirmHealth(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(200)
	var _, err = w.Write([]byte("ok"))

	if err != nil {
		return apiValues.ApiErrorInternalServerError
	}

	return nil
}

func confirmJsonHealth(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(200)
	var _, err = w.Write([]byte("{\"status\":\"ok\"}"))

	if err != nil {
		return apiValues.ApiErrorInternalServerError
	}

	return nil
}

func confirmNoContentHealth(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(204)
	var _, err = w.Write(nil)

	if err != nil {
		return apiValues.ApiErrorInternalServerError
	}

	return nil
}

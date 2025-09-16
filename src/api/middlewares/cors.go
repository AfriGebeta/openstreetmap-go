package middlewares

import (
	"net/http"
	"openstreetmap-go/src/api/model"
	"openstreetmap-go/src/api/utils"
	"openstreetmap-go/src/config"
	"slices"
)

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var origin = r.Header.Get("Origin")

		if origin == "" {
			utils.UpdateContextAndServe(w, r,
				model.RequestContextTagAlong{},
				next.ServeHTTP,
			)

			return
		}

		if !slices.Contains(config.ServerConfigObject.AllowedOrigins, origin) &&
			!slices.Contains(config.ServerConfigObject.AllowedOrigins, "*") {
			w.WriteHeader(403)
			return
		}

		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Expose-Headers", "*")

		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
			w.Header().Set("Access-Control-Allow-Headers", "*") // This is intentional. I do not think it is necessary to restrict this, and it is just a hassle to work with

			w.WriteHeader(http.StatusNoContent)
			return
		}

		utils.UpdateContextAndServe(w, r,
			model.RequestContextTagAlong{},
			next.ServeHTTP,
		)
	})
}

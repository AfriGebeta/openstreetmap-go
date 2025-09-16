package middlewares

import (
	"net/http"
	apiModel "openstreetmap-go/src/api/model"
	apiUtils "openstreetmap-go/src/api/utils"
	"time"
)

func LogAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var startTime = time.Now()

		var wrappedWriter = &apiModel.WrappedWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusNotImplemented,
		}

		var wrappedContext = &apiModel.WrappedContext{
			Context:  r.Context(),
			TagAlong: apiModel.RequestContextTagAlong{},
			Parent:   nil,
		}

		var req = r.WithContext(wrappedContext)

		next.ServeHTTP(wrappedWriter, req)

		apiUtils.LogAPIAccess(wrappedWriter, req, startTime, "api", []string{})
	})
}

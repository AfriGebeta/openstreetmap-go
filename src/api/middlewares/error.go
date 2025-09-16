package middlewares

import (
	"encoding/json"
	"net/http"
	apiModel "openstreetmap-go/src/api/model"
	apiUtils "openstreetmap-go/src/api/utils"
	utils2 "openstreetmap-go/src/common/utils"
	"openstreetmap-go/src/utils"
)

func HandleError(handler apiModel.RequestHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err = handler(w, r)

		if err != nil {
			var errRole = apiUtils.GetErrorRole(r)

			var orderedLocalePreference = apiUtils.ParseAcceptLanguageHeader(r.Header.Get("Accept-Language"))

			var errRes, errorList, status = apiUtils.GetError(err, errRole, orderedLocalePreference...)

			utils.LogError("api", utils2.GetPrintableErrorFromErrorList("api", errorList, errRole, orderedLocalePreference...))

			var payload, _ = json.Marshal(apiModel.GetErrorResponsePayload(errRes...))

			w.Header().Set("Content-Type", "application/json")

			w.WriteHeader(status)

			_, err = w.Write(payload)

			if err != nil {
				utils.LogError("error-handler", "could not send payload")
			}
		}
	}
}

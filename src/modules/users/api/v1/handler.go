package v1

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"openstreetmap-go/src/utils"
	"strings"
)

// parse and validate post request into the expected format
func ParseRequest(r *http.Request, expectedFormat interface{}) (interface{}, *utils.ApiError) {

	if r.Method != http.MethodPost {
		return nil, utils.NewApiError(http.StatusMethodNotAllowed, "invalid method: only post allowed", nil, nil)
	}
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(expectedFormat)
	if err != nil {
		return nil, utils.NewApiError(422, "parsing request body failed", err, nil)
	}

	validate := validator.New()
	err = validate.Struct(expectedFormat)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			fieldErrors := make(map[string]string)
			for _, fieldError := range validationErrors {
				field := fieldError.Field()
				fieldErrors[field] = fieldError.Tag()

			}
			return nil, utils.NewApiError(422, "validation failed", err, fieldErrors)
		}
		return nil, utils.NewApiError(422, "Validation error", err, nil)
	}

	return expectedFormat, nil
}

func CreateUser(w http.ResponseWriter, r *http.Request) error {
	var req UserCreate
	body, err := ParseRequest(r, &req)
	if err != nil {
		return err
	}
	createReq := body.(*UserCreate)
	// check acl for domain restriction
	domain := strings.Split(createReq.Email, "@")[1]

	return nil
}

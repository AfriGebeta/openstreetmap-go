package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"openstreetmap-go/src/api/utils"
	apiValues "openstreetmap-go/src/api/values"
	cacheClient "openstreetmap-go/src/cache/client"
	"openstreetmap-go/src/domain/model"
)

func createClient(w http.ResponseWriter, r *http.Request) error {
	var ctx, err = utils.GetWrappedContext(r)
	if err != nil {
		return err
	}
	body, err := utils.DecodeBody[CreateClient](r)
	if err != nil {
		return err
	}

	clientId, err := GenerateRandomHex(16)
	if err != nil {
		return errors.Join(apiValues.ApiErrorFailedValidation, err)
	}
	secret, err := GenerateRandomHex(16)
	if err != nil {
		return errors.Join(apiValues.ApiErrorFailedValidation, err)
	}

	code := model.Client{
		ClientId:     clientId,
		ClientSecret: secret,
		RedirectURI:  body.RedirectURI,
	}
	clientModelUnmarshaled, err := json.Marshal(code)
	if err != nil {
		return errors.Join(apiValues.ApiErrorFailedValidation, err)
	}

	err = cacheClient.RedisClient.Set(ctx, clientId, string(clientModelUnmarshaled), 0).Err()
	if err != nil {
		return errors.Join(apiValues.ApiErrorInternalServerError, err)
	}

	clientModel := CreateReturnType{
		RedirectURI: body.RedirectURI,
		ClientID:    clientId,
	}

	response, err := json.Marshal(clientModel)
	if err != nil {
		return errors.Join(apiValues.ApiErrorFailedValidation, err)
	}

	if err = utils.EncodeBody(w, response, 200); err != nil {
		return apiValues.ApiErrorInternalServerError
	}

	return nil
}

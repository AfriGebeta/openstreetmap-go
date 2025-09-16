package v1

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"net/http"
	"net/url"
	"openstreetmap-go/src/api/utils"
	apiValues "openstreetmap-go/src/api/values"
	cacheClient "openstreetmap-go/src/cache/client"
	"openstreetmap-go/src/config"
	domainModel "openstreetmap-go/src/domain/model"
	domainUtils "openstreetmap-go/src/domain/utils"
	"openstreetmap-go/src/domain/values"
	"time"
)

func generateRandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func authorize(w http.ResponseWriter, r *http.Request) error {
	var ctx, err = utils.GetWrappedContext(r)
	if err != nil {
		return errors.Join(apiValues.ApiErrorInternalServerError, err)
	}
	queryParams, err := parseAuthorizeQuery(r)
	if err != nil {
		return errors.Join(apiValues.ApiErrorFailedValidation, err)
	}

	clientInformation, err := cacheClient.RedisClient.Get(ctx, domainUtils.GetCachePrefix(values.CachePrefixValue.Account)+queryParams.ClientID).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return errors.Join(apiValues.ApiErrorNotFound, errors.New("Invalid client id"))
		} else {
			return errors.Join(apiValues.ApiErrorInternalServerError, err)
		}
	}

	var client domainModel.Client
	err = json.Unmarshal([]byte(clientInformation), &client)
	if err != nil {
		return errors.Join(apiValues.ApiErrorInternalServerError, err)
	}
	if queryParams.RedirectURI != client.RedirectURI {
		return errors.Join(apiValues.ApiErrorFailedValidation, errors.New("Invalid redirect uri"))
	}

	if !isHashingAllowed(queryParams.CodeChallengeMethod) {
		return errors.Join(apiValues.ApiErrorFailedValidation, errors.New("Invalid code challenge_method"))
	}
	flowId, err := generateRandomHex(32)
	codeModel := CodeModel{
		ClientId:            queryParams.ClientID,
		State:               queryParams.State,
		CodeChallenge:       queryParams.CodeChallenge,
		CodeChallengeMethod: queryParams.CodeChallengeMethod,
		RedirectURI:         queryParams.RedirectURI,
	}

	marshalled, err := json.Marshal(codeModel)
	if err != nil {
		return errors.Join(apiValues.ApiErrorInternalServerError, err)
	}

	err = cacheClient.RedisClient.Set(ctx, domainUtils.GetCachePrefix(values.CachePrefixValue.AuthorizationRequest)+flowId, string(marshalled), time.Minute*5).Err()
	if err != nil {
		return errors.Join(apiValues.ApiErrorInternalServerError, err)
	}

	http.Redirect(w, r, "/v1/oauth2/login?flow_id="+flowId, http.StatusTemporaryRedirect)

	return nil
}
func authenticate(w http.ResponseWriter, r *http.Request) error {

	var ctx, err = utils.GetWrappedContext(r)
	if err != nil {
		return errors.Join(apiValues.ApiErrorInternalServerError, err)
	}

	body, err := utils.DecodeBody[AuthenticateModel](r)
	if err != nil {
		return errors.Join(apiValues.ApiErrorFailedValidation, err)
	}

	existingOtp, err := cacheClient.RedisClient.Get(ctx, domainUtils.GetCachePrefix(values.CachePrefixValue.Otp)+body.FlowId).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return errors.Join(apiValues.ApiErrorFailedValidation, errors.New("Invalid otp"))
		} else {
			return errors.Join(apiValues.ApiErrorInternalServerError, err)
		}
	}
	for _, challenge := range body.Challenges {
		if challenge.Challenge != existingOtp {
			return errors.Join(apiValues.ApiErrorFailedValidation, errors.New("Invalid otp"))
		}
	}

	flowModel, err := cacheClient.RedisClient.Get(ctx, domainUtils.GetCachePrefix(values.CachePrefixValue.AuthorizationRequest)+body.FlowId).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return errors.Join(apiValues.ApiErrorFailedValidation, errors.New("Invalid code"))
		} else {
			return errors.Join(apiValues.ApiErrorInternalServerError, err)
		}
	}

	var _flowModel CodeModel
	err = json.Unmarshal([]byte(flowModel), &_flowModel)
	if err != nil {
		return errors.Join(apiValues.ApiErrorFailedValidation, err)
	}

	authCode, err := generateRandomHex(16)
	if err != nil {
		return errors.Join(apiValues.ApiErrorInternalServerError, err)
	}

	authCodeModel := CodeModel{
		ClientId:            _flowModel.ClientId,
		State:               _flowModel.State,
		CodeChallenge:       _flowModel.CodeChallenge,
		CodeChallengeMethod: _flowModel.CodeChallengeMethod,
		RedirectURI:         _flowModel.RedirectURI,
	}

	marshalled, err := json.Marshal(authCodeModel)
	if err != nil {
		return errors.Join(apiValues.ApiErrorInternalServerError, err)
	}

	cacheClient.RedisClient.Set(ctx, domainUtils.GetCachePrefix(values.CachePrefixValue.Authentication)+authCode, string(marshalled), time.Minute*5)
	redirectURL, err := url.Parse(_flowModel.RedirectURI)
	if err != nil {
		return errors.Join(apiValues.ApiErrorFailedValidation, err)
	}

	q := redirectURL.Query()
	q.Set("code", authCode)
	q.Set("state", _flowModel.State)

	redirectURL.RawQuery = q.Encode()

	http.Redirect(w, r, fmt.Sprintf("/v1/oauth2/consent?redirect_uri=", redirectURL), http.StatusTemporaryRedirect)

	return nil
}

func getToken(w http.ResponseWriter, r *http.Request) error {
	var ctx, err = utils.GetWrappedContext(r)
	if err != nil {
		return errors.Join(apiValues.ApiErrorInternalServerError, err)
	}

	body, err := utils.DecodeBody[GetTokenModel](r)
	if err != nil {
		return errors.Join(apiValues.ApiErrorFailedValidation, err)
	}

	if body.GrantType != "authorization_code" {
		return errors.Join(apiValues.ApiErrorFailedValidation, errors.New("Invalid grant type"))
	}

	codeModel, err := cacheClient.RedisClient.Get(ctx, domainUtils.GetCachePrefix(values.CachePrefixValue.Authentication)+body.Code).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return errors.Join(apiValues.ApiErrorFailedValidation, errors.New("Invalid code"))
		} else {
			return errors.Join(apiValues.ApiErrorInternalServerError, err)
		}
	}

	var _codeModel CodeModel
	err = json.Unmarshal([]byte(codeModel), &_codeModel)
	if err != nil {
		return errors.Join(apiValues.ApiErrorFailedValidation, err)
	}

	if body.ClientID == "" || body.ClientID != _codeModel.ClientId {
		return errors.Join(apiValues.ApiErrorFailedValidation, errors.New("Invalid client id"))
	}

	err = ValidatePKCE(body.CodeVerifier, _codeModel.CodeChallenge, _codeModel.CodeChallengeMethod)
	if err != nil {
		return errors.Join(apiValues.ApiErrorFailedValidation, errors.New("Invalid code verifier"))
	}
	err = cacheClient.RedisClient.Del(ctx, domainUtils.GetCachePrefix(values.CachePrefixValue.Authentication)+body.Code).Err()
	if err != nil {
		return errors.Join(apiValues.ApiErrorFailedValidation, err)
	}

	if body.RedirectURI != _codeModel.RedirectURI {
		return errors.Join(apiValues.ApiErrorFailedValidation, errors.New("Invalid redirect uri"))
	}

	claims := AccessTokenClaims{
		Scope: "read",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   _codeModel.State,
			Audience:  jwt.ClaimStrings{body.ClientID},
			Issuer:    config.ServerConfigObject.SelfBaseUrl,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(config.ServerConfigObject.JwtSecret)
	if err != nil {
		return errors.Join(apiValues.ApiErrorFailedValidation, err)
	}

	type ValidateTokenResponse struct {
		AccessToken string `json:"accessToken"`
		TokenType   string `json:"tokenType"`
		ExpiresIn   int    `json:"expiresIn"`
	}

	response := ValidateTokenResponse{
		AccessToken: signedToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(time.Hour.Seconds()),
	}

	if err = utils.EncodeBody(w, response, 200); err != nil {
		return apiValues.ApiErrorInternalServerError
	}

	return nil
}

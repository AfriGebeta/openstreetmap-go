package v1

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
)

func parseAuthorizeQuery(r *http.Request) (*AuthorizeModel, error) {
	queryParams := r.URL.Query()
	requiredFields := []string{"client_id", "redirect_uri", "state", "scope", "code_challenge", "code_challenge_method"}
	for _, field := range requiredFields {
		if queryParams.Get(field) == "" {
			return nil, fmt.Errorf("missing required query parameter: %s", field)
		}
	}

	model := &AuthorizeModel{
		ClientID:            queryParams.Get("client_id"),
		RedirectURI:         queryParams.Get("redirect_uri"),
		State:               queryParams.Get("state"),
		Scope:               queryParams.Get("scope"),
		CodeChallenge:       queryParams.Get("code_challenge"),
		CodeChallengeMethod: queryParams.Get("code_challenge_method"),
	}

	return model, nil
}

func isHashingAllowed(method string) bool {
	//Todo make this configurable
	if method != "S256" {
		return false
	}
	return true
}

func ValidatePKCE(codeVerifier string, codeChallenge string, codeChallengeGenerationMethod string) error {
	hash := sha256.Sum256([]byte(codeVerifier))
	encoded := base64.RawURLEncoding.EncodeToString(hash[:])
	if encoded != codeChallenge {
		return errors.New("invalid_grant: Code verifier does not match challenge")
	}
	return nil
}

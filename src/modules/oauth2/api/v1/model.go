package v1

import "github.com/golang-jwt/jwt/v5"

type GetTokenModel struct {
	GrantType    string `json:"grantType"`
	Code         string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
	ClientID     string `json:"client_id"`
	CodeVerifier string `json:"code_verifier"`
}

type AuthorizeModel struct {
	ClientID            string `json:"clientId"`
	RedirectURI         string `json:"redirectUri"`
	State               string `json:"state"`
	Scope               string `json:"scope"`
	CodeChallenge       string `json:"codeChallenge"`
	CodeChallengeMethod string `json:"codeChallengeMethod"`
}

type AuthenticateModel struct {
	Challenges []OauthChallenge `json:"challenges"`
	FlowId     string           `json:"flowId"`
}

type OauthChallenge struct {
	AuthFactorType string `json:"authFactorType"` //OTP
	Challenge      string `json:"challenge"`      // otp
	Format         string `json:"format"`         // format oidc read
}

type CodeModel struct {
	ClientId            string `json:"clientId"`
	State               string `json:"state"`
	CodeChallenge       string `json:"code_challenge"`
	CodeChallengeMethod string `json:"code_challenge_method"`
	RedirectURI         string `json:"redirect_uri"`
}

type AccessTokenClaims struct {
	Scope string `json:"scope"`
	jwt.RegisteredClaims
}

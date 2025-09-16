package model

type Client struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	RedirectURI  string `json:"redirectUri"`
}

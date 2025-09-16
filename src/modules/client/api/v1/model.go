package v1

type CreateClient struct {
	Name        string `json:"name"`
	RedirectURI string `json:"redirectUri"`
}

type CreateReturnType struct {
	RedirectURI string `json:"redirectUri"`
	ClientID    string `json:"clientId"`
}

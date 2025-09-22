package v1

import (
	_ "github.com/joho/godotenv"
	_ "github.com/shopspring/decimal"
	"net/http"
	"openstreetmap-go/src/api/middlewares"
)

func Setup() *http.ServeMux {
	route := http.NewServeMux()
	route.HandleFunc("POST /create", middlewares.HandleError(CreateUser))
	route.HandleFunc("POST /login", middlewares.HandleError(Login))
	route.HandleFunc("GET /{display_name}/confirm", middlewares.HandleError(Confirm))
	route.HandleFunc("GET /grant_role", middlewares.HandleError(GrantRole))
	return route
}

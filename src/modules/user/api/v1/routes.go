package v1

import (
	_ "github.com/joho/godotenv"
	_ "github.com/shopspring/decimal"
	"net/http"
	"openstreetmap-go/src/api/middlewares"
)

func Setup() *http.ServeMux {
	route := http.NewServeMux()
	route.HandleFunc("/create", middlewares.HandleError(CreateUser))
	route.HandleFunc("/login", middlewares.HandleError(Login))
	route.HandleFunc("/{display_name}/confirm", middlewares.HandleError(Confirm))
	route.HandleFunc("/grant_role", middlewares.HandleError(GrantRole))
	return route
}

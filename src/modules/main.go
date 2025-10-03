package modules

import (
	"github.com/a-h/templ"
	"net/http"
	changesetv1 "openstreetmap-go/src/modules/changeset/api/v1"
	clientApi "openstreetmap-go/src/modules/client/api/v1"
	healthApi "openstreetmap-go/src/modules/health/api/v1"
	mapApiV1 "openstreetmap-go/src/modules/map/api/v1"
	oauth2Api "openstreetmap-go/src/modules/oauth2/api/v1"
	user "openstreetmap-go/src/modules/user/api/v1"
	"openstreetmap-go/template"
)

func SetupApiRoutes() *http.ServeMux {
	var router = http.NewServeMux()
	router.Handle("/", Setup_uiRoutes())
	router.Handle("/api/", http.StripPrefix("/api", setupVersion1()))

	return router
}

func setupVersion1() *http.ServeMux {
	var router = http.NewServeMux()
	router.Handle("/0.6/", http.StripPrefix("/0.6", Setup0_6Routes()))
	return router
}

func Setup_uiRoutes() *http.ServeMux {
	var router = http.NewServeMux()

	component := template.LoginPage("login page")
	router.Handle("/login", templ.Handler(component))
	return router
}

func Setup0_6Routes() *http.ServeMux {
	var router = http.NewServeMux()
	router.Handle("/users/", http.StripPrefix("/users", user.Setup()))
	router.Handle("/health/", http.StripPrefix("/health", healthApi.Setup()))
	router.Handle("/oauth2/", http.StripPrefix("/oauth2", oauth2Api.Setup()))
	router.Handle("/client/", http.StripPrefix("/client", clientApi.Setup()))
	router.Handle("/changeset/", http.StripPrefix("/changeset", changesetv1.Setup()))
	router.Handle("/map", mapApiV1.Setup())

	return router
}

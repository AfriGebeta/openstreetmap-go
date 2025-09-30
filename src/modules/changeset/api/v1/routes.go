package v1

import (
	"net/http"
	"openstreetmap-go/src/api/middlewares"
)

func Setup() *http.ServeMux {
	route := http.NewServeMux()
	//api/changesets#index
	route.HandleFunc("GET /", middlewares.HandleError(getChangesets))

	// api/changesetes#create
	route.HandleFunc("POST /", middlewares.HandleError(createChangeset))

	// api/changesetes#create
	route.HandleFunc("PUT /create", middlewares.HandleError(create))

	//api/changesets/uploads#create
	route.HandleFunc("POST /{id}/upload", middlewares.HandleError(upload))

	// api/changesets/closes#update`
	route.HandleFunc("PUT /{id}/close", middlewares.HandleError(closeChangeset))

	// api/changesets#show
	route.HandleFunc("GET /{id}", middlewares.HandleError(getChangesetById))

	// `api/changesets#update`
	route.HandleFunc("PUT /{id}", middlewares.HandleError(updateChangeSet))

	// api/changesets/downloads#show
	route.HandleFunc("GET /{id}/download", middlewares.HandleError(downloadChangeset))

	// api/changeset_subscriptions#create`
	route.HandleFunc("POST /{id}/subscribe", middlewares.HandleError(subscribeToChangeset))

	// `api/changeset_subscriptions#destroy`
	route.HandleFunc("DELETE /{id}/unsubscribe", middlewares.HandleError(unsubscribeToChangeset))

	// api/changeset_comments#create`
	route.Handle("POST /{id}/comment", middlewares.HandleError(createComment))

	return route
}

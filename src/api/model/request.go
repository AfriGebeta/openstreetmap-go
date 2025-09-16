package model

import (
	"context"
	"net/http"
)

type Middleware func(next http.Handler) http.Handler

type WrappedContext struct {
	context.Context
	Parent   *WrappedContext
	TagAlong RequestContextTagAlong
}

type WrappedWriter struct {
	http.ResponseWriter
	StatusCode int
}

type WrappedRequest struct {
	*http.Request
}

func (w *WrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.StatusCode = statusCode
}

func (w *WrappedWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *WrappedWriter) Write(b []byte) (int, error) {
	return w.ResponseWriter.Write(b)
}

func (r *WrappedRequest) Context() WrappedContext {
	var tagAlong = RequestContextTagAlong{}

	var ctx, ok = r.Request.Context().(WrappedContext)

	if ok {
		tagAlong = ctx.TagAlong
	} else {
		ctx = WrappedContext{}
	}

	return WrappedContext{
		Context:  r.Request.Context(),
		TagAlong: tagAlong,
		Parent:   ctx.Parent,
	}
}

type RequestContextTagAlong struct {
	TracingId string // This has to be a combination of the actor type and the actor id followed by a random string representing the scenario. Something like (SYSTEM|USER)_actor-id_scenario-id_start-time-in-unix-milli
	Data      map[string]interface{}
	// TODO add any other required
}

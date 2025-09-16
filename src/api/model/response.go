package model

import "net/http"

type RequestHandler func(w http.ResponseWriter, r *http.Request) error

type ResponsePayload = map[string]any

func GetSuccessResponsePayload[T any](v T) ResponsePayload {
	return ResponsePayload{
		"data": v,
	}
}

func GetErrorResponsePayload[T any](v ...T) ResponsePayload {
	var tmp = ResponsePayload{
		"error":  v[0],
		"errors": v,
	}

	return tmp
}

type SuccessResponsePayload[T any] struct {
	Data T `json:"data"`
}

type ErrorResponsePayload[T any] struct {
	Error T `json:"error"`
}

type ErrorResponse struct {
	Namespace    string          `json:"-"`
	Code         string          `json:"code"`
	Message      string          `json:"message"`
	Elaborations []ErrorResponse `json:"elaborations,omitempty"`
}

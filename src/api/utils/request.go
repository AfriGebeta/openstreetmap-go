package utils

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	apiModel "openstreetmap-go/src/api/model"
	apiValues "openstreetmap-go/src/api/values"
	commonModel "openstreetmap-go/src/common/model"
	"openstreetmap-go/src/utils"
	"strings"
)

func CreateMiddlewareStack(ms ...apiModel.Middleware) apiModel.Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(ms) - 1; i >= 0; i-- {
			var x = ms[i]
			next = x(next)
		}

		return next
	}
}

func EncodeBody[T any](w http.ResponseWriter, v T, status int, format ...string) error {
	var to = "json"

	if len(format) > 0 {
		to = format[0]
	}

	switch to {
	case "json":
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(status)
		type Response[T any] struct {
			Data T `json:"data"`
		}

		if err := json.NewEncoder(w).Encode(Response[T]{Data: v}); err != nil {
			utils.LogError("encoder", "could not encode json %s", err.Error())
			return err
		}

		break
	default:
		w.WriteHeader(status)

		var err error

		if byteArr := utils.SafeCastValue[[]byte](v); byteArr == nil || *byteArr == nil {
			_, err = w.Write(nil)
		} else {
			_, err = w.Write(*byteArr)
		}

		if err != nil {
			utils.LogError("encoder", "could not encode json %s", err.Error())
			return err
		}
	}

	return nil
}

func GetBodyAsText(r *http.Request) (string, error) {
	var body = r.Body

	defer func(body io.ReadCloser) {
		err := body.Close()

		if err != nil {
			utils.LogError("api", "could not close body: %s", err.Error())
		}
	}(body)

	var bodyAsText, err = io.ReadAll(body)

	if err != nil {
		utils.LogError("api", "could not read body as text: %s", err.Error())
		return "", err
	}

	return string(bodyAsText), nil
}

func DecodeBody[T any](r *http.Request, format ...string) (T, error) {
	var v T

	var to = "json"

	if len(format) > 0 {
		to = format[0]
	}

	switch to {
	case "json":
		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			utils.LogError("decoder", "could not decode json %s", err.Error())

			return v, errors.Join(apiValues.ApiErrorUnprocessableEntity, errors.New("body is not a valid json"))
		}
		break
	}

	return v, nil
}

func GetWrappedContext(r *http.Request) (*apiModel.WrappedContext, error) {
	if value, ok := r.Context().(*apiModel.WrappedContext); ok {
		return value, nil
	} else {
		return &apiModel.WrappedContext{}, errors.Join(apiValues.ApiErrorInternalServerError, errors.New("could not get request context"))
	}
}

func UpdateContextAndServe(w http.ResponseWriter, r *http.Request, tagAlong apiModel.RequestContextTagAlong, handler http.HandlerFunc) {
	var parent = &apiModel.WrappedContext{}

	var currentContext, err = GetWrappedContext(r)

	if err == nil {
		parent = currentContext
	}

	var req = r.WithContext(&apiModel.WrappedContext{Context: currentContext.Context, TagAlong: tagAlong, Parent: parent})

	handler(w, req)

	if parent != nil {
		curr, err := GetWrappedContext(req)

		if err == nil && curr != nil {
			parent.TagAlong = curr.TagAlong
		}
	}
}

func StoreValueInContext(key string, value interface{}, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(
			w,
			r.WithContext(
				context.WithValue(r.Context(), key, value),
			),
		)
	})
}

func GetRequestIP(r *http.Request) string {
	var xForwardedFor = r.Header.Get("X-Forwarded-For")

	if xForwardedFor != "" {
		var ips = strings.Split(xForwardedFor, ",")
		return strings.TrimSpace(ips[0])
	}

	var xRealIP = r.Header.Get("X-Real-IP")

	if xRealIP != "" {
		return xRealIP
	}

	return r.RemoteAddr
}

func GetSessionData(r *http.Request) (userAgent, ip string, clientPlatform commonModel.ClientPlatform, err error) {
	userAgent = r.Header.Get("User-Agent")

	ip, _, err = net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		err = errors.Join(apiValues.ApiErrorFailedValidation, errors.New("suspicious activity detected"))
	}

	clientPlatform = IdentifyClientPlatform(userAgent)

	return
}

func IdentifyClientPlatform(userAgent string) commonModel.ClientPlatform {
	userAgent = strings.ToLower(userAgent)

	switch {
	case strings.Contains(userAgent, "mozilla") || strings.Contains(userAgent, "chrome") || strings.Contains(userAgent, "safari"):
		return commonModel.ClientPlatformValue.Web
	case strings.Contains(userAgent, "windows"):
		return commonModel.ClientPlatformValue.Windows
	case strings.Contains(userAgent, "macintosh") || strings.Contains(userAgent, "mac os x"):
		return commonModel.ClientPlatformValue.Darwin
	case strings.Contains(userAgent, "linux"):
		return commonModel.ClientPlatformValue.Linux
	case strings.Contains(userAgent, "android"), strings.Contains(userAgent, "postdroid"):
		return commonModel.ClientPlatformValue.Android
	case strings.Contains(userAgent, "iphone") || strings.Contains(userAgent, "ipad") || strings.Contains(userAgent, "parcelpod"):
		return commonModel.ClientPlatformValue.Ios
	default:
		return commonModel.ClientPlatformValue.Unknown
	}
}

func ParseAcceptLanguageHeader(header string) []commonModel.Locale {
	languages := strings.Split(header, ",")

	var cleanLanguages []commonModel.Locale

	for _, lang := range languages {
		// Remove any quality values if present things like q=0.8
		if index := strings.Index(lang, ";"); index != -1 {
			lang = strings.TrimSpace(lang[:index])
		} else {
			lang = strings.TrimSpace(lang)
		}

		cleanLanguages = append(cleanLanguages, commonModel.Locale(lang))
	}

	return cleanLanguages
}

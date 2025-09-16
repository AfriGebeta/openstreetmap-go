package model

import (
	"errors"
	"fmt"
)

type ErrorList map[ErrorCode]*Error

type ErrorCode string

func (n ErrorCode) String() string {
	return string(n)
}

func (n ErrorCode) StringPtr() *string {
	var tmp = n.String()
	return &tmp
}

func (n ErrorCode) Join(err ...error) error {
	var _err error

	if len(err) == 0 || err[0] == nil {
		return errors.New(n.String())
	}

	_err = err[0]

	return fmt.Errorf("%s%w", n.String(), _err)
}

type Error struct {
	Layer          ErrorLayer         `json:"layer"`
	Namespace      ErrorNamespace     `json:"namespace"`
	Code           ErrorCode          `json:"code"`
	Message        ErrorMessage       `json:"message"`
	RoleVisibility map[ErrorRole]bool `json:"roleVisibility"`
	Additional     interface{}        `json:"additional"`
}

func (e *Error) Clone() *Error {
	if e == nil {
		return nil
	}

	return &Error{
		Layer:          e.Layer,
		Namespace:      e.Namespace,
		Code:           e.Code,
		Message:        e.Message,
		RoleVisibility: e.RoleVisibility,
		Additional:     e.Additional,
	}
}

func (e *Error) SetNamespace(namespace ErrorNamespace) *Error {
	if e == nil {
		return nil
	}

	e.Namespace = namespace

	return e
}

func (e *Error) SetCode(code ErrorCode) *Error {
	if e == nil {
		return nil
	}

	e.Code = code

	return e
}

func (e *Error) SetMessage(message ErrorMessage) *Error {
	if e == nil {
		return nil
	}

	e.Message = message

	return e
}

func (e *Error) SetAdditional(additional interface{}) *Error {
	if e == nil {
		return nil
	}

	e.Additional = additional

	return e
}

func (e *Error) SetRoleVisibility(role ErrorRole, visible bool) *Error {
	if e == nil {
		return nil
	}

	e.RoleVisibility[role] = visible

	return e
}

func (e *Error) SetLayer(layer ErrorLayer) *Error {
	if e == nil {
		return nil
	}

	e.Layer = layer

	return e
}

type ErrorLayer string

var ErrorLayerValue = struct {
	Common     ErrorLayer
	System     ErrorLayer
	Data       ErrorLayer
	Repository ErrorLayer
	Domain     ErrorLayer
	Api        ErrorLayer
}{
	Common:     "COMMON",
	System:     "SYSTEM",
	Data:       "DATA",
	Repository: "REPOSITORY",
	Domain:     "DOMAIN",
	Api:        "API",
}

var ErrorLayerValues = []ErrorLayer{
	ErrorLayerValue.Common,
	ErrorLayerValue.System,
	ErrorLayerValue.Data,
	ErrorLayerValue.Repository,
	ErrorLayerValue.Domain,
	ErrorLayerValue.Api,
}

func (n ErrorLayer) String() string {
	return string(n)
}

func (n ErrorLayer) StringPtr() *string {
	var tmp = n.String()
	return &tmp
}

type ErrorNamespace string

func (n ErrorNamespace) String() string {
	return string(n)
}

func (n ErrorNamespace) StringPtr() *string {
	var tmp = n.String()
	return &tmp
}

type ErrorMessage map[ErrorRole]IntlText

func (em ErrorMessage) RoleMessage(role ...ErrorRole) IntlText {
	if em == nil {
		return nil
	}

	if len(role) > 0 {
		return em[role[0]]
	}

	return em["COMMON"]
}

type ErrorRole string

func (n ErrorRole) String() string {
	return string(n)
}

func (n ErrorRole) StringPtr() *string {
	var tmp = n.String()
	return &tmp
}

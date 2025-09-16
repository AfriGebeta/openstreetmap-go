package utils

import (
	"errors"
	"net/http"
	apiModel "openstreetmap-go/src/api/model"
	apiValues "openstreetmap-go/src/api/values"
	commonModel "openstreetmap-go/src/common/model"
	commonValues "openstreetmap-go/src/common/values"
	domainValues "openstreetmap-go/src/domain/values"
	"openstreetmap-go/src/utils"
	"strings"
)

func GetError(err error, role commonModel.ErrorRole, orderedLocalePreference ...commonModel.Locale) (errRes []apiModel.ErrorResponse, errorList []error, status int) {
	status = ObtainErrorStatus(err)

	errorList = utils.RecursivelyUnwrapJoin(err)

	if len(errorList) == 0 || errorList[0] == nil || (status >= 500 && status != http.StatusServiceUnavailable && status != http.StatusNotImplemented) {
		var message = commonValues.ErrorList[apiValues.ErrorCodeApiInternalServerError].Message

		var roleMessage = message.RoleMessage(role)

		if len(roleMessage) == 0 {
			roleMessage = message.RoleMessage("DEFAULT")
		}

		errRes = []apiModel.ErrorResponse{
			{
				Code:      commonValues.ErrorList[apiValues.ErrorCodeApiInternalServerError].Code.String(),
				Message:   roleMessage.LocaleText(orderedLocalePreference...),
				Namespace: commonValues.ErrorList[apiValues.ErrorCodeApiInternalServerError].Namespace.String(),
			},
		}

		return
	}

	var prevErr *commonModel.Error

	for _, e := range errorList {
		var _err, ok = commonValues.ErrorList[commonModel.ErrorCode(e.Error())]

		var tmp = apiModel.ErrorResponse{}

		if strings.HasPrefix(e.Error(), "RE") {
			break
		}

		var shouldBreak = errors.Is(e, domainValues.DomainErrorUnknown)

		if !ok {
			tmp.Code = "UNKNOWN"
			tmp.Message = e.Error()
			tmp.Namespace = "UNKNOWN"
		} else if !shouldBreak {
			tmp.Code = _err.Code.String()
			tmp.Message = utils.NonEmptyOr(_err.Message.RoleMessage(role).LocaleText(orderedLocalePreference...), _err.Message.RoleMessage("COMMON").LocaleText(orderedLocalePreference...))
			tmp.Namespace = _err.Namespace.String()
		}

		if prevErr != nil &&
			len(prevErr.Code) > 1 &&
			len(tmp.Code) > 1 &&
			prevErr.Code.String()[:2] == tmp.Code[:2] {
			errRes[len(errRes)-1].Elaborations = append(errRes[len(errRes)-1].Elaborations, tmp)
		} else {
			errRes = append(errRes, tmp)
		}

		if shouldBreak {
			break
		}

		prevErr = _err
	}

	return
}

func ObtainErrorStatus(err error) int {
	switch true {
	case errors.Is(err, apiValues.ApiErrorUnauthenticatedAccess):
		return http.StatusUnauthorized
	case errors.Is(err, apiValues.ApiErrorFailedValidation):
		return http.StatusBadRequest
	case errors.Is(err, apiValues.ApiErrorNotFound):
		return http.StatusNotFound
	case errors.Is(err, apiValues.ApiErrorUnauthorizedAction),
		errors.Is(err, apiValues.ApiErrorForbiddenAction):
		return http.StatusForbidden
	case errors.Is(err, apiValues.ApiErrorServiceUnavailable):
		return http.StatusServiceUnavailable
	case errors.Is(err, apiValues.ApiErrorDuplicatePost):
		return http.StatusConflict
	case errors.Is(err, apiValues.ApiErrorUnprocessableEntity):
		return http.StatusUnprocessableEntity
	case errors.Is(err, apiValues.ApiErrorNotImplemented):
		return http.StatusNotImplemented
	default:
		return http.StatusInternalServerError
	}
}

func GetErrorRole(r *http.Request) commonModel.ErrorRole {

	var (
		_default  = commonModel.ErrorRole("DEFAULT")
		ctx, _err = GetWrappedContext(r)
	)

	if _err != nil {
		return _default
	}

	var role = utils.SafeCastValue[string](ctx.TagAlong.Data["errorRole"])

	if role != nil {
		errRole := utils.Ternary(*role != "", commonModel.ErrorRole(*role), _default)
		return errRole
	}
	return _default
}

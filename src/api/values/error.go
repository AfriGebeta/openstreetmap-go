package values

import (
	"errors"
	commonModel "openstreetmap-go/src/common/model"
	commonUtils "openstreetmap-go/src/common/utils"
	commonValues "openstreetmap-go/src/common/values"
)

var (
	ErrorCodeApiUnknown                    commonModel.ErrorCode = "AE00001"
	ErrorCodeApiInternalServerError        commonModel.ErrorCode = "AE00002"
	ErrorCodeApiInternalServerConflict     commonModel.ErrorCode = "AE00003"
	ErrorCodeApiNotFound                   commonModel.ErrorCode = "AE00004"
	ErrorCodeApiNotImplemented             commonModel.ErrorCode = "AE00005"
	ErrorCodeApiDuplicatePost              commonModel.ErrorCode = "AE00006"
	ErrorCodeApiUnprocessableEntity        commonModel.ErrorCode = "AE00007"
	ErrorCodeApiFailedValidation           commonModel.ErrorCode = "AE00008"
	ErrorCodeApiUnauthenticatedAccess      commonModel.ErrorCode = "AE00009"
	ErrorCodeApiUnauthorizedAction         commonModel.ErrorCode = "AE0000A"
	ErrorCodeApiForbiddenAction            commonModel.ErrorCode = "AE0000B"
	ErrorCodeApiInternalServiceUnavailable commonModel.ErrorCode = "AE0000C"
)

var (
	ApiErrorUnknown                = errors.New(ErrorCodeApiUnknown.String())
	ApiErrorInternalServerError    = errors.New(ErrorCodeApiInternalServerError.String())
	ApiErrorInternalServerConflict = errors.New(ErrorCodeApiInternalServerConflict.String())
	ApiErrorNotFound               = errors.New(ErrorCodeApiNotFound.String())
	ApiErrorNotImplemented         = errors.New(ErrorCodeApiNotImplemented.String())
	ApiErrorDuplicatePost          = errors.New(ErrorCodeApiDuplicatePost.String())
	ApiErrorFailedValidation       = errors.New(ErrorCodeApiFailedValidation.String())
	ApiErrorUnprocessableEntity    = errors.New(ErrorCodeApiUnprocessableEntity.String())
	ApiErrorUnauthenticatedAccess  = errors.New(ErrorCodeApiUnauthenticatedAccess.String())
	ApiErrorUnauthorizedAction     = errors.New(ErrorCodeApiUnauthorizedAction.String())
	ApiErrorForbiddenAction        = errors.New(ErrorCodeApiForbiddenAction.String())
	ApiErrorServiceUnavailable     = errors.New(ErrorCodeApiInternalServiceUnavailable.String())
)

var ApiErrors = commonModel.ErrorList{
	ErrorCodeApiUnknown: commonUtils.GetCommonError(
		commonValues.ErrorCodeCommonUnknown,
		commonModel.ErrorLayerValue.Api,
	),
	ErrorCodeApiInternalServerError: {
		Layer:     commonModel.ErrorLayerValue.Api,
		Namespace: "DEFAULT",
		Code:      ErrorCodeApiInternalServerError,
		Message: commonModel.ErrorMessage{
			"DEFAULT": commonModel.IntlText{
				commonModel.LocaleValue.EnUs: "Something went wrong",
				commonModel.LocaleValue.AmEt: "ችግር ተፈጥራል",
			},
		},
		RoleVisibility: map[commonModel.ErrorRole]bool{
			"DEFAULT": true,
		},
	},
	ErrorCodeApiInternalServerConflict: {
		Layer:     commonModel.ErrorLayerValue.Api,
		Namespace: "DEFAULT",
		Code:      ErrorCodeApiInternalServerConflict,
		Message: commonModel.ErrorMessage{
			"DEFAULT": commonModel.IntlText{
				commonModel.LocaleValue.EnUs: "Internal server conflict",
				commonModel.LocaleValue.AmEt: "ውስጣዊ ግጭት",
			},
		},
	},
	ErrorCodeApiNotFound: {
		Layer:     commonModel.ErrorLayerValue.Api,
		Namespace: "DEFAULT",
		Code:      ErrorCodeApiNotFound,
		Message: commonModel.ErrorMessage{
			"DEFAULT": commonModel.IntlText{
				commonModel.LocaleValue.EnUs: "Resource not found",
				commonModel.LocaleValue.AmEt: "አልተገኘም",
			},
		},
	},
	ErrorCodeApiNotImplemented: {
		Layer:     commonModel.ErrorLayerValue.Api,
		Namespace: "DEFAULT",
		Code:      ErrorCodeApiNotImplemented,
		Message: commonModel.ErrorMessage{
			"DEFAULT": commonModel.IntlText{
				commonModel.LocaleValue.EnUs: "Not implemented",
				commonModel.LocaleValue.AmEt: "አልተተገበረም",
			},
		},
	},
	ErrorCodeApiDuplicatePost: {
		Layer:     commonModel.ErrorLayerValue.Api,
		Namespace: "DEFAULT",
		Code:      ErrorCodeApiDuplicatePost,
		Message: commonModel.ErrorMessage{
			"DEFAULT": commonModel.IntlText{
				commonModel.LocaleValue.EnUs: "Resource already exists",
				commonModel.LocaleValue.AmEt: "ከዚህ ቀደም የገባ",
			},
		},
	},
	ErrorCodeApiFailedValidation: {
		Layer:     commonModel.ErrorLayerValue.Api,
		Namespace: "DEFAULT",
		Code:      ErrorCodeApiFailedValidation,
		Message: commonModel.ErrorMessage{
			"DEFAULT": commonModel.IntlText{
				commonModel.LocaleValue.EnUs: "Data validation failed",
				commonModel.LocaleValue.AmEt: "የውሂብ ማረጋገጫ አልተሳካም",
			},
		},
	},
	ErrorCodeApiUnprocessableEntity: {
		Layer:     commonModel.ErrorLayerValue.Api,
		Namespace: "DEFAULT",
		Code:      ErrorCodeApiUnprocessableEntity,
		Message: commonModel.ErrorMessage{
			"DEFAULT": commonModel.IntlText{
				commonModel.LocaleValue.EnUs: "Unprocessable entity",
				commonModel.LocaleValue.AmEt: "አልተሳካም ንብረት",
			},
		},
	},
	ErrorCodeApiUnauthenticatedAccess: {
		Layer:     commonModel.ErrorLayerValue.Api,
		Namespace: "DEFAULT",
		Code:      ErrorCodeApiUnauthenticatedAccess,
		Message: commonModel.ErrorMessage{
			"DEFAULT": commonModel.IntlText{
				commonModel.LocaleValue.EnUs: "Unauthenticated access",
				commonModel.LocaleValue.AmEt: "የማንነት ማረጋገጫ አልተሳካም",
			},
		},
	},
	ErrorCodeApiUnauthorizedAction: {
		Layer:     commonModel.ErrorLayerValue.Api,
		Namespace: "DEFAULT",
		Code:      ErrorCodeApiUnauthorizedAction,
		Message: commonModel.ErrorMessage{
			"DEFAULT": commonModel.IntlText{
				commonModel.LocaleValue.EnUs: "Unauthorized action",
				commonModel.LocaleValue.AmEt: "ያልተፈቀደ ተግባር",
			},
		},
	},
	ErrorCodeApiForbiddenAction: {
		Layer:     commonModel.ErrorLayerValue.Api,
		Namespace: "DEFAULT",
		Code:      ErrorCodeApiForbiddenAction,
		Message: commonModel.ErrorMessage{
			"DEFAULT": commonModel.IntlText{
				commonModel.LocaleValue.EnUs: "Forbidden action",
				commonModel.LocaleValue.AmEt: "የተከለከለ ተግባር",
			},
		},
	},
	ErrorCodeApiInternalServiceUnavailable: {
		Layer:     commonModel.ErrorLayerValue.Api,
		Namespace: "DEFAULT",
		Code:      ErrorCodeApiInternalServiceUnavailable,
		Message: commonModel.ErrorMessage{
			"DEFAULT": commonModel.IntlText{
				commonModel.LocaleValue.EnUs: "Service unavailable",
				commonModel.LocaleValue.AmEt: "የቆመ አገልግሎት",
			},
		},
	},
}

func init() {
	for key, value := range ApiErrors {
		commonValues.ErrorList[key] = value
	}
}

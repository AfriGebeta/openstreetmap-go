package values

import (
	"errors"
	commonModel "openstreetmap-go/src/common/model"
	"openstreetmap-go/src/common/utils"
	commonValues "openstreetmap-go/src/common/values"
)

var (
	ErrorCodeDomainUnknown          commonModel.ErrorCode = "DE00001"
	ErrorCodeDomainForbidden        commonModel.ErrorCode = "DE00002"
	ErrorCodeDomainFailedValidation commonModel.ErrorCode = "DE00003"
	ErrorCodeDomainNotFound         commonModel.ErrorCode = "DE00004"
)

var (
	DomainErrorUnknown          = errors.New(ErrorCodeDomainUnknown.String())
	DomainErrorForbidden        = errors.New(ErrorCodeDomainForbidden.String())
	DomainErrorFailedValidation = errors.New(ErrorCodeDomainFailedValidation.String())
	DomainErrorNotFound         = errors.New(ErrorCodeDomainNotFound.String())
)

var DomainErrorList = commonModel.ErrorList{
	ErrorCodeDomainUnknown: utils.GetCommonError(
		commonValues.ErrorCodeCommonUnknown,
		commonModel.ErrorLayerValue.Domain,
	),
	ErrorCodeDomainForbidden: {
		Layer:     commonModel.ErrorLayerValue.Domain,
		Namespace: "DEFAULT",
		Code:      ErrorCodeDomainForbidden,
		Message: commonModel.ErrorMessage{
			"DEFAULT": commonModel.IntlText{
				commonModel.LocaleValue.EnUs: "Forbidden action",
				commonModel.LocaleValue.AmEt: "የተከለከለ ተግባር",
			},
		},
	},
	ErrorCodeDomainFailedValidation: {
		Layer:     commonModel.ErrorLayerValue.Domain,
		Namespace: "DEFAULT",
		Code:      ErrorCodeDomainFailedValidation,
		Message: commonModel.ErrorMessage{
			"DEFAULT": commonModel.IntlText{
				commonModel.LocaleValue.EnUs: "Data validation failed",
				commonModel.LocaleValue.AmEt: "የተሳሳተ የውሂብ ቅርጸት",
			},
		},
	},
	ErrorCodeDomainNotFound: {
		Layer:     commonModel.ErrorLayerValue.Domain,
		Namespace: "DEFAULT",
		Code:      ErrorCodeDomainNotFound,
		Message: commonModel.ErrorMessage{
			"DEFAULT": commonModel.IntlText{
				commonModel.LocaleValue.EnUs: "Resource Not found",
				commonModel.LocaleValue.AmEt: "ውሂቡ አልተገኘም",
			},
		},
	},
}

func init() {
	for _, v := range DomainErrorList {
		commonValues.ErrorList[v.Code] = v
	}
}

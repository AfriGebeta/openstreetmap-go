package values

import (
	"errors"
	commonModel "openstreetmap-go/src/common/model"
	commonUtils "openstreetmap-go/src/common/utils"
	commonValues "openstreetmap-go/src/common/values"
)

var (
	ErrorCodeDataUnknown                  commonModel.ErrorCode = "RE00001"
	ErrorCodeDataNotFound                 commonModel.ErrorCode = "RE00002"
	ErrorCodeDataInvalid                  commonModel.ErrorCode = "RE00003"
	ErrorCodeDataDuplicate                commonModel.ErrorCode = "RE00004"
	ErrorCodeDataMissingRequiredField     commonModel.ErrorCode = "RE00005"
	ErrorCodeDataReferencedRecordNotFound commonModel.ErrorCode = "RE00006"
)

var (
	DataErrorUnknown                  = errors.New(ErrorCodeDataUnknown.String())
	DataErrorNotFound                 = errors.New(ErrorCodeDataNotFound.String())
	DataErrorInvalid                  = errors.New(ErrorCodeDataInvalid.String())
	DataErrorDuplicate                = errors.New(ErrorCodeDataDuplicate.String())
	DataErrorMissingRequiredField     = errors.New(ErrorCodeDataMissingRequiredField.String())
	DataErrorReferencedRecordNotFound = errors.New(ErrorCodeDataReferencedRecordNotFound.String())
)

var DataErrors = commonModel.ErrorList{
	ErrorCodeDataUnknown: commonUtils.GetCommonError(
		commonValues.ErrorCodeCommonUnknown,
		commonModel.ErrorLayerValue.Data,
	),
	ErrorCodeDataNotFound: {
		Layer:     commonModel.ErrorLayerValue.Data,
		Namespace: "DEFAULT",
		Code:      ErrorCodeDataNotFound,
		Message: commonModel.ErrorMessage{
			"DEFAULT": commonModel.IntlText{
				commonModel.LocaleValue.EnUs: "Record not found",
				commonModel.LocaleValue.AmEt: "መዝገቡ አልተገኘም",
			},
		},
	},
	ErrorCodeDataInvalid: {
		Layer:     commonModel.ErrorLayerValue.Data,
		Namespace: "DEFAULT",
		Code:      ErrorCodeDataInvalid,
		Message: commonModel.ErrorMessage{
			"DEFAULT": commonModel.IntlText{
				commonModel.LocaleValue.EnUs: "Invalid data format",
				commonModel.LocaleValue.AmEt: "የማይታወቅ ውሂብ ፎርማት",
			},
		},
	},
	ErrorCodeDataDuplicate: {
		Layer:     commonModel.ErrorLayerValue.Data,
		Namespace: "DEFAULT",
		Code:      ErrorCodeDataDuplicate,
		Message: commonModel.ErrorMessage{
			"DEFAULT": commonModel.IntlText{
				commonModel.LocaleValue.EnUs: "Duplicate data",
				commonModel.LocaleValue.AmEt: "የተደገመ ውሂብ",
			},
		},
	},
	ErrorCodeDataMissingRequiredField: {
		Layer:     commonModel.ErrorLayerValue.Data,
		Namespace: "DEFAULT",
		Code:      ErrorCodeDataMissingRequiredField,
		Message: commonModel.ErrorMessage{
			"DEFAULT": commonModel.IntlText{
				commonModel.LocaleValue.EnUs: "Missing required field",
				commonModel.LocaleValue.AmEt: "ውሂቡ አልተማላም",
			},
		},
	},
	ErrorCodeDataReferencedRecordNotFound: {
		Layer:     commonModel.ErrorLayerValue.Data,
		Namespace: "DEFAULT",
		Code:      ErrorCodeDataReferencedRecordNotFound,
		Message: commonModel.ErrorMessage{
			"DEFAULT": commonModel.IntlText{
				commonModel.LocaleValue.EnUs: "Referenced record not found",
				commonModel.LocaleValue.AmEt: "የተጠቆመበት መዝገብ አልተገኘም",
			},
		},
	},
}

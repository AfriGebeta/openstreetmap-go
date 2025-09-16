package values

import (
	"errors"
	"openstreetmap-go/src/common/model"
)

var (
	ErrorCodeCommonUnknown = model.ErrorCode("SE00001")
)

var (
	CommonErrorUnknown = errors.New(ErrorCodeCommonUnknown.String())
)

var CommonErrors = map[model.ErrorCode]*model.Error{
	ErrorCodeCommonUnknown: {
		Layer:     model.ErrorLayerValue.System,
		Namespace: "DEFAULT",
		Code:      ErrorCodeCommonUnknown,
		Message: map[model.ErrorRole]model.IntlText{
			"DEFAULT": {
				model.LocaleValue.EnUs: "Unknown error",
				model.LocaleValue.AmEt: "ያልታወቀ ድክመት",
			},
		},
		RoleVisibility: map[model.ErrorRole]bool{
			"DEFAULT": true,
		},
	},
}

var ErrorList = model.ErrorList{}

func init() {
	for _, v := range CommonErrors {
		ErrorList[v.Code] = v
	}
}

package utils

import (
	"openstreetmap-go/src/common/model"
	"openstreetmap-go/src/common/values"
	"openstreetmap-go/src/utils"
	"strings"
)

func GetCommonError(code model.ErrorCode, layer ...model.ErrorLayer) *model.Error {
	var err = values.CommonErrors[code].Clone()

	if len(layer) > 0 {
		err.Layer = layer[0]
	}

	return err
}

func GetPrintableError(tag string, err error, role model.ErrorRole, orderedLocalePreference ...model.Locale) string {
	return GetPrintableErrorFromErrorList(tag, utils.RecursivelyUnwrapJoin(err), role, orderedLocalePreference...)
}

func GetPrintableErrorFromErrorList(tag string, errList []error, role model.ErrorRole, orderedLocalePreference ...model.Locale) string {
	var printable string

	for i, e := range errList {
		var _err, ok = values.ErrorList[model.ErrorCode(e.Error())]

		if i > 0 {
			printable += "\n" + strings.Repeat(" ", len(tag)+21)
		}

		if !ok {
			printable += "(" + utils.GreyString("UNKNOWN") + ") " + e.Error()
		} else {
			printable += "(" + utils.GreyString(_err.Code.String()) + ") " + _err.Message.RoleMessage(role).LocaleText(orderedLocalePreference...)
		}
	}

	return printable
}

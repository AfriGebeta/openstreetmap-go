package utils

import (
	"fmt"
	"html/template"
	"openstreetmap-go/src/utils"
	"strconv"
	"strings"
	"time"
)

var TemplateFunctions = template.FuncMap{
	"formatDate": func(date time.Time) string {
		return date.Format("January 2, 2006")
	},
	"formatDateTime": func(date time.Time) string {
		return date.Format("January 2, 2006 at 15:04:05")
	},
	"formatCurrencyValue64": func(amount float64) string {
		return strconv.FormatFloat(amount, 'f', 2, 64)
	},
	"formatCurrencyValue32": func(amount float32) string {
		return strconv.FormatFloat(float64(amount), 'f', 2, 64)
	},
	"toUpper":   strings.ToUpper,
	"toLower":   strings.ToLower,
	"trimSpace": strings.TrimSpace,
	"sprintf":   fmt.Sprintf,
	"valueOrString": func(value *string, replacement string) string {
		return utils.ValueOr(value, replacement)
	},
}

package model

import (
	"openstreetmap-go/src/utils"
	"strings"
)

// IntlText is a map of locale to string. It is used to store internationalized text.
// The key __accepts__ is a reserved key that is used to store a comma separated list
// of keys that serve as placeholders in the string. These keys should be all small-
// letter, snake-case strings.
type IntlText map[Locale]string

// LocaleText returns the text in the given locale. If the locale is not found, it
// falls back to en-us.
func (it IntlText) LocaleText(localePreferenceOrder ...Locale) string {
	var tmp = it[LocaleValue.EnUs]

	for _, locale := range localePreferenceOrder {
		if msg, ok := it[Locale(strings.ToLower(locale.String()))]; ok {
			return msg
		}
	}

	return tmp
}

// LocaleConstructedText constructs a string from the template in the given locale
// using the given arguments.
func (it IntlText) LocaleConstructedText(args map[string]string, localePreferenceOrder ...Locale) string {
	return utils.ConstructStringFromTemplate(
		it.LocaleText(localePreferenceOrder...),
		args,
	)
}

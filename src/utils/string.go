package utils

import "strings"

func GetEquivalentWhiteSpace(string string) string {
	var tmp = ""

	for range string {
		tmp = tmp + " "
	}

	return tmp
}

func ConstructStringFromTemplate(template string, args map[string]string) string {
	for key, value := range args {
		template = strings.ReplaceAll(template, "{{"+key+"}}", value)
	}

	return template
}

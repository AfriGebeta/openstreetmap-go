package model

type Locale string

type localeValue struct {
	AmEt Locale `json:"am-ET" ,json:"am"`
	EnUs Locale `json:"en-US" ,json:"en"`
}

var LocaleValue = localeValue{
	AmEt: "am-ET",
	EnUs: "en-US",
}

var LocaleValues = []Locale{
	LocaleValue.AmEt,
	LocaleValue.EnUs,
}

func (l Locale) String() string {
	return string(l)
}

func (l Locale) StringPtr() *string {
	var tmp = l.String()
	return &tmp
}

package values

type CachePrefix string

var CachePrefixValue = struct {
	AuthorizationRequest CachePrefix
	Otp                  CachePrefix
	Authentication       CachePrefix
	Account              CachePrefix
}{
	AuthorizationRequest: "authorization_request",
	Otp:                  "otp",
	Authentication:       "authorization_code",
	Account:              "account",
}

var CachePrefixValues = []CachePrefix{
	CachePrefixValue.AuthorizationRequest,
	CachePrefixValue.Otp,
	CachePrefixValue.Authentication,
}

func (c CachePrefix) String() string {
	return string(c)
}

func (c CachePrefix) StringPtr() *string {
	var tmp = c.String()
	return &tmp
}

func (c CachePrefix) SurroundWithColons() string {
	return ":" + c.String() + ":"
}

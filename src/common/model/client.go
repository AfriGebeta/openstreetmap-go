package model

type ClientPlatform string

var ClientPlatformValue = struct {
	Web     ClientPlatform
	Android ClientPlatform
	Ios     ClientPlatform
	Windows ClientPlatform
	Darwin  ClientPlatform
	Linux   ClientPlatform
	Unknown ClientPlatform
}{
	Web:     ClientPlatform("WEB"),
	Android: ClientPlatform("ANDROID"),
	Ios:     ClientPlatform("IOS"),
	Windows: ClientPlatform("WINDOWS"),
	Darwin:  ClientPlatform("DARWIN"),
	Linux:   ClientPlatform("LINUX"),
	Unknown: ClientPlatform("UNKNOWN"),
}

var ClientPlatformValues = []ClientPlatform{
	ClientPlatform("WEB"),
	ClientPlatform("ANDROID"),
	ClientPlatform("IOS"),
	ClientPlatform("WINDOWS"),
	ClientPlatform("DARWIN"),
	ClientPlatform("LINUX"),
	ClientPlatform("UNKNOWN"),
}

func (e ClientPlatform) String() string { return string(e) }

func (e ClientPlatform) StringPtr() *string {
	var tmp = e.String()
	return &tmp
}

package config

// these are and should always be injected by the build system
var (
	Version    = "v0.0.1"
	CommitHash = ""
)

var ServerConfigObject = ServerConfig{
	EnvFilePath:              "deployment-files/local/.env",
	YmlFilePath:              "deployment-files/local/settings.yml",
	Env:                      "local",
	SelfBaseUrl:              "http://0.0.0.0:5050",
	DefaultDataPublic:        true,
	DefaultHomeLat:           0.0,
	DefaultHomeLon:           0.0,
	DefaultHomeZoom:          0,
	DefaultLanguages:         "en-GB,en-US,en",
	DefaultConsiderPd:        false,
	DefaultDescriptionFormat: "markdown",
	DefaultImageUseGravatar:  false,
	SecretKeyBase:            "",
	GeoRecordScale:           10000000,
	MaxRequestArea:           0.25,
	MaxNumberOfWayNodes:      50000,
}

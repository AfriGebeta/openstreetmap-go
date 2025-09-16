package config

// ServerConfig is a configuration for the technical operations of
// the server.

type OSMSettings struct {
	ServerProtocol                    string   `yaml:"server_protocol"`
	ServerURL                         string   `yaml:"server_url"`
	Generator                         string   `yaml:"generator"`
	CopyrightOwner                    string   `yaml:"copyright_owner"`
	AttributionURL                    string   `yaml:"attribution_url"`
	LicenseURL                        string   `yaml:"license_url"`
	SupportEmail                      string   `yaml:"support_email"`
	EmailFrom                         string   `yaml:"email_from"`
	EmailReturnPath                   string   `yaml:"email_return_path"`
	APIVersion                        string   `yaml:"api_version"`
	Status                            string   `yaml:"status"`
	MaxRequestArea                    float64  `yaml:"max_request_area"`
	TracepointsPerPage                int      `yaml:"tracepoints_per_page"`
	DefaultChangesetQueryLimit        int      `yaml:"default_changeset_query_limit"`
	MaxChangesetQueryLimit            int      `yaml:"max_changeset_query_limit"`
	DefaultChangesetCommentQueryLimit int      `yaml:"default_changeset_comment_query_limit"`
	MaxChangesetCommentQueryLimit     int      `yaml:"max_changeset_comment_query_limit"`
	DefaultNoteQueryLimit             int      `yaml:"default_note_query_limit"`
	MaxNoteQueryLimit                 int      `yaml:"max_note_query_limit"`
	MaxNumberOfNodes                  int      `yaml:"max_number_of_nodes"`
	MaxNumberOfWayNodes               int      `yaml:"max_number_of_way_nodes"`
	MaxNumberOfRelationMembers        int      `yaml:"max_number_of_relation_members"`
	PostcodeZoom                      int      `yaml:"postcode_zoom"`
	APITimeout                        int      `yaml:"api_timeout"`
	WebTimeout                        int      `yaml:"web_timeout"`
	DefaultEditor                     string   `yaml:"default_editor"`
	NominatimURL                      string   `yaml:"nominatim_url"`
	GraphhopperURL                    string   `yaml:"graphhopper_url"`
	FOSSGISOSRMURL                    string   `yaml:"fossgis_osrm_url"`
	FOSSGISValhallaURL                string   `yaml:"fossgis_valhalla_url"`
	WikidataAPIURL                    string   `yaml:"wikidata_api_url"`
	WikimediaCommonsURL               string   `yaml:"wikimedia_commons_url"`
	LinkifyHosts                      []string `yaml:"linkify_hosts"`
	LinkifyHostsReplacement           string   `yaml:"linkify_hosts_replacement"`
}

type ServerConfig struct {
	EnvFilePath              string
	YmlFilePath              string
	Env                      string
	Port                     string
	PrimaryDbUrl             string
	ReplicaDbUrls            []string
	CacheUrl                 string
	CachePassword            string
	CacheDatabase            int
	CacheSystemPrefix        string
	AllowedOrigins           []string
	SelfBaseUrl              string
	JwtSecret                string
	OSMSettings              *OSMSettings
	DefaultDataPublic        bool
	DefaultHomeLat           float64
	DefaultHomeLon           float64
	DefaultHomeZoom          int
	DefaultLanguages         string
	DefaultConsiderPd        bool
	DefaultDescriptionFormat string
	DefaultImageUseGravatar  bool
	SecretKeyBase            string
	GeoRecordScale           int
	MaxRequestArea           float64
	MaxNumberOfWayNodes      int
}

package config

import (
	"errors"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path"
	"strconv"
)

var config ServerConfig

func SetupEnv() {
	SetupConfigurationFromEnv()
	SetupConfigurationFromYaml()
}

func SetupConfigurationFromEnv() {

	if envLoadErr := godotenv.Load(); envLoadErr != nil {
		panic("could not load .env file")
	}
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		ServerConfiguration.DBHost = dbHost
	} else {
		panic(errors.New("DB_HOST environment variable not set"))
	}
	if dbUsername := os.Getenv("DB_USERNAME"); dbUsername != "" {
		ServerConfiguration.DBUserName = dbUsername
	} else {
		panic(errors.New("DB_USERNAME environment variable not set"))
	}
	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		ServerConfiguration.DBPassword = dbPassword
	} else {
		panic(errors.New("DB_PASSWORD environment variable not set"))
	}

	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		port, err := strconv.Atoi(dbPort)
		if err != nil {
			panic("Invalid DB_PORT")
		}
		ServerConfiguration.DBPort = port
	} else {
		panic(errors.New("DB_PORT environment variable not set"))
	}

	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		ServerConfiguration.DatabaseName = dbName
	} else {
		panic(errors.New("DB_NAME environment variable not set"))
	}

}

// read from yaml
func SetupConfigurationFromYaml() {
	getwd, err := os.Getwd()
	if err != nil {
		return
	}
	ymlPath := path.Join(getwd, "/src/config/settings.yml")
	f, err := os.Open(ymlPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	ServerConfiguration.ServerProtocol = config.ServerProtocol
	ServerConfiguration.ServerUrl = config.ServerUrl
	ServerConfiguration.Generator = config.Generator
	ServerConfiguration.CopyrightOwner = config.CopyrightOwner
	ServerConfiguration.AttributionUrl = config.AttributionUrl
	ServerConfiguration.LicenseUrl = config.LicenseUrl
	ServerConfiguration.SupportEmail = config.SupportEmail
	ServerConfiguration.EmailFrom = config.EmailFrom
	ServerConfiguration.EmailReturnPath = config.EmailReturnPath
	ServerConfiguration.ApiVersion = config.ApiVersion
	ServerConfiguration.Status = config.Status
	ServerConfiguration.MaxRequestArea = config.MaxRequestArea
	ServerConfiguration.TracepointsPerPage = config.TracepointsPerPage
	ServerConfiguration.DefaultChangesetQueryLimit = config.DefaultChangesetQueryLimit
	ServerConfiguration.MaxChangesetQueryLimit = config.MaxChangesetQueryLimit
	ServerConfiguration.DefaultChangesetCommentQueryLimit = config.DefaultChangesetCommentQueryLimit
	ServerConfiguration.MaxChangesetCommentQueryLimit = config.MaxChangesetCommentQueryLimit
	ServerConfiguration.DefaultChangesetCommentsFeedQueryLimit = config.DefaultChangesetCommentsFeedQueryLimit
	ServerConfiguration.MaxChangesetCommentsFeedQueryLimit = config.MaxChangesetCommentsFeedQueryLimit
	ServerConfiguration.MaxNumberOfNodes = config.MaxNumberOfNodes
	ServerConfiguration.MaxNumberOfWayNodes = config.MaxNumberOfWayNodes
	ServerConfiguration.MaxNumberOfRelationMembers = config.MaxNumberOfRelationMembers
	ServerConfiguration.MaxNoteRequestArea = config.MaxNoteRequestArea
	ServerConfiguration.DefaultNoteQueryLimit = config.DefaultNoteQueryLimit
	ServerConfiguration.MaxNoteQueryLimit = config.MaxNoteQueryLimit
	ServerConfiguration.MaxIssuesCount = config.MaxIssuesCount
	ServerConfiguration.MaxTraceSize = config.MaxTraceSize
	ServerConfiguration.PostcodeZoom = config.PostcodeZoom
	ServerConfiguration.ApiTimeout = config.ApiTimeout
	ServerConfiguration.WebTimeout = config.WebTimeout
	ServerConfiguration.UserBlockPeriods = config.UserBlockPeriods
	ServerConfiguration.UserAccountDeletionDelay = config.UserAccountDeletionDelay
	ServerConfiguration.MaxMessagesPerHour = config.MaxMessagesPerHour
	ServerConfiguration.DefaultMessageQueryLimit = config.DefaultMessageQueryLimit
	ServerConfiguration.MaxMessageQueryLimit = config.MaxMessageQueryLimit
	ServerConfiguration.MaxFollowsPerHour = config.MaxFollowsPerHour
	ServerConfiguration.MinChangesetCommentsPerHour = config.MinChangesetCommentsPerHour
	ServerConfiguration.InitialChangesetCommentsPerHour = config.InitialChangesetCommentsPerHour
	ServerConfiguration.MaxChangesetCommentsPerHour = config.MaxChangesetCommentsPerHour
	ServerConfiguration.CommentsToMaxChangesetComments = config.CommentsToMaxChangesetComments
	ServerConfiguration.ModeratorChangesetCommentsPerHour = config.ModeratorChangesetCommentsPerHour
	ServerConfiguration.MinChangesPerHour = config.MinChangesPerHour
	ServerConfiguration.InitialChangesPerHour = config.InitialChangesPerHour
	ServerConfiguration.MaxChangesPerHour = config.MaxChangesPerHour
	ServerConfiguration.DaysToMaxChanges = config.DaysToMaxChanges
	ServerConfiguration.ImporterChangesPerHour = config.ImporterChangesPerHour
	ServerConfiguration.ModeratorChangesPerHour = config.ModeratorChangesPerHour
	ServerConfiguration.MinSizeLimit = config.MinSizeLimit
	ServerConfiguration.InitialSizeLimit = config.InitialSizeLimit
	ServerConfiguration.MaxSizeLimit = config.MaxSizeLimit
	ServerConfiguration.DaysToMaxSizeLimit = config.DaysToMaxSizeLimit
	ServerConfiguration.ImporterSizeLimit = config.ImporterSizeLimit
	ServerConfiguration.ModeratorSizeLimit = config.ModeratorSizeLimit
	ServerConfiguration.MessagesDomain = config.MessagesDomain
	ServerConfiguration.MaxmindDatabase = config.MaxmindDatabase
	ServerConfiguration.NearbyUsers = config.NearbyUsers
	ServerConfiguration.NearbyRadius = config.NearbyRadius
	ServerConfiguration.SpamThreshold = config.SpamThreshold
	ServerConfiguration.DiaryFeedDelay = config.DiaryFeedDelay
	ServerConfiguration.DefaultLegale = config.DefaultLegale
	ServerConfiguration.AttachmentsDir = config.AttachmentsDir
	ServerConfiguration.LogPath = config.LogPath
	ServerConfiguration.LogstashPath = config.LogstashPath
	ServerConfiguration.MaptilerKey = config.MaptilerKey
	ServerConfiguration.MemcacheServers = config.MemcacheServers
	ServerConfiguration.NominatimUrl = config.NominatimUrl
	ServerConfiguration.DefaultEditor = config.DefaultEditor
	ServerConfiguration.OauthApplication = config.OauthApplication
	ServerConfiguration.IdApplication = config.IdApplication
	ServerConfiguration.ImageryBlacklist = config.ImageryBlacklist
	ServerConfiguration.OverpassUrl = config.OverpassUrl
	ServerConfiguration.OverpassCredentials = config.OverpassCredentials
	ServerConfiguration.GraphhopperUrl = config.GraphhopperUrl
	ServerConfiguration.FossgisOsrmUrl = config.FossgisOsrmUrl
	ServerConfiguration.FossgisValhallaUrl = config.FossgisValhallaUrl
	ServerConfiguration.WikidataApiUrl = config.WikidataApiUrl
	ServerConfiguration.WikimediaCommonsUrl = config.WikimediaCommonsUrl
	ServerConfiguration.LinkifyHosts = config.LinkifyHosts
	ServerConfiguration.LinkifyHostsReplacement = config.LinkifyHostsReplacement
	ServerConfiguration.LinkifyWikiHosts = config.LinkifyWikiHosts
	ServerConfiguration.LinkifyWikiHostsReplacement = config.LinkifyWikiHostsReplacement
	ServerConfiguration.LinkifyWikiOptionalPathPrefix = config.LinkifyWikiOptionalPathPrefix
	ServerConfiguration.GoogleAuthId = config.GoogleAuthId
	ServerConfiguration.GoogleAuthSecret = config.GoogleAuthSecret
	ServerConfiguration.GoogleOpenidRealm = config.GoogleOpenidRealm
	ServerConfiguration.FacebookAuthId = config.FacebookAuthId
	ServerConfiguration.FacebookAuthSecret = config.FacebookAuthSecret
	ServerConfiguration.GithubAuthId = config.GithubAuthId
	ServerConfiguration.GithubAuthSecret = config.GithubAuthSecret
	ServerConfiguration.MicrosoftAuthId = config.MicrosoftAuthId
	ServerConfiguration.MicrosoftAuthSecret = config.MicrosoftAuthSecret
	ServerConfiguration.WikipediaAuthId = config.WikipediaAuthId
	ServerConfiguration.WikipediaAuthSecret = config.WikipediaAuthSecret
	ServerConfiguration.ThunderforestKey = config.ThunderforestKey
	ServerConfiguration.TracestrackKey = config.TracestrackKey
	ServerConfiguration.TotpKey = config.TotpKey
	ServerConfiguration.CspEnforce = config.CspEnforce
	ServerConfiguration.CspReportUrl = config.CspReportUrl
	ServerConfiguration.AvatarStorage = config.AvatarStorage
	ServerConfiguration.TraceFileStorage = config.TraceFileStorage
	ServerConfiguration.TraceImageStorage = config.TraceImageStorage
	ServerConfiguration.TraceIconStorage = config.TraceIconStorage
	ServerConfiguration.AvatarStorageUrl = config.AvatarStorageUrl
	ServerConfiguration.TraceImageStorageUrl = config.TraceImageStorageUrl
	ServerConfiguration.TraceIconStorageUrl = config.TraceIconStorageUrl
	ServerConfiguration.TileCdnUrl = config.TileCdnUrl
	ServerConfiguration.SmtpAddress = config.SmtpAddress
	ServerConfiguration.SmtpPort = config.SmtpPort
	ServerConfiguration.SmtpDomain = config.SmtpDomain
	ServerConfiguration.SmtpEnableStarttlsAuto = config.SmtpEnableStarttlsAuto
	ServerConfiguration.SmtpTlsVerifyMode = config.SmtpTlsVerifyMode
	ServerConfiguration.SmtpAuthentication = config.SmtpAuthentication
	ServerConfiguration.SmtpUserName = config.SmtpUserName
	ServerConfiguration.SmtpPassword = config.SmtpPassword
	ServerConfiguration.Matomo = config.Matomo
	ServerConfiguration.SignupIpPerDay = config.SignupIpPerDay
	ServerConfiguration.SignupIpMaxBurst = config.SignupIpMaxBurst
	ServerConfiguration.SignupEmailPerDay = config.SignupEmailPerDay
	ServerConfiguration.SignupEmailMaxBurst = config.SignupIpMaxBurst
	ServerConfiguration.DoorkeeperSigningKey = config.DoorkeeperSigningKey

}

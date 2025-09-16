package config

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
	"openstreetmap-go/src/utils"
	"os"
	"strconv"
	"strings"
)

func Load(configObject *ServerConfig) {
	LoadArgs(configObject)

	utils.LogInfo("config", "version is "+utils.GreyString(Version))

	if CommitHash != "" {
		utils.LogInfo("config", "head is "+utils.GreyString(CommitHash))
	}

	LoadEnvFile(configObject)
	LoadYmlFile(configObject)
}

func LoadArgs(configObject *ServerConfig) {
	var (
		port = flag.Int("port", 0, "Port number to use")
		p    = flag.Int("p", 0, "Port number (shorthand)")

		env = flag.String("env", "local", "Environment to use")
		e   = flag.String("e", "local", "Environment to use (shorthand)")

		version = flag.Bool("version", false, "Print version")
		v       = flag.Bool("v", false, "Print version (shorthand)")

		help = flag.Bool("help", false, "Print help")
		h    = flag.Bool("h", false, "Print help (shorthand)")

		commitHash = flag.Bool("commit-hash", false, "Commit hash")
		c          = flag.Bool("c", false, "Commit hash (shorthand)")
	)

	flag.Parse()

	if (help != nil && *help) || (h != nil && *h) {
		flag.Usage()
		os.Exit(0)
	}

	if (version != nil && *version) || (v != nil && *v) {
		fmt.Println(Version)
		os.Exit(0)
	}

	if (commitHash != nil && *commitHash) || (c != nil && *c) {
		fmt.Println(CommitHash)
		os.Exit(0)
	}

	if e != nil && *e != "" {
		env = e
	}

	if *env != "" {
		configObject.EnvFilePath = "deployment-files/" + *env + "/.env"
	}

	if p != nil && *p != 0 {
		port = p
	}

	if *port != 0 {
		configObject.Port = strconv.Itoa(*port)
	}
}

func LoadEnvFile(configObject *ServerConfig) {
	if envLoadErr := godotenv.Load(configObject.EnvFilePath); envLoadErr != nil {
		if envLoadErr = godotenv.Load(); envLoadErr != nil {
			utils.LogFatal("server", "could not load .env file: "+envLoadErr.Error())
		}
	}

	if env := os.Getenv("ENV"); env != "" {
		configObject.Env = env
		utils.LogInfo("server", "ENV: "+utils.GreyString(env))
	} else {
		utils.LogWarn("server", "ENV not set in .env file: using "+utils.GreyString(configObject.Env))
	}

	if configObject.Port == "" {
		if port := os.Getenv("PORT"); port != "" {
			configObject.Port = port
			utils.LogInfo("server", "PORT: "+utils.GreyString(port))
		} else {
			utils.LogWarn("server", "PORT not set in .env file: using "+utils.GreyString(configObject.Port))
		}
	}

	if dbUrl := os.Getenv("PRIMARY_DB_URL"); dbUrl != "" {
		configObject.PrimaryDbUrl = dbUrl
	} else {
		utils.LogFatal("server", "DB_URL not set in .env file")
	}

	if secondaryDbUrls := os.Getenv("REPLICA_DB_URLS"); secondaryDbUrls != "" {
		var tmp []string

		for _, item := range strings.Split(secondaryDbUrls, ",") {
			tmp = append(tmp, item)
		}

		configObject.ReplicaDbUrls = tmp
	} else {
		utils.LogWarn("server", "SECONDARY_DB_URLS not set in .env file")
	}

	if cacheUrl := os.Getenv("CACHE_URL"); cacheUrl != "" {
		configObject.CacheUrl = cacheUrl
	} else {
		utils.LogFatal("server", "CACHE_URL not set in .env file")
	}

	if cachePassword := os.Getenv("CACHE_PASSWORD"); cachePassword != "" {
		configObject.CachePassword = cachePassword
	} else {
		utils.LogWarn("server", "CACHE_PASSWORD not set in .env file")
	}

	if cacheDatabase := os.Getenv("CACHE_DATABASE"); cacheDatabase != "" {
		if value, err := strconv.Atoi(cacheDatabase); err == nil {
			configObject.CacheDatabase = value
		} else {
			utils.LogWarn("server", "could not parse CACHE_DATABASE from .env file")
		}
	} else {
		utils.LogWarn("server", "CACHE_DATABASE not set in .env file")
	}

	if cacheSystemPrefix := os.Getenv("CACHE_SYSTEM_PREFIX"); cacheSystemPrefix != "" {
		configObject.CacheSystemPrefix = cacheSystemPrefix
	} else {
		utils.LogWarn("server", "CACHE_SYSTEM_PREFIX not set in .env file")
	}

	if allowedOrigins := os.Getenv("ALLOWED_ORIGINS"); allowedOrigins != "" {
		var split = strings.Split(allowedOrigins, ";")

		for i, origin := range split {
			split[i] = strings.TrimSpace(origin)
		}

		if len(split) == 0 {
			utils.LogFatal("server", "ALLOWED_ORIGINS not set in .env file")
		}

		ServerConfigObject.AllowedOrigins = split
	} else {
		utils.LogFatal("server", "ALLOWED_ORIGINS not set in .env file")
	}

	if selfBaseUrl := os.Getenv("SELF_BASE_URL"); selfBaseUrl != "" {
		ServerConfigObject.SelfBaseUrl = selfBaseUrl
	} else {
		utils.LogFatal("server", "SELF_BASE_URL not found in .env")
	}
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		configObject.JwtSecret = jwtSecret
	} else {
		utils.LogFatal("server", "JWT_SECRET not set in .env file")
	}

	if secretKeybase := os.Getenv("SECRET_KEY_BASE"); secretKeybase != "" {
		configObject.SecretKeyBase = secretKeybase
	} else {
		utils.LogFatal("server", "SECRET_KEY_BASE not set in .env file")
	}
}

func LoadYmlFile(configObject *ServerConfig) {
	currentDir, err := os.Getwd()
	if err != nil {
		utils.LogFatal("server", "could not get current directory")
	}

	ymlPath := currentDir + "/" + configObject.YmlFilePath
	ymlFile, err := os.Open(ymlPath)
	if err != nil {
		utils.LogFatal("server", "could not open yml file")
	}
	defer ymlFile.Close()

	var osmSettings OSMSettings
	decoder := yaml.NewDecoder(ymlFile)
	if err := decoder.Decode(&osmSettings); err != nil {
		utils.LogFatal("server", "could not parse yml file: "+err.Error())
	}

	// Assign the parsed settings to the ServerConfigObject
	configObject.OSMSettings = &osmSettings
}

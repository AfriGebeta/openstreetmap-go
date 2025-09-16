package client

import (
	"context"
	"github.com/redis/go-redis/v9"
	"openstreetmap-go/src/config"
	"openstreetmap-go/src/utils"
	"strconv"
)

type CacheConfig struct {
	Database int
	Url      string
	Password string
}

var RedisClient *CacheClient

type CacheClient struct {
	*redis.Client
	Config CacheConfig
}

func Load(serverConfig *config.ServerConfig, client **CacheClient) {
	if client == nil {
		utils.LogFatal("cache-connection-loader", "cache client not provided")
		return
	}

	*client = NewCacheClient(CacheConfig{
		Database: serverConfig.CacheDatabase,
		Url:      serverConfig.CacheUrl,
		Password: serverConfig.CachePassword,
	})
}

func NewCacheClient(config CacheConfig) *CacheClient {
	var cacheLocation = utils.GreyString(strconv.Itoa(config.Database)) + " at " + utils.GreyString(config.Url)

	utils.LogInfo(
		"cache-connection-loader",
		"connecting to cache "+
			utils.GreyString(strconv.Itoa(config.Database))+" at "+utils.GreyString(config.Url),
	)

	var client = redis.NewClient(&redis.Options{
		Addr:     config.Url,
		Password: config.Password,
		DB:       config.Database,
	})

	_, err := client.Ping(context.Background()).Result()

	if err != nil {
		utils.LogFatal("cache", "could not connect to Redis: "+err.Error())
	}

	utils.LogSuccess("cache-connection-loader", "connected to cache at "+cacheLocation)

	return &CacheClient{Client: client, Config: config}
}

func (cc CacheClient) CloseCacheConnection() {
	var err = cc.Close()

	if err != nil {
		utils.LogFatal("cache-connection-closer", "could not close cache connection: "+err.Error())
	}

	utils.LogSuccess(
		"cache-connection-closer",
		"closed cache connection "+
			utils.GreyString(strconv.FormatInt(int64(cc.Config.Database), 10))+" at "+utils.GreyString(cc.Config.Url),
	)
}

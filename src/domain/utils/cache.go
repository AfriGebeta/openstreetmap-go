package utils

import (
	"openstreetmap-go/src/config"
	"openstreetmap-go/src/domain/values"
)

func GetCachePrefix(prefix values.CachePrefix) string {
	return config.ServerConfigObject.CacheSystemPrefix + prefix.SurroundWithColons()
}

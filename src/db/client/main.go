package client

import (
	"fmt"
	"math/rand"
	"net/url"
	"openstreetmap-go/src/config"
	"openstreetmap-go/src/utils"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseConfig struct {
	PrimaryUrl  string
	ReplicaUrls []string
}

var DatabaseClients DatabaseConnections

type DatabaseConnections []*gorm.DB

func Load(configObject *config.ServerConfig, clients *DatabaseConnections) {
	if clients == nil {
		utils.LogFatal("database-connection-loader", "database clients not provided")
		return
	}

	*clients = NewDatabaseClients(DatabaseConfig{
		PrimaryUrl:  configObject.PrimaryDbUrl,
		ReplicaUrls: configObject.ReplicaDbUrls,
	})
}

func NewDatabaseClients(config DatabaseConfig) DatabaseConnections {
	var dbUrls = []string{config.PrimaryUrl}
	dbUrls = append(dbUrls, config.ReplicaUrls...)

	var (
		urlsString        = "urls"
		connectionsString = "connections"
	)

	if len(dbUrls) == 1 {
		urlsString = "url"
		connectionsString = "connection"
	}

	utils.LogInfo("database-connection-loader", fmt.Sprintf("found %s database %s", utils.YellowString(fmt.Sprint(len(dbUrls))), urlsString))
	utils.LogInfo("database-connection-loader", fmt.Sprintf("creating database %s...", connectionsString))

	var connections DatabaseConnections

	for _, connString := range dbUrls {
		db, err := gorm.Open(postgres.Open(connString), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			utils.LogFatal("database-connection-loader", "database connection error:"+err.Error())
			return nil
		}

		sqlDB, err := db.DB()
		if err != nil {
			utils.LogFatal("database-connection-loader", "failed to get underlying sql.DB: "+err.Error())
		}

		sqlDB.SetMaxOpenConns(19)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(30 * time.Minute)

		if err := sqlDB.Ping(); err != nil {
			utils.LogFatal("database-connection-loader", "database ping error: "+err.Error())
		}

		var dbInfo string
		if info, parseErr := url.Parse(connString); parseErr == nil {
			dbInfo = " " + utils.GreyString(strings.TrimPrefix(info.Path, "/")) + " at " + utils.GreyString(info.Host)
		}

		utils.LogSuccess("database-connection-loader", fmt.Sprintf("successfully connected to database%s", dbInfo))
		connections = append(connections, db)
	}

	return connections
}

func (dc DatabaseConnections) GetRandomConnection() *gorm.DB {
	return dc[rand.Intn(len(dc))]
}

func (dc DatabaseConnections) GetMasterConnection() *gorm.DB {
	return dc[0]
}

func (dc DatabaseConnections) CloseDbConnections() {
	var connectionsString = "connections"
	if len(dc) == 1 {
		connectionsString = "connection"
	}

	utils.LogInfo("database-connection-cleaner", "closing database "+connectionsString+"...")

	for _, db := range dc {
		sqlDB, err := db.DB()
		if err != nil {
			utils.LogError("database-connection-cleaner", "could not get sql.DB from gorm: "+err.Error())
			continue
		}

		if err := sqlDB.Close(); err != nil {
			utils.LogError("database-connection-cleaner", "could not close db connection: "+err.Error())
		}
	}

	utils.LogInfo("database-connection-cleaner", "finished closing database "+connectionsString)
}

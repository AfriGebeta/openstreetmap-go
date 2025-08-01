package db

import (
	"database/sql"
	"log"
	"math/rand"
	"openstreetmap-go/src/config"
)

var DbConnections []*sql.DB

func InitDb() {
	pgUrl := "postgres://" +
		config.ServerConfiguration.SmtpUserName + ":" +
		config.ServerConfiguration.SmtpPassword + "@" +
		config.ServerConfiguration.SmtpPassword + ":" +
		config.ServerConfiguration.SmtpPassword + "/" + config.ServerConfiguration.SmtpPassword
	db, err := sql.Open("postgres", pgUrl)
	if err != nil {
		log.Fatal(err)
	}
	DbConnections = append(DbConnections, db)
}

func ReturnRandomDbConnection() *sql.DB {
	return DbConnections[rand.Intn(len(DbConnections))]
}

func ReturnMasterDbConnection() *sql.DB {
	return DbConnections[0]
}

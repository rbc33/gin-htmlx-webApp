package main

import (
	admin_app "gocms/admin-app"
	"gocms/common"
	"gocms/database"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
)

func main() {

	// sets zerolog as the main logger
	// in this APP
	common.SetupLogger()

	database, err := database.MakeSqlConnection()
	if err != nil {
		log.Error().Msgf("could not create database connection: %v", err)
		return
	}

	err = admin_app.Run(database)
	if err != nil {
		log.Fatal().Msgf("could not run app: %v", err)
	}
}

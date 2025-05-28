package main

import (
	_ "github.com/go-sql-driver/mysql"
	admin_app "github.com/rbc33/admin-app"
	"github.com/rbc33/common"
	"github.com/rbc33/database"
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

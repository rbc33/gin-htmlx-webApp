package main

import (
	"gocms/app"
	"gocms/common"
	"gocms/database"

	_ "github.com/go-sql-driver/mysql"

	"github.com/rs/zerolog/log"
)

func main() {

	db_connection, err := database.MakeSqlConnection()
	if err != nil {
		log.Error().Msgf("could not create database connection: %v", err)
		return
	}
	common.SetupLogger()
	err = app.Run(&db_connection)
	if err != nil {
		log.Error().Msgf("%s", err)
	}

}

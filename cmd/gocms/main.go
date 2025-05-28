package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/rbc33/app"
	"github.com/rbc33/common"
	"github.com/rbc33/database"

	"github.com/rs/zerolog/log"
)

func main() {

	db_connection, err := database.MakeSqlConnection()
	if err != nil {
		log.Error().Msgf("could not create database connection: %v", err)
		return
	}
	common.SetupLogger()
	app.Run(db_connection)
}

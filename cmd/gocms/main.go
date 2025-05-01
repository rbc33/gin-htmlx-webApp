package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/rbc33/app"
	"github.com/rbc33/database"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func setupLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().Msg("Logger created")
}

func main() {

	db_connection, err := database.MakeSqlConnection()
	if err != nil {
		log.Error().Msgf("could not create database connection: %v", err)
	}
	setupLogger()
	app.Run(db_connection)
}

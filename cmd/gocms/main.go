package main

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/rbc33/gocms/app"
	"github.com/rbc33/gocms/common"
	"github.com/rbc33/gocms/database"

	"github.com/rs/zerolog/log"
)

func main() {

	common.SetupLogger()

	err := godotenv.Load()
	if err != nil {
		log.Error().Msgf("%s", err)
		os.Exit(-1)
	}

	db_connection, err := database.MakeSqlConnection()
	if err != nil {
		log.Error().Msgf("could not create database connection: %v", err)
		return
	}
	r := app.SetupRoutes(&db_connection)
	err = r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))

	if err != nil {
		log.Error().Msgf("could not run app: %v", err)
		os.Exit(-1)

	}
}

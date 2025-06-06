package main

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	admin_app "github.com/rbc33/admin-app"
	"github.com/rbc33/common"
	"github.com/rbc33/database"
	"github.com/rs/zerolog/log"
)

func main() {

	// sets zerolog as the main logger
	// in this APP
	common.SetupLogger()

	err := godotenv.Load()
	if err != nil {
		log.Error().Msgf("%s", err)
		os.Exit(-1)
	}

	database, err := database.MakeSqlConnection()
	if err != nil {
		log.Error().Msgf("could not create database connection: %v", err)
		os.Exit(-1)

	}

	r := admin_app.SetupRoutes(&database)
	err = r.Run(fmt.Sprintf(":%s", os.Getenv("PORT_ADMIN")))
	if err != nil {
		log.Error().Msgf("could not run app: %v", err)
		os.Exit(-1)

	}
}

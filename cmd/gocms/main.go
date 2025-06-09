package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rbc33/gocms/app"
	"github.com/rbc33/gocms/common"
	"github.com/rbc33/gocms/database"

	"github.com/rs/zerolog/log"
)

func main() {
	config_toml := flag.String("config", "", "path to the config file")
	flag.Parse()

	common.SetupLogger()

	if (*config_toml) != "" {
		log.Info().Msgf("Reading config file %s", *config_toml)
		settings, err := common.ReadConfigToml(*config_toml)
		if err != nil {
			log.Error().Msgf("Could not read config file: %s", err)
			os.Exit(-1)
		}
		common.GetSettings(settings)

	}

	// err := godotenv.Load()
	// if err != nil {
	// 	log.Error().Msgf("%s", err)
	// 	os.Exit(-1)
	// }

	db_connection, err := database.MakeSqlConnection(common.Settings)
	if err != nil {
		log.Error().Msgf("could not create database connection: %v", err)
		return
	}
	Port := os.Getenv("PORT")
	if Port == "" {
		Port = common.Settings.WebserverPort
	}
	r := app.SetupRoutes(common.Settings, &db_connection)
	err = r.Run(fmt.Sprintf(":%s", Port))

	if err != nil {
		log.Error().Msgf("could not run app: %v", err)
		os.Exit(-1)

	}
}

package main

import (
	"gocms/app"
	"gocms/common"
	"gocms/database"

	_ "github.com/go-sql-driver/mysql"

	"github.com/rs/zerolog/log"
)

func main() {

	common.SetupLogger()

	db_connection, err := database.MakeSqlConnection()
	if err != nil {
		log.Warn().Msgf("DB connection failed (expected on first deploy): %v", err)
		// TEMPORAL: Continuar sin DB para que la app inicie
		log.Info().Msg("Starting app without DB connection...")
		// Crear conexión dummy o usar una struct vacía
		db_connection = database.Database{}
	}

	err = app.Run(&db_connection)
	if err != nil {
		log.Error().Msgf("%s", err)
	}

}

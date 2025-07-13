package database

import (
	"database/sql"
	"time"

	"github.com/rbc33/gocms/common"
	"github.com/rs/zerolog/log"
)

func MakeSqlConnection(appSettings common.AppSettings) (database SqlDatabase, err error) {
	/// Checking the DB connection
	// err := godotenv.Load()
	// if err != nil {
	// 	return SqlDatabase{}, err
	// }
	connection_str := appSettings.DatabaseUri
	db, err := sql.Open("mysql", connection_str)
	if err != nil {
		return SqlDatabase{}, err
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Second * 5)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	log.Info().Msg(connection_str)

	return SqlDatabase{
		MY_SQL_URL: connection_str,
		Connection: db,
	}, nil
}

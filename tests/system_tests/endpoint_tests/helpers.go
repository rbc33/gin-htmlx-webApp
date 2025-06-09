package index_test

import (
	"context"
	_ "database/sql"
	"embed"
	"fmt"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/rbc33/gocms/common"
	"github.com/rbc33/gocms/database"
)

//go:generate ../../../migrations ./migrations

//go:embed migrations/*.sql
var embedMigrations embed.FS

var current_port int = 0

func runDatabaseServer(app_settings common.AppSettings) {
	// Parse the database URI to get the database name
	dsn, err := mysql.ParseDSN(app_settings.DatabaseUri)
	if err != nil {
		panic(fmt.Errorf("invalid database URI: %v", err))
	}

	pro := createTestDatabase(dsn.DBName)
	engine := sqle.NewDefault(pro)
	engine.Analyzer.Catalog.MySQLDb.AddRootAccount()

	session := memory.NewSession(sql.NewBaseSession(), pro)
	ctx := sql.NewContext(context.Background(), sql.WithSession(session))
	ctx.SetCurrentDatabase(dsn.DBName)

	config := server.Config{
		Protocol: "tcp",
		Address:  fmt.Sprintf("%s:%s", dsn.Addr, dsn.Params["port"]),
	}
	s, err := server.NewServer(
		config,
		engine,
		nil, // use default ContextFactory
		memory.NewSessionBuilder(pro),
		nil, // no ServerEventListener
	)
	if err != nil {
		panic(err)
	}
	if err = s.Start(); err != nil {
		panic(err)
	}
}

func createTestDatabase(name string) *memory.DbProvider {
	db := memory.NewDatabase(name)
	db.BaseDatabase.EnablePrimaryKeyIndexes()

	pro := memory.NewDBProvider(db)
	return pro
}

func waitForDb(app_settings common.AppSettings) (database.SqlDatabase, error) {

	for range 400 {
		database, err := database.MakeSqlConnection(app_settings)

		if err == nil {
			return database, nil
		}

		time.Sleep(25 * time.Millisecond)
	}

	return database.SqlDatabase{}, fmt.Errorf("database did not start")
}

func getAppSettings() common.AppSettings {
	if os.Getenv("GITHUB_ACTIONS") != "" {
		return common.AppSettings{
			DatabaseUri:   "root:root@tcp(mysql:3306)/gocms",
			WebserverPort: "8080",
		}
	}
	return common.AppSettings{
		DatabaseUri:   "root:secret@tcp(192.168.0.100:33060)/gocms",
		WebserverPort: "8080",
	}
}

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
// @title        GoCMS Admin API
// @version      1.0.0
// @description  This is the admin API for the GoCMS app.
// @schemes   http
// @host      localhost:8081
// @BasePath  /
// @contact.name   Ricardo
// @contact.email  ricardobenthem@gmail.com
// @license.name  MIT
// @consumes  application/json
// @consumes  multipart/form-data
// @produces  application/json
package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	// lua "github.com/yuin/gopher-lua"

	// "github.com/joho/godotenv"
	admin_app "github.com/rbc33/gocms/admin-app"
	"github.com/rbc33/gocms/common"
	"github.com/rbc33/gocms/database"
	"github.com/rbc33/gocms/plugins"

	// "github.com/rbc33/gocms/plugins"
	"github.com/joho/godotenv"
	_ "github.com/rbc33/gocms/docs"
	"github.com/rs/zerolog/log"
)

func main() {

	config_toml := flag.String("config", "", "path to the config file")
	flag.Parse()

	common.SetupLogger("error-admin.log")
	err := godotenv.Load()
	if err != nil {
		log.Error().Msgf("Error loading .env file")
	}
	fmt.Println("TOKEN_HOUR_LIFESPAN:", os.Getenv("TOKEN_HOUR_LIFESPAN"))

	if (*config_toml) != "" {
		log.Info().Msgf("Reading config file %s", *config_toml)
		settings, err := common.ReadConfigToml(*config_toml)
		if err != nil {
			log.Error().Msgf("Could not read config file %s", err)
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
	Port := (os.Getenv("PORT_ADMIN"))
	if Port == "" {
		Port = common.Settings.WebserverPortAdmin
	}
	shortcode_handlers, err := admin_app.LoadShortcodesHandlers(common.Settings.Shortcodes)
	if err != nil {
		log.Error().Msgf("%s", err)
		os.Exit(-1)
	}

	// TODO : we probably want to refactor loadShortcodeHandler
	// TODO : into loadPluginHandlers instead

	post_hook := &plugins.PostHook{}
	image_plugin := plugins.Plugin{
		ScriptName: "img",
		Id:         "img-plugin",
	}
	post_hook.Register(image_plugin)
	// img, _ := shortcode_handlers["img"]
	hooks_map := map[string]plugins.Hook{
		"add_post": post_hook,
	}

	r := admin_app.SetupRoutes(common.Settings, shortcode_handlers, &db_connection, hooks_map)
	// r := admin_app.SetupRoutes(common.Settings, shortcode_handlers, &db_connection)// Esta línea añade la ruta para la UI de Swagger.
	// // La URL será: http://localhost:8081/swagger/index.html
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	err = r.Run(fmt.Sprintf(":%s", Port))
	if err != nil {
		log.Error().Msgf("could not run app: %v", err)
		os.Exit(-1)

	}
}

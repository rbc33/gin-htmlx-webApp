package common

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

type CardSchema struct {
	Name string `toml:"schema_name"`
}
type Shortcode struct {
	// name for the shortcode {{name:...:...}}
	Name string `toml:"name"`
	// the lua plugin path
	Plugin string `toml:"plugin"`
}

type AppSettings struct {
	DatabaseUri        string       `toml:"MY_SQL_URL"`
	MediaDir           string       `toml:"MEDIA_DIR"`
	WebserverPort      string       `toml:"PORT"`
	WebserverPortAdmin string       `toml:"PORT_ADMIN"`
	CardSchema         []CardSchema `toml:"card_schema"`
	Shortcodes         []Shortcode  `toml:"shortcodes"`
}

type Navbar struct {
	Links []Link `toml:"links"`
}

var Settings AppSettings

func GetSettings(settings AppSettings) {
	Settings = settings
}
func ReadConfigToml(filepath string) (AppSettings, error) {
	var config AppSettings
	// Initialize CardSchema as empty slice before decoding
	config.CardSchema = []CardSchema{}

	meta, err := toml.DecodeFile(filepath, &config)
	if err != nil {
		return AppSettings{}, fmt.Errorf("failed to decode config file: %w", err)
	}

	// Improve error handling for undecoded keys
	if undecoded := meta.Undecoded(); len(undecoded) > 0 {
		var keys []string
		for _, key := range undecoded {
			keys = append(keys, strings.Join(key, "."))
		}
		return config, fmt.Errorf("undecoded keys found: %s", strings.Join(keys, ", "))
	}

	// Validate required fields
	if config.DatabaseUri == "" {
		return config, fmt.Errorf("MY_SQL_URL is required")
	}
	if config.WebserverPort == "" {
		return config, fmt.Errorf("PORT is required")
	}

	return config, nil
}

func IsGithubActions() bool {
	return os.Getenv("GITHUB_ACTIONS") != ""
}

func GetTestDatabaseUri() string {
	if IsGithubActions() {
		return "root:root@tcp(mysql:3306)/gocms"
	}
	// Local Docker MySQL instance
	return "root:secret@tcp(192.168.0.100:33060)/gocms"
}

func GetTestServerAddress() string {
	if IsGithubActions() {
		return "mysql:3306"
	}
	// For local testing, bind to localhost but connect to Docker
	return "localhost:8080"
}

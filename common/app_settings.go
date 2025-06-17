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
	WebserverPort      string       `toml:"PORT"`
	WebserverPortAdmin string       `toml:"PORT_ADMIN"`
	CardSchema         []CardSchema `toml:"card_schema"`
	Shortcodes         []Shortcode  `toml:"shortcodes"`
	ImageDirectory     string       `toml:"image_dir"`
	CacheEnabled       bool         `toml:"cache_enabled"`
	AppNavbar          Navbar       `toml:"navbar"`
	RecaptchaSiteKey   string       `toml:"recaptcha_sitekey, omitempty"`
	RecaptchaSecret    string       `toml:"recaptcha_secret, omitempty"`
	AppDomain          string       `toml:"app_domain, omitempty"`
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
	if db_uri := GetTestDatabaseUri(); db_uri != "" {
		config.DatabaseUri = db_uri
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
func IsDocker() string {
	return os.Getenv("DOCKER_DB_URI")
}

func GetTestDatabaseUri() string {
	if MY_SQL_URL := os.Getenv("MY_SQL_URL"); MY_SQL_URL != "" {
		return MY_SQL_URL
	} else if IsGithubActions() {
		return "root:root@tcp(mysql:3306)/gocms"
	} else if docker_uri := IsDocker(); docker_uri != "" {
		return docker_uri
	}

	return "root:secret@tcp(192.168.0.100:33060)/gocms"
}

func GetTestServerAddress() string {
	if IsGithubActions() {
		return "mysql:3306"
	}
	// For local testing, bind to localhost but connect to Docker
	return "localhost:8080"
}

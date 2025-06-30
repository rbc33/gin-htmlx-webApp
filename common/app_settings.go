package common

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

type Shortcode struct {
	// name for the shortcode {{name:...:...}}
	Name string `toml:"name"`
	// the lua plugin path
	Plugin string `toml:"plugin"`
}

type AppSettings struct {
	DatabaseUri        string             `toml:"MY_SQL_URL"`
	WebserverPort      string             `toml:"PORT"`
	WebserverPortAdmin string             `toml:"PORT_ADMIN"`
	CardSchema         []CardSchema       `toml:"card_schema"`
	Shortcodes         []Shortcode        `toml:"shortcodes"`
	ImageDirectory     string             `toml:"image_dir"`
	CacheEnabled       bool               `toml:"cache_enabled"`
	AppNavbar          Navbar             `toml:"navbar"`
	RecaptchaSiteKey   string             `toml:"recaptcha_sitekey, omitempty"`
	RecaptchaSecret    string             `toml:"recaptcha_secret, omitempty"`
	AppDomain          string             `toml:"app_domain, omitempty"`
	Galleries          map[string]Gallery `toml:"gallery"`
	StickyPosts        []int              `toml:"sticky_posts"`
}

type Navbar struct {
	Links     []Link            `toml:"links"`
	Dropdowns map[string][]Link `toml:"dropdowns"`
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
	// Override with environment variables if they exist, which is common in containers.
	if dbURI := GetDatabaseURIFromEnv(); dbURI != "" {
		config.DatabaseUri = dbURI
	}
	// Validate required fields
	if config.DatabaseUri == "" {
		return config, fmt.Errorf("DatabaseUri is required in config or DOCKER_DB_URI/MY_SQL_URL env var must be set")
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

func IsKubernetes() bool {
	return os.Getenv("KUBERNETES_SERVICE_HOST") != ""
}

// GetDatabaseURIFromEnv checks environment variables for a database URI in a specific order of precedence.
// It returns an empty string if no relevant environment variable is found.
func GetDatabaseURIFromEnv() string {
	// DOCKER_DB_URI is the primary override for Kubernetes/Docker Compose
	if dockerURI := os.Getenv("DOCKER_DB_URI"); dockerURI != "" {
		return dockerURI
	}
	if mySQLURL := os.Getenv("MY_SQL_URL"); mySQLURL != "" && os.Getenv("ANDROID_HOME") == "" {
		return mySQLURL
	}
	if IsGithubActions() {
		return "root:root@tcp(mysql:3306)/gocms"
	}
	return "" // No environment override found
}

func GetTestServerAddress() string {
	if IsGithubActions() {
		return "mysql:3306"
	}
	// For local testing, bind to localhost but connect to Docker
	return "localhost:8080"
}

package common

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
)

type CardSchema struct {
	Name string `toml:"schema_name"`
}

type AppSettings struct {
	DatabaseUri        string       `toml:"MY_SQL_URL"`
	MediaDir           string       `toml:"MEDIA_DIR"`
	WebserverPort      string       `toml:"PORT"`
	WebserverPortAdmin string       `toml:"PORT_ADMIN"`
	CardSchema         []CardSchema `toml:"card_schema"`
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
	// Initialize CardSchema as empty slice instead of nil
	config.CardSchema = []CardSchema{}

	meta, err := toml.DecodeFile(filepath, &config)
	if err != nil {
		return AppSettings{}, err
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

package common

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog/log"
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
	meta, err := toml.DecodeFile(filepath, &config)
	if err != nil {
		return AppSettings{}, err
	}
	if undecoded_keys := meta.Undecoded(); len(undecoded_keys) > 0 {
		err := fmt.Errorf("cound not decode keys: ")
		for _, key := range undecoded_keys {
			log.Error().Msgf("%v, %v", err, strings.Join(key, ", "))
		}
	}
	return config, nil
}

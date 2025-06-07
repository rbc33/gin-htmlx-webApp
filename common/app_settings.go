package common

import (
	"github.com/BurntSushi/toml"
)

type CardSchema struct {
	Name string `toml:"schema_name"`
	Path string `toml:"schema_path"`
}

type AppSettings struct {
	DatabaseUri   string       `toml:"database_uri"`
	MediaDir      string       `toms:"MEDIA_DIR"`
	WebserverPost string       `toml:"webserver_port"`
	CardSchema    []CardSchema `toml:"card_schema"`
}

type Navbar struct {
	Links []Link `toml:"links"`
}

func ReadConfigToml(filepath string) (AppSettings, error) {
	var config AppSettings
	_, err := toml.DecodeFile(filepath, &config)
	if err != nil {
		return AppSettings{}, err
	}

	return config, nil
}

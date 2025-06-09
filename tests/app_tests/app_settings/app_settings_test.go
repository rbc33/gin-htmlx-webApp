package app_settings_tests

import (
	"errors"
	"os"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/rbc33/gocms/common"
	"github.com/stretchr/testify/assert"
)

// Writes the contents into a temporary
// toml file
func writeToml(contents []byte) (s string, err error) {
	file, err := os.CreateTemp(os.TempDir(), "*.toml")
	if err != nil {
		return "", err
	}
	defer func() {
		err = errors.Join(file.Close(), err)
	}()

	_, err = file.Write(contents)
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}

func TestCorrectToml(t *testing.T) {
	expected := common.AppSettings{
		DatabaseUri:   "test_database_address",
		WebserverPort: "99999",
	}
	bytes, err := toml.Marshal(expected)
	assert.Nil(t, err)

	filepath, err := writeToml(bytes)
	assert.Nil(t, err)

	actual, err := common.ReadConfigToml(filepath)
	assert.Nil(t, err)
	assert.Equal(t, actual, expected)
}

func TestMissingDatabaseAddress(t *testing.T) {

	missing_database_address := struct {
		WebserverPort string `toml:"webserver_port"`
	}{
		WebserverPort: "99999",
	}

	bytes, err := toml.Marshal(missing_database_address)
	assert.Nil(t, err)

	filepath, err := writeToml(bytes)
	assert.Nil(t, err)

	_, err = common.ReadConfigToml(filepath)
	assert.NotNil(t, err)
}

func TestWrongwebserverPortValueType(t *testing.T) {
	missing_database_address := struct {
		DatabaseUri   string `toml:"database_usri"`
		WebserverPort int    `toml:"webserver_port"`
	}{
		DatabaseUri:   "test_database_uri",
		WebserverPort: 99999,
	}

	bytes, err := toml.Marshal(missing_database_address)
	assert.Nil(t, err)

	filepath, err := writeToml(bytes)
	assert.Nil(t, err)

	_, err = common.ReadConfigToml(filepath)
	assert.NotNil(t, err)
}

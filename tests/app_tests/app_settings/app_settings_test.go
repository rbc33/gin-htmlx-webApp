package app_settings_tests

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/rbc33/gocms/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
func gitAction() string {
	if os.Getenv("GITHUB_ACTIONS") != "" {
		return "root:root@tcp(mysql:3306)/gocms"
	} else {
		return "root:secret@tcp(192.168.0.100:33060)/gocms"
	}

}
func TestCorrectToml(t *testing.T) {
	db := gitAction()
	content := fmt.Sprintf(`
MY_SQL_URL = "%v"
PORT = "99999"
`, db)
	tmpfile, err := os.CreateTemp("", "config.*.toml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.WriteString(content)
	require.NoError(t, err)

	settings, err := common.ReadConfigToml(tmpfile.Name())
	require.NoError(t, err)

	expected := common.AppSettings{
		DatabaseUri:   db,
		WebserverPort: "99999",
		CardSchema:    []common.CardSchema{}, // Initialize as empty slice
	}
	assert.Equal(t, expected, settings)
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

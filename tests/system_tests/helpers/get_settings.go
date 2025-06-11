package test

import (
	"os"

	"github.com/rbc33/gocms/common"
)

func GetAppSettings() common.AppSettings {
	if os.Getenv("GITHUB_ACTIONS") != "" {
		return common.AppSettings{
			DatabaseUri:   "root:root@tcp(mysql:3306)/gocms",
			WebserverPort: "8080",
		}
	}

	// Local Docker MySQL instance settings
	return common.AppSettings{
		DatabaseUri:   "root:secret@tcp(192.168.0.100:33060)/gocms",
		WebserverPort: "8080",
	}
}

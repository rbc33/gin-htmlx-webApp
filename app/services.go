package app

import (
	"github.com/gin-gonic/gin"
	"github.com/rbc33/gocms/common"
	"github.com/rbc33/gocms/database"
	"github.com/rbc33/gocms/views"
)

func servicesHandler(c *gin.Context, db database.Database) ([]byte, error) {
	return renderHtml(c, views.MakeServicesPage(common.Settings.AppNavbar.Links, common.Settings.AppNavbar.Dropdowns))
}

package admin_app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rbc33/gocms/common"
	"github.com/rbc33/gocms/database"
	"github.com/rs/zerolog/log"
)

// @Summary      Add a new permalink
// @Description  Adds a new permalink to the database.
// @Tags         permalink
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        permalink path string true "Permalink"
// @Param        post_id path int true "Post ID"
// @Success      200 {object} PermalinkIdResponse
// @Failure      400 {object} common.ErrorResponse "Invalid request body or missing data"
// @Router       /permalinks/{permalink}/{post_id} [post]
func postPermalinkHandler(database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		// add_permalink_request := struct {
		// 	Permalink string `uri:"permalink" binding:"required"`
		// 	PostId    int    `uri:"post_id" binding:"required"`
		// }{}
		var add_permalink_request AddPermalinkRequest

		err := c.ShouldBindUri(&add_permalink_request)
		if err != nil {
			log.Error().Msgf("invalid request for adding permalink: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("could not add permalink", err))
			return
		}

		permalinkDb := common.Permalink{
			Path:   add_permalink_request.Permalink,
			PostId: add_permalink_request.PostId,
		}

		id, err := database.AddPermalink(permalinkDb)

		if err != nil {
			log.Error().Msgf("failed to add post: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("could not add post", err))
			return
		}

		c.JSON(http.StatusOK, PermalinkIdResponse{
			id,
		})
	}
}

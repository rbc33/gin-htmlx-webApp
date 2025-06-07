package admin_app

import (
	"encoding/json"
	"net/http"

	"github.com/fossoreslp/go-uuid-v4"
	"github.com/gin-gonic/gin"
	"github.com/rbc33/gocms/common"
	"github.com/rbc33/gocms/database"
	"github.com/rs/zerolog/log"
)

type AddCardRequest struct {
	ImageLocation string `json:"image_location"`
	JsonData      string `json:"json_data"`
	SchemaName    string `json:"json_schema"`
}

type ChangeCardRequest struct {
	Id            string `json:"uuid"`
	ImageLocation string `json:"image_location"`
	JsonData      string `json:"json_data"`
	SchemaName    string `json:"json_schema"`
}

type DeleteCardRequest struct {
	Id string `json:"uuid"`
}

func getCardHandler(database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		// localhost:8080/Card/{id}
		var Card_binding common.CardIdBinding
		if err := c.ShouldBindUri(&Card_binding); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "could not get Card id",
				"msg":   err.Error(),
			})
			return
		}

		card, err := database.GetCard(Card_binding.Id)
		if err != nil {
			log.Warn().Msgf("could not get Card from DB: %v", err)
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Card id not found",
				"msg":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"uuid":           card.Uuid,
			"image_location": card.ImageLocation,
			"json_data":      card.JsonData,
			"json_schema":    card.SchemaName,
		})
	}
}

func postCardHandler(database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		var add_card_request AddCardRequest
		decoder := json.NewDecoder(c.Request.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&add_card_request)

		if err != nil {
			log.Warn().Msgf("could not get post from DB: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
				"msg":   err.Error(),
			})
			return
		}

		uuid, err := uuid.New()
		if err != nil {
			log.Error().Msgf("error generating UUID: %v", err)
		}

		err = database.AddCard(
			uuid.String(),
			add_card_request.ImageLocation,
			add_card_request.JsonData,
			add_card_request.SchemaName,
		)
		if err != nil {
			log.Error().Msgf("failed to add Card: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "could not add Card",
				"msg":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id": uuid.String(),
		})
	}
}

func putCardHandler(database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		var change_card_request ChangeCardRequest
		decoder := json.NewDecoder(c.Request.Body)
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&change_card_request)
		if err != nil {
			log.Warn().Msgf("could not get post from DB: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
				"msg":   err.Error(),
			})
			return
		}

		err = database.ChangeCard(
			change_card_request.Id,
			change_card_request.ImageLocation,
			change_card_request.JsonData,
			change_card_request.SchemaName,
		)
		if err != nil {
			log.Error().Msgf("failed to change card: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "could not change card",
				"msg":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id": change_card_request.Id,
		})
	}
}

func deleteCardHandler(database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		var delete_card_request DeleteCardRequest
		decoder := json.NewDecoder(c.Request.Body)
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&delete_card_request)
		if err != nil {
			log.Warn().Msgf("could not delete card: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
				"msg":   err.Error(),
			})
			return
		}

		err = database.DeleteCard(delete_card_request.Id)
		if err != nil {
			log.Error().Msgf("failed to delete card: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "could not delete card",
				"msg":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id": delete_card_request.Id,
		})
	}
}

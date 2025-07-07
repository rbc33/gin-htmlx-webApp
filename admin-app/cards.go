package admin_app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kaptinlin/jsonschema"
	"github.com/rbc33/gocms/common"
	"github.com/rbc33/gocms/database"
	"github.com/rs/zerolog/log"
)

// @Summary      Get a card list
// @Description  Retrieves a paginated list of cards by schema UUID.
// @Tags         cards
// @Produce      json
// @Param        schema path int true "schema UUID"
// @Success      200 {object} GetCardRequest
// @Failure      400 {object} common.ErrorResponse "Invalid post ID"
// @Failure      404 {object} common.ErrorResponse "Post not found"
// @Router       /cards/{schema} [get]
func getCardHandler(database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {

		var get_card_request GetCardRequest

		err := c.ShouldBindUri(&get_card_request)
		if err != nil {
			log.Error().Msgf("could not bind url params: %v", err)
			c.JSON(http.StatusBadRequest, common.MsgErrorRes("invalid card request, missing information"))
			return
		}

		if (get_card_request.Limit == 0) && (get_card_request.Page != 0) {
			log.Error().Msgf("card limit is 0 but pages is %d", get_card_request.Page)
			c.JSON(http.StatusBadRequest, common.MsgErrorRes("card limit is 0 but page is not"))
			return
		}

		if (get_card_request.Page == 0) && (get_card_request.Limit != 0) {
			log.Error().Msgf("card page is 0 but limit is %d", get_card_request.Limit)
			c.JSON(http.StatusBadRequest, common.MsgErrorRes("card page is 0 but limit is not"))
			return
		}

		limit := get_card_request.Limit
		page := get_card_request.Page
		if (get_card_request.Limit == 0) && (get_card_request.Page == 0) {
			limit = 10
			page = 0
		}

		cards, err := database.GetCards(get_card_request.Schema, int(limit), int(page))
		if err != nil {
			log.Error().Msgf("could not get cards: %v", err)
			c.JSON(http.StatusBadRequest, common.MsgErrorRes("invalid card schema uuid"))
			return
		}

		c.JSON(http.StatusOK, cards)
	}
}

// @Summary      Add a new post
// @Description  Adds a new post to the database.
// @Tags         cards
// @Accept       json
// @Produce      json
// @Param        post body AddCardRequest true "Post to add"
// @Success      200 {object} CardIdResponse
// @Failure      400 {object} common.ErrorResponse "Invalid request body or missing data"
// @Router       /cards [post]
func postCardHandler(database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		var add_card_request AddCardRequest
		if c.Request.Body == nil {
			c.JSON(http.StatusBadRequest, common.MsgErrorRes("no request body provided"))
			return
		}
		decoder := json.NewDecoder(c.Request.Body)
		err := decoder.Decode(&add_card_request)

		if err != nil {
			log.Warn().Msgf("invalid post request: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("invalid request body", err))
			return
		}

		// TODO : Sanity checks that everything inside the
		// TODO : request makes sense. I.e. content is json,
		// TODO : i.e json content matches the schema, etc.
		// err = checkRequiredData(add_card_request)
		// if err != nil {
		// 	log.Error().Msgf("failed to add post required data is missing: %v", err)
		// 	c.JSON(http.StatusBadRequest, common.ErrorRes("missing required data", err))
		// 	return
		// }

		// Check that the schema exists
		schema, err := database.GetCardSchema(add_card_request.Schema)
		if err != nil {
			log.Error().Msgf("card schema does not exist: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("card schema does not exist", err))
			return
		}
		// log.Info().Msgf("schema en req: %v", add_card_request.Schema)
		// log.Info().Msgf("schema en post: %v", schema)

		err = validateCardAgainstSchema(add_card_request.Content, schema.Schema)
		if err != nil {
			log.Error().Msgf("%v", err.Error())
			c.JSON(http.StatusBadRequest, common.ErrorRes("could not add card", err))
			return
		}

		id, err := database.AddCard(
			add_card_request.Image,
			add_card_request.Schema,
			add_card_request.Content,
		)
		if err != nil {
			log.Error().Msgf("failed to add card: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("could not add card", err))
			return
		}

		c.JSON(http.StatusOK, CardIdResponse{
			id,
		})
	}
}

func validateCardAgainstSchema(card_data string, json_schema string) error {

	// Parse the schema here
	schema_compiler := jsonschema.NewCompiler()
	schema, err := schema_compiler.Compile([]byte(json_schema))

	// log.Info().Msgf("schema en validate: %v", schema)

	if err != nil {
		return fmt.Errorf("failed to compile the json_schema from db: %v", err)
	}

	json_map := make(map[string]interface{})
	err = json.Unmarshal([]byte(card_data), &json_map)
	if err != nil {
		return fmt.Errorf("failed to parse card json : %v", err)
	}

	result := schema.Validate(json_map)
	if !result.IsValid() {
		details, _ := json.MarshalIndent(result.ToList(), "", "  ")
		return fmt.Errorf("failed to check vard data against schema: %v", string(details))
	}

	return nil
}

// @Summary      Update an existing card
// @Description  Updates an existing card with new data.
// @Tags         cards
// @Accept       json
// @Produce      json
// @Param        card body ChangeCardRequest true "Card data to update"
// @Success      200 {object} CardIdResponse
// @Failure      400 {object} common.ErrorResponse "Invalid request body or could not change card"
// @Router       /cards [put]
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

// @Summary      Delete a card
// @Description  Deletes a card by its ID.
// @Tags         cards
// @Produce      json
// @Param        id body DeletePostRequest true "Card ID to delete"
// @Success      200 {object} PostIdResponse
// @Failure      400 {object} common.ErrorResponse "Invalid ID provided"
// @Failure      404 {object} common.ErrorResponse "Card not found"
// @Router       /cards [delete]
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

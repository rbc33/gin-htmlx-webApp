package admin_app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kaptinlin/jsonschema"
	"github.com/rbc33/gocms/common"
	"github.com/rbc33/gocms/database"
	"github.com/rs/zerolog/log"
)

// @Summary      Add a new card schema
// @Description  Creates a new card schema by validating and storing the provided JSON schema and title.
// @Tags         card_schema
// @Accept       json
// @Produce      json
// @Param        schema body AddCardSchemaRequest true "Card schema to add"
// @Success      200 {object} CardIdResponse
// @Failure      400 {object} common.ErrorResponse "Invalid request body or schema"
// @Router       /card-schemas [post]
func postSchemaHandler(database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		var add_schema_request AddCardSchemaRequest
		if c.Request.Body == nil {
			c.JSON(http.StatusBadRequest, common.MsgErrorRes("no request body provided"))
			return
		}

		// Validate the content of the schema
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.MsgErrorRes("could not read request body"))
			return
		}

		if !json.Valid(body) {
			c.JSON(http.StatusBadRequest, common.MsgErrorRes("invalid json given in request body"))
			return
		}

		err = json.Unmarshal(body, &add_schema_request)
		if err != nil {
			error_msg := fmt.Errorf("could not unmarshall json request: %v", err)
			c.JSON(http.StatusBadRequest, common.MsgErrorRes(error_msg.Error()))
			return
		}

		if err = checkSchemaValues(add_schema_request); err != nil {
			c.JSON(http.StatusBadRequest, common.MsgErrorRes(err.Error()))
			return
		}

		id, err := database.AddCardSchema(
			add_schema_request.JsonSchema,
			add_schema_request.JsonTitle,
		)
		if err != nil {
			log.Error().Msgf("failed to add card schema: %v", err)
			c.JSON(http.StatusBadRequest, common.ErrorRes("could not add card schema", err))
			return
		}

		c.JSON(http.StatusOK, CardIdResponse{
			id,
		})
	}
}

// @Summary      Get a card schema by ID
// @Description  Retrieves a card schema by its UUID.
// @Tags         card_schema
// @Produce      json
// @Param        id path string true "Card schema UUID"
// @Success      200 {object} common.CardSchema
// @Failure      400 {object} common.ErrorResponse "Invalid schema ID"
// @Failure      404 {object} common.ErrorResponse "Schema not found"
// @Router       /card-schemas/{id} [get]
func getSchemaHandler(database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		// localhost:8080/post/{id}
		var card_schema common.CardSchemaIdBinding
		if err := c.ShouldBindUri(&card_schema); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "could not get post id",
				"msg":   err.Error(),
			})
			return
		}

		schema, err := database.GetCardSchema(card_schema.Id)
		if err != nil {
			log.Warn().Msgf("could not get post from DB: %v", err)
			c.JSON(http.StatusNotFound, gin.H{
				"error": "post id not found",
				"msg":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, common.CardSchema{
			Uuid:   card_schema.Id,
			Cards:  schema.Cards,
			Title:  schema.Title,
			Schema: schema.Schema,
		})
	}
}

// @Summary      Get list of card schemas
// @Description  Retrieves a paginated list of card schemas.
// @Tags         card_schema
// @Produce      json
// @Param        offset query int false "Pagination offset" default(0)
// @Param        limit query int false "Pagination limit (0 means no limit)" default(0)
// @Success      200 {object} GetSchemaasResponse
// @Failure      400 {object} common.ErrorResponse "Invalid query parameters"
// @Failure      500 {object} common.ErrorResponse "Server error"
// @Router       /card-schemas [get]
func getSchemasHandler(database database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lee offset y limit de la query (?offset=0&limit=10)
		// Un valor de 0 para el límite significa "sin límite".
		offsetStr := c.DefaultQuery("offset", "0")
		limitStr := c.DefaultQuery("limit", "0")

		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset parameter"})
			return
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
			return
		}

		// Si no se especifica un límite, se obtienen todos los posts.
		// Si se especifica, se usa para la paginación.
		schemas, err := database.GetCardSchemas(limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, GetSchemaasResponse{Schemas: schemas})
	}
}

func checkSchemaValues(add_schema_request AddCardSchemaRequest) error {

	if add_schema_request.JsonSchema == "" {
		return fmt.Errorf("`schema` cannot be empty")
	}
	if add_schema_request.JsonTitle == "" {
		return fmt.Errorf("`title` cannot be empty")
	}

	schema_compiler := jsonschema.NewCompiler()
	_, err := schema_compiler.Compile([]byte(add_schema_request.JsonSchema))

	if err != nil {
		return fmt.Errorf("`schema` is invalid: %v", err)
	}

	return nil
}

// @Summary      Delete a card schema
// @Description  Deletes a card schema by its ID.
// @Tags         card_schema
// @Accept       json
// @Produce      json
// @Param        schema body DeleteSchemaBinding true "ID of the card schema to delete"
// @Success      200 {object} map[string]string "Deleted schema ID"
// @Failure      400 {object} common.ErrorResponse "Invalid request or deletion error"
// @Router       /card-schemas [delete]
func deleteCardSchemaHandler(database database.Database) func(*gin.Context) {
	return func(c *gin.Context) {
		var delete_schema_request DeleteSchemaBinding
		decoder := json.NewDecoder(c.Request.Body)
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&delete_schema_request)
		if err != nil {
			log.Warn().Msgf("could not delete post: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
				"msg":   err.Error(),
			})
			return
		}
		err = database.DeleteCardSchema(delete_schema_request.Id)
		if err != nil {
			log.Error().Msgf("failed to delete post: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "could not delete post",
				"msg":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id": delete_schema_request.Id,
		})
	}
}

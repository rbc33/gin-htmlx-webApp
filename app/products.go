package app

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rbc33/gocms/common"
	"github.com/rbc33/gocms/database"
	"github.com/rbc33/gocms/views"
)

func productHandler(c *gin.Context, db database.Database) ([]byte, error) {

	// Used only for destructing the URI params
	var params struct {
		Schema string `uri:"schema" binding:"required"`
		Limit  int    `uri:"limit"`
		Page   int    `uri:"page"`
	}

	err := c.ShouldBindUri(&params)

	if err != nil {
		return []byte{}, fmt.Errorf("could not bind url params: %v", err)
	}

	if (params.Limit == 0) && (params.Page != 0) {
		return []byte{}, fmt.Errorf("card limit is 0 but pages is %d", params.Page)
	}

	if (params.Page == 0) && (params.Limit != 0) {
		return []byte{}, fmt.Errorf("card page is 0 but limit is %d", params.Limit)
	}

	limit := params.Limit
	page := params.Page
	if (params.Limit == 0) && (params.Page == 0) {
		limit = 10
		page = 0
	}

	cards, err := db.GetCards(params.Schema, int(limit), int(page))
	if err != nil {
		return []byte{}, fmt.Errorf("could not get cards: %v", err)
	}

	// TODO : this isn't very efficient as we transform
	// TODO : the card data to a JSON string, only to
	// TODO : deserialise it into a map of interface
	cards_data := make([]map[string]interface{}, 0)
	for _, card := range cards {
		var card_data map[string]interface{}
		err = json.Unmarshal([]byte(card.Content), &card_data)
		if err != nil {
			return []byte{}, fmt.Errorf("could not parse card json")
		}
		// Añade el campo image al mapa
		card_data["image"] = card.Image
		cards_data = append(cards_data, card_data)
	}

	return renderHtml(c, views.MakeProductPage(common.Settings.AppNavbar.Links, cards_data, common.Settings.AppNavbar.Dropdowns))
}

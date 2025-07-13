package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rbc33/gocms/common"
	"github.com/rbc33/gocms/database"
)

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func LoginHandler(db database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var input LoginInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		u := common.User{}

		u.Username = input.Username
		u.Password = input.Password

		token, err := common.LoginCheck(u.Username, u.Password, db)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})

	}
}

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary      Create new User
// @Description  Adds a new User to the database.
// @Tags         register
// @Accept       json
// @Produce      json
// @Param        post body RegisterInput true "Post to add"
// @Success      201 {object} PostIdResponse
// @Failure      400 {object} common.ErrorResponse "Invalid request body or missing data"
// @Router       /post [post]

func CreateRegisterHandler(db database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var input RegisterInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		u := common.User{}

		u.Username = input.Username
		u.Password = input.Password

		_, err := u.SaveUser(db)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "registration success"})

	}
}

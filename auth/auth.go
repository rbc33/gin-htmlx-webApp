package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rbc33/gocms/common"
	"github.com/rbc33/gocms/database"
	"github.com/rbc33/gocms/utils/token"
)

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary      Get current user
// @Description  Returns the currently authenticated user based on JWT token.
// @Tags         auth
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} common.User
// @Failure      400 {object} common.ErrorResponse
// @Router       /user [get]
func GetCurrentUserHandler(db database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		user_id, err := token.ExtractTokenID(c)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		u, err := db.GetUserById(user_id)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		u.Password = ""

		c.JSON(http.StatusOK, gin.H{"message": "success", "data": u})
	}
}

type TokenResponse struct {
	Token string `json:"token"`
}

// @Summary      Login user
// @Description  Authenticates user and returns a JWT token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials body LoginInput true "User credentials"
// @Success      200 {object} TokenResponse
// @Failure      400 {object} common.ErrorResponse
// @Router       /login [post]
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
type RegisterResponse struct {
	Id int `json:"user_id"`
}

// // @Summary      Create new User
// // @Description  Adds a new User to the database.
// // @Tags         auth
// // @Accept       json
// // @Produce      json
// // @Param        post body RegisterInput true "User to create"
// // @Success      201 {object} RegisterResponse
// // @Failure      400 {object} common.ErrorResponse
// // @Router       /register [post]
// comment
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

		c.JSON(http.StatusConflict, gin.H{"message": "registration success"})

	}
}

package middleware

import (
	"api/service/auth"
	"api/service/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoggedMiddleware(c *gin.Context) {

	token, err := c.Cookie("token")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, responses.APIResponse{Success: false, Data: "Token não encontrado!"})
		return
	}

	validation, user_id := auth.ValidarToken(token)

	if !validation {
		c.AbortWithStatusJSON(http.StatusUnauthorized, responses.APIResponse{Success: false, Data: "Token não é válido!"})
		return
	} else {
		c.Set("user_id", user_id)
		c.Next()
	}

}

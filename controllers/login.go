package controllers

import (
	"api/service/auth"
	"api/service/db"
	"api/service/models"
	"api/service/responses"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type BodyLogin struct {
	Email    string
	Password string
}

func LoginController(c *gin.Context) {

	var loginUser BodyLogin

	err := c.BindJSON(&loginUser)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, responses.APIResponse{Success: false, Data: "Informações inválidas na requisição!"})
		return
	}

	var user models.User

	response := db.DB.Find(&user, models.User{Email: loginUser.Email})

	if response.Error != nil || response.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, responses.APIResponse{Success: false, Data: "Erro ao encontrar usuário!"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.Senha), []byte(loginUser.Password))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, responses.APIResponse{Success: false, Data: "Senha invalida!"})
		return
	} else {
		token, err := auth.GerarToken(user.ID)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, responses.APIResponse{Success: false, Data: "Erro ao gerar token!"})
			return
		}
		fmt.Println("dominimo", os.Getenv("DOMAIN"))
		c.SetCookie("token", token, 3600*24*7, "/", os.Getenv("DOMAIN"), os.Getenv("DOMAIN") != "localhost", os.Getenv("DOMAIN") != "localhost")

		c.JSON(http.StatusAccepted, responses.APIResponse{Success: true, Data: token})

	}
}

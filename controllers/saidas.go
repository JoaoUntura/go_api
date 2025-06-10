package controllers

import (
	"api/service/db"
	"api/service/models"
	"api/service/responses"
	"api/service/utils"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSaidas(c *gin.Context) {
	userID := utils.HandleUserID(c)

	if userID == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, responses.APIResponse{Success: false, Data: "Você não tem autorização!"})
		return
	}

	limitParam := c.DefaultQuery("limit", "10")
	pageParam := c.DefaultQuery("page", "1")

	limit, err1 := strconv.Atoi(limitParam)
	page, err2 := strconv.Atoi(pageParam)

	// Validação básica
	if err1 != nil || err2 != nil || limit <= 0 || page <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, responses.APIResponse{Success: false, Data: "Parâmetros de paginação inválidos"})
		return
	}

	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	var totalItems int64
	db.DB.Model(&models.Saida{}).Where("user_id = ?", userID).Count(&totalItems)

	var saidas []models.Saida
	result := db.DB.
		Where("user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Order("date DESC").
		Find(&saidas)

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, responses.APIResponse{Success: false, Data: result.Error.Error()})
		return
	}

	// Cálculo do total de páginas
	totalPages := int((totalItems + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, responses.APIResponsePagination{Success: true, Data: saidas, Page: page, Limit: limit, TotalItems: totalItems, TotalPages: totalPages})
}

func PostSaidas(c *gin.Context) {

	userID := utils.HandleUserID(c)

	if userID == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, responses.APIResponse{Success: false, Data: "Você não tem autorização!"})
		return
	}

	var newSaida models.Saida

	err := c.ShouldBindJSON(&newSaida)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, responses.APIResponse{Success: false, Data: "Requisição Inválida"})
		return
	}

	if !utils.IsOwner(newSaida.UserID, userID) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, responses.APIResponse{Success: false, Data: "Você não tem autorização para esse usuário!"})
		return
	}
	newSaida.UserID = userID
	result := db.DB.Create(&newSaida)

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, responses.APIResponse{Success: false, Data: result.Error})
		return
	}

	c.JSON(http.StatusCreated, responses.APIResponse{Success: true, Data: newSaida.ID})

}

func DeleteSaidas(c *gin.Context) {

	userID := utils.HandleUserID(c)

	if userID == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, responses.APIResponse{Success: false, Data: "Você não tem autorização!"})
		return
	}

	param := c.Param("id")

	idSaida, err := strconv.ParseUint(param, 10, 64)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, responses.APIResponse{Success: false, Data: "ID inválido"})
		return
	}

	idSaidaUint := uint(idSaida)

	var saidaDeletar models.Saida

	resultFindDelete := db.DB.First(&saidaDeletar, idSaidaUint)

	if resultFindDelete.Error != nil || resultFindDelete.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, responses.APIResponse{Success: false, Data: "Registro não encontrado"})
		return
	}

	if !utils.IsOwner(saidaDeletar.UserID, userID) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, responses.APIResponse{Success: false, Data: "Você não tem autorização para esse usuário!"})
		return
	}

	result := db.DB.Delete(&saidaDeletar)

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, responses.APIResponse{Success: false, Data: result.Error})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Success: true, Data: "Deletado com sucesso!"})

}

func PutSaidas(c *gin.Context) {

	userID := utils.HandleUserID(c)

	if userID == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, responses.APIResponse{Success: false, Data: "Você não tem autorização!"})
		return
	}

	param := c.Param("id")

	idSaida, err := strconv.ParseUint(param, 10, 64)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, responses.APIResponse{Success: false, Data: "ID inválido"})
		return
	}

	idSaidaUint := uint(idSaida)

	var saidaValidation models.Saida

	resultValidation := db.DB.First(&saidaValidation, idSaidaUint)

	if resultValidation.Error != nil || resultValidation.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, responses.APIResponse{Success: false, Data: "Registro não encontrado"})
		return
	}

	if !utils.IsOwner(saidaValidation.UserID, userID) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, responses.APIResponse{Success: false, Data: "Você não tem autorização para esse usuário!"})
		return
	}

	var saidaUpdate models.Saida

	errJson := c.ShouldBindJSON(&saidaUpdate)

	if errJson != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, responses.APIResponse{Success: false, Data: "Requisição Inválida"})
		return
	}

	result := db.DB.Where(&saidaValidation).Updates(saidaUpdate)

	if result.Error != nil || result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, responses.APIResponse{Success: false, Data: "Erro ao atualizar!"})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Success: true, Data: "Atualizado com sucesso!"})
}

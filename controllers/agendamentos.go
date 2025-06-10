package controllers

import (
	"api/service/db"
	"api/service/models"
	"api/service/responses"
	"strconv"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestBody struct {
	Agendamento models.Agendamento `json:"agendamento"`
	Servico     *[]models.Servico  `json:"servico"`
	PaymentDbId *models.Pagamento  `json:"paymentDbId"`
}

func GetAgendamentos(c *gin.Context) {
	param := c.Param("id")
	id_user, err := strconv.ParseInt(param, 0, 36)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{Success: false, Data: "Invalid user ID"})
		return
	}
	idUserUint := uint(id_user)

	var agendamentos []models.Agendamento
	result := db.DB.Where(models.Agendamento{IDUser: &idUserUint}).Preload("Pagamento").Preload("Funcionario").Preload("ServicosAgendado.Servico").Find(&agendamentos)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, responses.APIResponse{Success: false, Data: result.Error})
	}

	c.JSON(http.StatusOK, responses.APIResponse{Success: true, Data: agendamentos})

}

func PostAgendamentos(c *gin.Context) {
	var reqBody RequestBody

	err := c.ShouldBind(&reqBody)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, responses.APIResponse{Success: false, Data: "Agendamento Inválido"})
	}

	newAgendamento := reqBody.Agendamento

	if !newAgendamento.Indisponivel && reqBody.PaymentDbId == nil {
		total := 0.0
		for _, value := range *reqBody.Servico {
			var service models.Servico
			db.DB.Find(&service, value.ID)
			total += *service.Preco
		}

		date := time.Now()
		status := "pending"
		statusDetail := "Esperando Confirmação"
		externo := false
		paymentStruct := models.Pagamento{DateCreated: &date, Status: &status, StatusDetail: &statusDetail, UserID: newAgendamento.IDUser, TransactionAmount: &total, Externo: &externo}
		result := db.DB.Create(&paymentStruct)
		if result.Error != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, responses.APIResponse{Success: false, Data: result.Error})
		}

		newAgendamento.Pagamento.ID = paymentStruct.ID
	} else if !newAgendamento.Indisponivel {
		result := db.DB.Find(reqBody.PaymentDbId)

		if result.RowsAffected > 0 {
			newAgendamento.Pagamento.ID = reqBody.PaymentDbId.ID
		}

	}

	result := db.DB.Create(&newAgendamento)

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, responses.APIResponse{Success: false, Data: result.Error})
	}

	if len(*reqBody.Servico) > 0 {
		for _, value := range *reqBody.Servico {
			result := db.DB.Create(models.AgendamentoServico{IdAgendamento: int(newAgendamento.ID), IdServico: value.ID})
			if result.Error != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, responses.APIResponse{Success: false, Data: "Erro ao buscar serviço"})
				return
			}
		}
	}

	c.JSON(http.StatusCreated, responses.APIResponse{Success: true, Data: newAgendamento.ID})

}

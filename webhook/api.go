package webhook

import (
	"github.com/edimarlnx/rinha-de-backend-2024-q1-edimarlnx/app"
	"github.com/edimarlnx/rinha-de-backend-2024-q1-edimarlnx/types"
	"github.com/edimarlnx/rinha-de-backend-2024-q1-edimarlnx/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func wsTransacoes(e *gin.RouterGroup, transactionController *app.TransactionController) {
	e.POST("/clientes/:id/transacoes", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var transacao types.Transacao
		err := c.BindJSON(&transacao)
		if err != nil {
			msg := "Erro ao realizar o parser dos dados da transação."
			utils.Log.WithField("err", err).Error(msg)
			c.AbortWithStatusJSON(http.StatusBadRequest, &types.Response{
				Status: http.StatusBadRequest,
				Error:  msg,
			})
			return
		}
		saldo, errWithCode := transactionController.Transacao(id, transacao)
		if errWithCode != nil {
			msg := errWithCode.Error.Error()
			utils.Log.Error(msg)
			c.AbortWithStatusJSON(errWithCode.Code, &types.Response{
				Status: errWithCode.Code,
				Error:  msg,
			})
			return
		}
		c.IndentedJSON(http.StatusOK, saldo)
	})
}

func wsExtrato(e *gin.RouterGroup, transactionController *app.TransactionController) {
	e.GET("/clientes/:id/extrato", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		extrato, errWithCode := transactionController.Extrato(id)
		if errWithCode != nil {
			msg := errWithCode.Error.Error()
			utils.Log.Error(msg)
			c.AbortWithStatusJSON(errWithCode.Code, &types.Response{
				Status: errWithCode.Code,
				Error:  msg,
			})
			return
		}
		c.IndentedJSON(http.StatusOK, extrato)
	})
}

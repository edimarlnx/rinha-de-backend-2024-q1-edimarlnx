package webhook

import (
	"github.com/edimarlnx/rinha-de-backend-2024-q1-edimarlnx/app"
	"github.com/gin-gonic/gin"
)

func CreateWebhook(transactionController *app.TransactionController) (*gin.Engine, *gin.RouterGroup) {
	engine := gin.Default()
	routerGroup := engine.Group("/")
	wsHealth(routerGroup)
	wsTransacoes(routerGroup, transactionController)
	wsExtrato(routerGroup, transactionController)
	return engine, routerGroup
}

package main

import (
	"github.com/edimarlnx/rinha-de-backend-2024-q1-edimarlnx/app"
	"github.com/edimarlnx/rinha-de-backend-2024-q1-edimarlnx/utils"
	"github.com/edimarlnx/rinha-de-backend-2024-q1-edimarlnx/webhook"
)

var (
	listenAddress string
	dbUri         string
)

func init() {
	listenAddress = utils.GetEnv("LISTEN_ADDRESS", ":8080")
	dbUri = utils.GetEnv("DB_URI", "postgres://postgres:postgres@localhost:55432/rinha-backend?sslmode=disable")
}

func main() {

	transactionController := app.New(dbUri)
	server, _ := webhook.CreateWebhook(transactionController)
	err := server.Run(listenAddress)
	if err != nil {
		utils.Log.WithError(err).Error("Error starting server")
		return
	}
}

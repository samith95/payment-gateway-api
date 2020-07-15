package main

import (
	"payment-gateway-api/api/app"
	"payment-gateway-api/api/config"
	"payment-gateway-api/api/data_access"
)

func main() {
	err := data_access.Db.Setup(config.DbStoreFilePath)
	if err != nil {
		panic("failed to connect to db: " + err.Error())
	}
	defer data_access.Db.Close()
	app.RunApp()
}

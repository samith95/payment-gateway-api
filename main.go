package main

import (
	"payment-gateway-api/api/app"
	"payment-gateway-api/api/config"
	"payment-gateway-api/api/data_access"
)

func main() {
	data_access.Db.Setup(config.DbStoreFilePath)
	defer data_access.Db.Close()
	app.RunApp()
}

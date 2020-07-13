package main

import (
	"payment-gateway-api/api/app"
	"payment-gateway-api/api/data_access"
)

func main() {
	data_access.Db.Init()
	defer data_access.Db.Close()
	app.RunApp()
}

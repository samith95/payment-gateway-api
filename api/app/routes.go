package app

import (
	"payment-gateway-api/api/controllers/authorisation_controller"
)

func routes() {
	router.POST("/authorize", authorisation_controller.HandleAuthorisationRequest)
	//Auxiliary handler for troubleshooting
	//TODO: remove once application is done
	router.GET("/authorize/getrecords", authorisation_controller.HandleGetAuthsRecords)
}

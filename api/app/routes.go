package app

import (
	"payment-gateway-api/api/controllers/authorisation_controller"
	"payment-gateway-api/api/controllers/capture_controller"
	"payment-gateway-api/api/controllers/refund_controller"
	"payment-gateway-api/api/controllers/void_controller"
)

func routes() {
	router.POST("/authorize", authorisation_controller.HandleAuthorisationRequest)
	router.PATCH("/void", void_controller.HandleVoidRequest)
	router.PATCH("/capture", capture_controller.HandleCaptureRequest)
	router.PATCH("/refund", refund_controller.HandleRefundRequest)
}

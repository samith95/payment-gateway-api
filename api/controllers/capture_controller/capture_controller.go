package capture_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"payment-gateway-api/api/domain/capture_domain"
	"payment-gateway-api/api/domain/error_domain"
	"payment-gateway-api/api/services/capture_service"
)

func HandleCaptureRequest(c *gin.Context) {
	request := capture_domain.CaptureRequest{}

	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, error_domain.GatewayError{
			Code:  http.StatusBadRequest,
			Error: "request body is invalid",
		})
		return
	}

	result, apiError := capture_service.CaptureService.CaptureTransactionAmount(request)
	if apiError != nil {
		c.JSON(apiError.Status(), apiError)
		return
	}
	c.JSON(http.StatusOK, result)
}

package refund_controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"payment-gateway-api/api/domain/error_domain"
	"payment-gateway-api/api/domain/refund_domain"
	"payment-gateway-api/api/services/refund_service"
)

//HandleRefundRequest handles request for the refund endpoint
func HandleRefundRequest(c *gin.Context) {
	request := refund_domain.RefundRequest{}

	err := c.BindJSON(&request)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, error_domain.GatewayError{
			Code:  http.StatusBadRequest,
			Error: "request body is invalid",
		})
		return
	}

	result, apiError := refund_service.RefundService.RefundTransactionAmount(request)
	if apiError != nil {
		c.JSON(apiError.Status(), apiError)
		return
	}
	c.JSON(http.StatusOK, result)
}

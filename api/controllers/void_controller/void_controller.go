package void_controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"payment-gateway-api/api/domain/error_domain"
	"payment-gateway-api/api/domain/void_domain"
	"payment-gateway-api/api/services/void_service"
)

//HandleVoidRequest handles request for the void endpoint
func HandleVoidRequest(c *gin.Context) {
	request := void_domain.VoidRequest{}

	err := c.BindJSON(&request)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, error_domain.GatewayError{
			Code:  http.StatusBadRequest,
			Error: "request body is invalid",
		})
		return
	}

	result, apiError := void_service.VoidService.VoidTransaction(request)
	if apiError != nil {
		c.JSON(apiError.Status(), apiError)
		return
	}
	c.JSON(http.StatusOK, result)
}

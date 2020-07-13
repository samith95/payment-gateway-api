package app

import (
	"github.com/gin-gonic/gin"
	"log"
)

var (
	router = gin.Default()
)

//RunApp will run constantly until the application is shutdown
func RunApp() {
	routes()

	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

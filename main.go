package main

import (
	"log"

	"vinam/handlers"

	"github.com/gin-gonic/gin"
)

func init() {

	handlers.LoadEnvVariables()
	handlers.Connect()
}
func main() {
	r := gin.Default()
	r.POST("/", handlers.Home)
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err.Error())
	}

}

package main

import (
	"PROTOBUF/config"
	"PROTOBUF/controller"

	"github.com/gin-gonic/gin"
)

func main() {

	config.DbConnect()
	r := gin.Default()

	r.POST("/add", controller.CreateUserHandle)
	r.GET("/get", controller.GetUser)

	r.Run(":8080")

}

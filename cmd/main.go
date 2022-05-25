package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/user", createUser)
	router.POST("/login", auth_middleware, loginUser)
	router.POST("/findPartner", auth_middleware, findPartner)

	router.Run("localhost:8080")
}

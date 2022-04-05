package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/user", createUser)
	router.POST("/login", loginUser)

	router.Run("localhost:8080")
}

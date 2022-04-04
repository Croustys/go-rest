package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/user", createUser)

	router.Run("localhost:8080")
}

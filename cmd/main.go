package main

import (
	//"github.com/Croustys/go-rest/pkg/db"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/user", createUser)
	router.POST("/login", loginUser)
	router.POST("/findPartner", findPartner)

	//db.Insert("asd")

	router.Run("localhost:8080")
}

package main

import (
	"log"

	"snaptalk-api/src/db"

	"github.com/gin-gonic/gin"
)

func createUser(c *gin.Context) {
	var user db.User

	if c.ShouldBind(&user) != nil {
		log.Fatalf("Error")
	}
	db.Insert(user.Name, user.Email)
}

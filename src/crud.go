package main

import (
	"log"

	"snaptalk-api/src/db"

	"github.com/gin-gonic/gin"
)

func createUser(c *gin.Context) {
	var user db.User

	if c.ShouldBind(&user) != nil {
		log.Fatalf("Error createUser")
	}
	if db.Insert(user.Name, user.Password, user.Email) {
		c.JSON(200, gin.H{
			"message": "user created",
		})
	}
}
func loginUser(c *gin.Context) {
	var user db.User

	if c.ShouldBind(&user) != nil {
		log.Fatalf("Error loginUser")
	}
	if db.Login(user.Email, user.Password) {
		c.JSON(200, gin.H{
			"message": "login successfull",
		})
	}
}

package main

import (
	"os"

	"github.com/Croustys/go-rest/pkg/chat"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func main() {
	router := gin.Default()
	setup_env()

	gothic.Store = setup_store()

	goth.UseProviders(
		google.New(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"), "http://localhost:3000/auth/google/callback", "email", "profile"),
	)

	hub := chat.CreateHub()
	go hub.Run()

	router.POST("/user", createUser)
	router.POST("/login", auth_middleware, loginUser)
	router.POST("/findPartner", auth_middleware, findPartner)
	router.GET("/auth/:provider/callback", oauthLogin)
	router.GET("/auth/:provider", authProvider)
	router.GET("/logout/:provider", logoutProvider)
	router.GET("/ws", func(c *gin.Context) {
		chat.ServeChat(hub, c.Writer, c.Request)
	})

	router.Run("localhost:3000")
}

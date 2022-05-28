package main

import (
	"os"
	"time"

	"github.com/Croustys/go-rest/pkg/chat"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

var Hub *chat.Hub

func main() {
	router := gin.Default()
	setup_env()

	gothic.Store = setup_store()

	goth.UseProviders(
		google.New(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"), "http://localhost:3000/auth/google/callback", "email", "profile"),
	)

	Hub := chat.CreateHub()
	go Hub.Run()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000/register", "http://localhost:3000", "http://localhost:3000/login", "*"},
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"},
		ExposeHeaders:    []string{"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	router.POST("/register", createUser)
	router.POST("/login", auth_middleware, loginUser)
	router.POST("/findPartner", auth_middleware, findPartner)
	router.GET("/auth/:provider/callback", oauthLogin)
	router.GET("/auth/:provider", authProvider)
	router.GET("/logout/:provider", logoutProvider)
	router.GET("/ws", auth_middleware, chatHandler)

	router.Run("localhost:3001")
}

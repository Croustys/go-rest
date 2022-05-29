package main

import (
	"os"
	"time"

	"github.com/Croustys/go-rest/pkg/chat"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
	Hub               *chat.Hub
)

func main() {
	router := gin.Default()
	setup_env()

	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:3001/auth/google/callback",
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

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
	router.GET("/auth/google", oauthLogin)
	router.GET("/auth/google/callback", authProvider)
	router.GET("/ws", auth_middleware, chatHandler)

	router.Run("localhost:3001")
}

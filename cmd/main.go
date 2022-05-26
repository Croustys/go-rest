package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func setSecretToken() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err.Error())
	}

	key = os.Getenv("SESSION_KEY")
}

var key string = ""

func main() {
	router := gin.Default()

	setSecretToken()
	maxAge := 86400 * 30
	isProd := false

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd

	gothic.Store = store

	goth.UseProviders(
		google.New("155971532808-osgvjdl31kht2bv9ifot11rsu99ao116.apps.googleusercontent.com", "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", "http://localhost:3000/auth/google/callback", "email", "profile"),
	)

	router.POST("/user", createUser)
	router.POST("/login", auth_middleware, loginUser)
	router.POST("/findPartner", auth_middleware, findPartner)
	router.GET("/auth/:provider/callback", oauthLogin)
	router.GET("/auth/:provider", authProvider)
	router.GET("/logout/:provider", logoutProvider)

	router.Run("localhost:3000")
}

type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

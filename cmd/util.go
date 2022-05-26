package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

func set_request_provider(c *gin.Context) {
	//Building query beacuse gothic doesn't handle gin request
	q := c.Request.URL.Query()
	q.Add("provider", c.Param("provider"))
	c.Request.URL.RawQuery = q.Encode()
}
func setup_store() *sessions.CookieStore {
	maxAge := 86400 * 30
	isProd := false

	tempStore := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

	tempStore.MaxAge(maxAge)
	tempStore.Options.Path = "/"
	tempStore.Options.HttpOnly = true
	tempStore.Options.Secure = isProd

	return tempStore
}

func setup_env() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err.Error())
	}
}

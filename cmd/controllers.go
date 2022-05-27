package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/Croustys/go-rest/pkg/auth"
	"github.com/Croustys/go-rest/pkg/db"
	"github.com/markbates/goth/gothic"

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
func auth_middleware(c *gin.Context) {
	isAuthorized := auth.AuthUser(c)
	isLoginRequest := c.Request.URL.String() == "/login"

	if isAuthorized && isLoginRequest {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"message": "Authorized",
		})
	} else if isAuthorized || isLoginRequest {
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unsuccessful Authorization",
		})
	}
}
func loginUser(c *gin.Context) {
	var user db.User

	if c.ShouldBind(&user) != nil {
		log.Fatalf("Error loginUser")
	}

	if user.Email != "" && db.Login(user.Email, user.Password) {
		auth.GenerateToken(c)

		c.JSON(http.StatusOK, gin.H{
			"message": "login successful",
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unsuccessful Auth",
		})
	}
}
func findPartner(c *gin.Context) {
	var user db.User

	if c.ShouldBind(&user) != nil {
		log.Fatalf("Error findPartner")
	}
	//business logic matching partners
	c.JSON(200, gin.H{
		"partnerId": user.ID.String(),
	})
}

func oauthLogin(c *gin.Context) {
	set_request_provider(c)

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		return
	}
	t, _ := template.New("foo").Parse(UserTemplate)
	t.Execute(c.Writer, user)
}
func authProvider(c *gin.Context) {
	set_request_provider(c)

	if gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request); err == nil {
		t, _ := template.New("foo").Parse(UserTemplate)
		t.Execute(c.Writer, gothUser)
	} else {
		gothic.BeginAuthHandler(c.Writer, c.Request)
	}
}
func logoutProvider(c *gin.Context) {
	gothic.Logout(c.Writer, c.Request)
	c.Writer.Header().Set("Location", "/")
	c.Writer.WriteHeader(http.StatusTemporaryRedirect)
}

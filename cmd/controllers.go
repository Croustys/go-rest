package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Croustys/go-rest/pkg/auth"
	"github.com/Croustys/go-rest/pkg/chat"
	"github.com/Croustys/go-rest/pkg/db"

	"github.com/gin-gonic/gin"
)

func createUser(c *gin.Context) {
	var user db.User

	if c.ShouldBind(&user) != nil {
		log.Fatalf("Error createUser")
	}
	if db.Insert(user.Name, user.Password, user.Email) {
		c.JSON(http.StatusOK, gin.H{
			"message": "user created",
		})
	}
}
func findUser(c *gin.Context) {
	email := c.Param("email")
	user := db.GetUser(email)
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
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

type userData struct {
	Id             int
	Email          string
	Verified_email bool
	Picture        string
}

//var accesToken string

func authProvider(c *gin.Context) {
	log.Println(c.Request.FormValue("code"))
	content, err := getUserInfo(c.Request.FormValue("state"), c.Request.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(c.Writer, c.Request, "/", http.StatusTemporaryRedirect)
		return
	}

	var user userData
	json.Unmarshal(content, &user)

	if !db.SaveOauth(user.Email) {
		log.Println("Couldn't save user")
		return
	}
	url := "http://localhost:3000/account" /* ?token=" + accesToken */
	if err != nil {
		log.Println(err)
	}
	http.Redirect(c.Writer, c.Request, url, http.StatusSeeOther)
}
func oauthLogin(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL("pseudo-random")
	http.Redirect(c.Writer, c.Request, url, http.StatusTemporaryRedirect)
}
func getUserInfo(state string, code string) ([]byte, error) {
	if state != "pseudo-random" {
		return nil, fmt.Errorf("invalid oauth state")
	}
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	//accesToken = token.AccessToken
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	return contents, nil
}

/* func logoutProvider(c *gin.Context) {
	gothic.Logout(c.Writer, c.Request)
	c.Writer.Header().Set("Location", "/")
	c.Writer.WriteHeader(http.StatusTemporaryRedirect)
} */
func chatHandler(c *gin.Context) {
	chat.ServeChat(Hub, c.Writer, c.Request)
}

package auth

import (
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var secret_token string

func AuthUser(c *gin.Context) bool {
	tok, err := c.Cookie("Auth_token")
	if err != nil {
		return false
	}

	return verifyToken(tok)
}

func verifyToken(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			log.Println(ok)
		}
		setSecretToken()
		return []byte(secret_token), nil
	})

	if err != nil {
		log.Println(err)
	}

	return token.Valid
}

func GenerateToken(c *gin.Context) {
	new_token := jwt.New(jwt.SigningMethodHS256)

	setSecretToken()
	tokenString, err := new_token.SignedString([]byte(secret_token))
	if err != nil {
		log.Println(err)
	}

	c.SetCookie("Auth_token", tokenString, 86400, "", "", false, true)
}

func setSecretToken() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err.Error())
	}

	secret_token = os.Getenv("JWT_TOKEN")
}

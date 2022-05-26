package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

var userTemplate = `
<p><a href="/logout/{{.Provider}}">logout</a></p>
<p>Name: {{.Name}} [{{.LastName}}, {{.FirstName}}]</p>
<p>Email: {{.Email}}</p>
<p>NickName: {{.NickName}}</p>
<p>Location: {{.Location}}</p>
<p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
<p>Description: {{.Description}}</p>
<p>UserID: {{.UserID}}</p>
<p>AccessToken: {{.AccessToken}}</p>
<p>ExpiresAt: {{.ExpiresAt}}</p>
<p>RefreshToken: {{.RefreshToken}}</p>
`

func oauthLogin(c *gin.Context) {
	q := c.Request.URL.Query()
	q.Add("provider", c.Param("provider"))
	c.Request.URL.RawQuery = q.Encode()

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		fmt.Fprintln(c.Writer, err)
		return
	}
	t, _ := template.New("foo").Parse(userTemplate)
	t.Execute(c.Writer, user)
}

func authProvider(c *gin.Context) {
	//Building query beacuse gothic doesn't handle gin request
	q := c.Request.URL.Query()
	q.Add("provider", c.Param("provider"))
	c.Request.URL.RawQuery = q.Encode()

	if gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request); err == nil {
		t, _ := template.New("foo").Parse(userTemplate)
		t.Execute(c.Writer, gothUser)
	} else {
		gothic.BeginAuthHandler(c.Writer, c.Request)
	}
}

func index(c *gin.Context) {
	m := make(map[string]string)
	m["google"] = "Google"
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	providerIndex := &ProviderIndex{Providers: keys, ProvidersMap: m}
	t, _ := template.New("foo").Parse(indexTemplate)
	t.Execute(c.Writer, providerIndex)
}

func logoutProvider(c *gin.Context) {
	gothic.Logout(c.Writer, c.Request)
	c.Writer.Header().Set("Location", "/")
	c.Writer.WriteHeader(http.StatusTemporaryRedirect)
}

var indexTemplate = `{{range $key,$value:=.Providers}}
    <p><a href="/auth/{{$value}}">Log in with {{index $.ProvidersMap $value}}</a></p>
{{end}}`

func main() {
	router := gin.Default()

	key := "Secret-session-key" // Replace with your SESSION_SECRET or similar
	maxAge := 86400 * 30        // 30 days
	isProd := false             // Set to true when serving over https

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
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
	router.GET("/", index)

	router.Run("localhost:3000")
}

type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

func GoogleAuthCallbackfunc(c *gin.Context) {

	provider := c.Param("provider")
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "provider", provider))
	fmt.Println("Request", c.Request)
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	fmt.Println("USER", user)
	if err != nil {
		fmt.Println("ERROR: ", c.Writer, c.Request)
		return
	}

	fmt.Println(user)

	c.Redirect(http.StatusAccepted, "http://localhost:8080/")
}

func OAuthLogout(c *gin.Context) {

	provider := c.Param("provider")

	type ContextKey string
	const ProviderKey ContextKey = "provider"

	c.Request = c.Request.WithContext(context.WithValue(c, ProviderKey, provider))
	fmt.Println(c.Request)

	gothic.Logout(c.Writer, c.Request)
	c.Writer.Header().Set("Location", "/")
	c.Writer.WriteHeader(http.StatusTemporaryRedirect)
}

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

func OAuthProvider(c *gin.Context) {

	provider := c.Param("provider")
	fmt.Println("provider", provider)

	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "provider", provider))
	fmt.Println(c.Request)

	// try to get the user without re-authenticating
	if gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request); err == nil {
		t, _ := template.New("foo").Parse(userTemplate)
		t.Execute(c.Writer, gothUser)
	} else {
		fmt.Println("HERE")
		gothic.BeginAuthHandler(c.Writer, c.Request)
	}
}

func OAuthGithub(c *gin.Context) {
	githubProvider := github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), "http://localhost:8080/callback")
	goth.UseProviders(githubProvider)

	q := c.Request.URL.Query()
	q.Add("provider", "github")
	c.Request.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func OAuthGoogle(c *gin.Context) {
	google := google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "http://127.0.0.1:8080/callback")
	goth.UseProviders(google)

	q := c.Request.URL.Query()
	q.Add("provider", "google")
	c.Request.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func OAuthGithubCallback(c *gin.Context) {
	q := c.Request.URL.Query()
	q.Add("provider", "google")
	c.Request.URL.RawQuery = q.Encode()
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	res, err := json.Marshal(user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	fmt.Println(user)
	fmt.Println(res)
	jsonString := string(res)
	c.JSON(http.StatusOK, gin.H{
		"user": jsonString,
	})
}

package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)


type ContextKey string
const ProviderKey ContextKey = "provider"

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

	res, err := json.Marshal(user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	
	jsonString := string(res)

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"user": jsonString,
	})
}

func OAuthLogout(c *gin.Context) {

	provider := c.Param("provider")

	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "provider", provider))
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

	// try to get the user without re-authenticating
	if gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request); err == nil {
		t, _ := template.New("foo").Parse(userTemplate)
		t.Execute(c.Writer, gothUser)
	} else {
		fmt.Println("HERE")
		gothic.BeginAuthHandler(c.Writer, c.Request)
	}
}

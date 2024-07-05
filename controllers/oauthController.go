package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
)




func OAuthGithub(c *gin.Context) {
	githubProvider := github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), "http://localhost:8080/callback")
	goth.UseProviders(githubProvider)

	q := c.Request.URL.Query()
	q.Add("provider", "github")
	c.Request.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func OAuthGithubCallback(c *gin.Context) {
	q := c.Request.URL.Query()
	q.Add("provider", "github")
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
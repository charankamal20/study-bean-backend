package controllers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"study-bean/database"
	"study-bean/models"
	"study-bean/responses"
	"time"
)

type NewSessionBody struct {
	Uid                string `json:"uid"`
	SessionName        string `json:"session_name"`
	SessionDescription string `json:"session_description"`
}

func CreateNewSession(c *gin.Context) {

	body := &NewSessionBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": responses.UserIdNotFound,
		})
	}

	if body.Uid == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": responses.UserIdNotFound,
		})
	}

	NewSession := &models.Session{
		ID:                 primitive.NewObjectID(),
		SessionID:          primitive.NewObjectID().Hex(),
		AdminID:            body.Uid,
		SessionName:        body.SessionName,
		SessionDescription: body.SessionDescription,
		SessionPhoto:       "",
		NumberOfMembers:    1,
		Banned:             []string{},
		UpdatedAt:          time.Now(),
		ExpiresAfter:       time.Now().Add(24 * time.Second),
	}

	if _, err := primitive.ObjectIDFromHex(body.Uid); err != nil {
		NewSession.TempMemberList = []string{body.Uid}
		NewSession.Members = []string{}
	} else {
		NewSession.Members = []string{body.Uid}
		NewSession.TempMemberList = []string{}
	}

	err := database.CreateNewSession(NewSession)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": responses.SessionCreationFailed,
		})
		return
	}

	c.SetCookie("Session", NewSession.SessionID, 3600*24, "/", ".classikh.me", true, true)
	// c.SetCookie("Session", NewSession.SessionID, 3600*24, "/", "", false, true)
	c.JSON(http.StatusCreated, gin.H{
		"success":   true,
		"sessionId": NewSession.SessionID,
	})
}

type JoinSessionBody struct {
	Uid        string `json:"uid"`
	Session_id string `json:"session_id"`
}

func JoinSession(c *gin.Context) {

	body := &JoinSessionBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": responses.UserIdNotFound,
		})
	}

	if body.Uid == "" || body.Session_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": responses.UserIdNotFound,
		})
	}

	isLoggedIn := true
	if _, err := primitive.ObjectIDFromHex(body.Uid); err != nil {
		isLoggedIn = false
	}

	err := database.JoinSession(body.Uid, body.Session_id, isLoggedIn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	//! BUG: TIME SHOULD BE CALCULATED FOR THE COOKIE
	c.SetCookie("Session", body.Session_id, 3600*24, "/", ".classikh.me", true, true)
	// c.SetCookie("Session", body.Uid, 3600*24, "/", "", false, true)
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": responses.JoinedSessionSuccessfully,
	})
}

type SessionTodoBody struct	{
	Todo string `json:"todo"`
	Priority models.Priority `json:"priority"`
}

func CreateSessionTodo(c *gin.Context) {

}

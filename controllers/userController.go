package controllers

import (
	"net/http"
	"study-bean/database"
	"study-bean/models"
	"study-bean/responses"
	"study-bean/tokens"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// func getUserByName(email string) (*models.User, error) {
//     for i, u := range models.Users {
//         if u.Email == email {
//             return &models.Users[i], nil
//         }
//     }

//     return nil, errors.New("user not found")
// }

func GetAllUsers(context *gin.Context) {
	users, err := database.FindAllUsers()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"success" : true,
			"message": responses.DatabaseError,
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":   users,
	})
}

func SignUp(c *gin.Context) {
	// Get the email/pass off the body
	var body struct {
		Email    string
		Password string
		Username string
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": responses.UserNotFound,
			"success": false,
		})
		return
	}

	// Check is user with same email exists
	user, _ := database.FindUserByEmail(body.Email)
	if user != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message":   responses.EmailTaken,
		})
		return
	}

	// Check if user with same username exists
	user, _ = database.FindUserByUsername(body.Username)
	if user != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message":   responses.UsernameTaken,
		})
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": responses.FailedToHash,
		})
		return
	}

	refreshTokenString, err := tokens.GenerateNewRefreshToken()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": responses.FailedRefreshToken,
		})
		return
	}

	var newUser models.User

	newUser.Password = string(hash)
	newUser.Email = body.Email
	newUser.Username = body.Username
	newUser.ID = primitive.NewObjectID()
	newUser.User_ID = newUser.ID.Hex()
	newUser.RefreshToken = refreshTokenString
	newUser.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	err = database.AddUserToDatabase(newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message":   responses.TryAgain,
		})
		return
	}

	// Respond
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"user":    newUser,
	})
}

func Login(c *gin.Context) {
	// Get the email/pass off the body
	var body struct {
		Email    string
		Password string
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": responses.TryAgain,
		})
		return
	}

	user, err := database.FindUserByEmail(body.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message":   responses.UserNotFound,
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message":   responses.InvalidCredentials,
		})
		return
	}

	tokenString, err := tokens.GenerateNewAuthToken(user.Email, user.User_ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message":   responses.ErrorTokenCreation,
		})
		return
	}

	refreshTokenString, err := tokens.GenerateNewRefreshToken()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message":   responses.ErrorRefreshTokenCreation,
		})
		return
	}

	err = database.UpdateUserRefreshToken(user.Email, refreshTokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": responses.InternalServerError,
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", "Bearer "+tokenString, 3600*24, "", "", true, true)
	c.SetCookie("refresh_token", refreshTokenString, 3600*24*5, "", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": responses.LoginSuccessful,
	})
}

func Validate(context *gin.Context) {

	email, exists := context.Get("email")

	if !exists {
		context.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "I'm Logged In",
		"email":   email,
	})
}

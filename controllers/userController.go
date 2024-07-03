package controllers

import (
	"net/http"
	"study-bean/database"
	"study-bean/models"
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
			"error": "Some Error Occured",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"success": true,
		"users":   users,
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
			"error": "Failed to parse body",
		})
		return
	}

	// Check is user with same email exists
	user, _ := database.FindUserByEmail(body.Email)
	if user != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "User with this Email Already Exists",
		})
		return
	}

	// Check if user with same username exists
	user, _ = database.FindUserByUsername(body.Username)
	if user != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Username Taken",
		})
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	refreshTokenString, err := tokens.GenerateNewRefreshToken()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Failed to generate a refresh token",
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
			"error":   "Please Try Again",
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
			"error": "Failed to parse body",
		})
		return
	}

	user, err := database.FindUserByEmail(body.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "User Does Not Exist",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid Credentials",
		})
		return
	}

	tokenString, err := tokens.GenerateNewAuthToken(user.Email, user.User_ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Error Creating Token",
		})
		return
	}

	refreshTokenString, err := tokens.GenerateNewRefreshToken()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Error Creating Refresh Token",
		})
		return
	}

	err = database.UpdateUserRefreshToken(user.Email, refreshTokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Some Error Occured",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", "Bearer "+tokenString, 3600*24, "", "", true, true)
	c.SetCookie("refresh_token", refreshTokenString, 3600*24*5, "", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login Successful",
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

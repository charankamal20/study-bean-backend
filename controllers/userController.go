package controllers

import (
	"errors"
	"net/http"
	"os"
	"study-bean/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func getUserByName(email string) (*models.User, error) {
    for i, u := range models.Users {
        if u.Email == email {
            return &models.Users[i], nil
        }
    }

    return nil, errors.New("user not found")
}

func SignUp(context *gin.Context) {
	// Get the email/pass off the body
	var body struct {
		Email string
		Password string
		Username string
	}

	if err := context.Bind(&body); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse body",
		})
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Create the user
	user, _ := getUserByName(body.Email)
	if user != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": "User Already Exists",
		})
		return
	}

	var newUser models.User
	newUser.Password = string(hash)
	newUser.Email = body.Email
	newUser.Username = body.Username

	newUser.ID = uuid.New().String()

	models.Users = append(models.Users, newUser)

	// Respond
	context.JSON(http.StatusCreated, gin.H{
		"success": true,
		"user": newUser,
	})
}


func Login(context *gin.Context) {
// Get the email/pass off the body
	var body struct {
		Email string
		Password string
	}

	if err := context.Bind(&body); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse body",
		})
		return
	}

	user, err := getUserByName(body.Email)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": "User Does Not Exist",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": "Invalid Credentials",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": "Error Creating Token",
		})
		return
	}

	context.SetSameSite(http.SameSiteNoneMode)
	context.SetCookie("Authorization", tokenString, 3600 * 24, "", "", true, true)

	context.JSON(http.StatusOK, gin.H{
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
		"email": email,
	})
}
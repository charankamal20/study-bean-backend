package tokens

import (
	"errors"
	"os"
	"study-bean/database"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateNewRefreshToken() (string, error) {
	refresh_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Duration(time.Now().Day()) * 5).Unix(),
	})
	refreshTokenString, err := refresh_token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}

	return refreshTokenString, err
}

func GenerateNewAuthToken(email string, user_id string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":   email,
		"user_id": user_id,
		"exp":     time.Now().Local().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func UpdateAuthTokenFromRefreshToken(email string, oldRefreshToken string) (string, error) {

	// find the user from email
	user, err := database.FindUserByEmail(email)
	if err != nil {
		return "", err
	}

	// get user refresh token
	// check if they are same
	if oldRefreshToken != user.RefreshToken {
		return "", errors.New("refresh token not valid")
	}

	//!TODO: Check if token is expired, in that case, prompt to login

	// if yes then generate new token and return
	authToken, err := GenerateNewAuthToken(email, user.User_ID)
	if err != nil {
		return "", err
	}

	// err = database.UpdateUserByKey(user.ID, "refresh_token", refreshToken)
	// if err != nil {
	// 	return "", err
	// }

	return authToken, nil
}

package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"study-bean/responses"
	"study-bean/tokens"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(context *gin.Context) {
	// Get Cookie
	auth_token, err := context.Cookie("Authorization")
	if err != nil {
		context.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	tokenString := strings.Split(auth_token, " ")[1]
	// Decode/ validate
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_KEY")), nil
	})

	if err != nil {
		fmt.Println(err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		context.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	fmt.Println(claims)

	expirationTime, ok := claims["exp"].(float64)
	if !ok {
		context.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	expirationTimeInt := int64(expirationTime)

	currentTime := time.Now().Unix()

	if currentTime > expirationTimeInt {
		refresh_token, err := context.Cookie("refresh_token")
		fmt.Println("TOKEN_EXPIRED")
		if err != nil {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		auth_token, err = tokens.UpdateAuthTokenFromRefreshToken(claims["email"].(string), refresh_token)
		if err != nil {
			context.AbortWithStatus(http.StatusForbidden)
			return
		}

		context.SetSameSite(http.SameSiteLaxMode)
		context.SetCookie("Authorization", "Bearer "+auth_token, 3600*24, "", "", true, true)
	}

	context.Set("email", claims["email"])
	context.Set("user_id", claims["user_id"])

	context.Next()
}

func SessionMiddleware(context *gin.Context) {
	// Get Cookie
	session_cookie, err := context.Cookie("Session")
	if err != nil {
		context.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	tokenString := strings.Split(session_cookie, " ")[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_KEY")), nil
	})

	if err != nil {
		fmt.Println(err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		context.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	fmt.Println(claims)

	expirationTime, ok := claims["exp"].(float64)
	if !ok {
		context.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	expirationTimeInt := int64(expirationTime)

	currentTime := time.Now().Unix()

	if currentTime > expirationTimeInt {
		context.JSON(http.StatusNotFound, responses.SessionExpired)
	}

	context.Set("session_id", claims["session_id"])

	context.Next()
}

package middlware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)



func RequireAuth(context *gin.Context) {
	// Get Cookie
	tokenString, err := context.Cookie("Authorization")

	if err != nil {
		context.AbortWithStatus(http.StatusUnauthorized)
		return
	}

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

	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		context.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// !TODO: Check is user exists

	context.Set("email", claims["email"])

	context.Next()
}
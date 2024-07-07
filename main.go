package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"study-bean/controllers"
	"study-bean/initializers"
	"study-bean/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	println("hello@kanishk")
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.NewAuth()
}

func CORSMiddleware() gin.HandlerFunc {

	var origin string
	env := os.Getenv("ENVIRONMENT")

	if env == "PRODUCTION" {
		origin = "https://studybean.classikh.me"
	} else {
		origin = "http://localhost:3000"
	}

	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	defer func() {
		if err := initializers.MongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	router := gin.Default()
	router.Use(CORSMiddleware())
	// router.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:3000"},
	// 	AllowMethods:     []string{"PUT", "GET", "POST", "DELETE", "OPTIONS", "PATCH", "HEAD"},
	// 	AllowHeaders:     []string{"Origin", "Authorization", "x-csrf-token"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))

	// router.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{origins},
	// 	AllowHeaders:     []string{"Origin", "Authorization", "x-csrf-token"},
	// 	AllowMethods:     []string{"PUT", "GET", "POST", "DELETE", "OPTIONS", "PATCH", "HEAD"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	AllowOriginFunc: func(origin string) bool {
	// 		return origin == origins
	// 	},
	// 	MaxAge: 12 * time.Hour,
	// }))
	htmlFormat := `<html><body>%v</body></html>`
	router.GET("/", func(c *gin.Context) {
		html := fmt.Sprintf(htmlFormat, `<a href="/auth/google">Login through Google</a> <a href="/auth/github">Login through Github</a> <a href="/auth/discord">Login through Discord</a>`)
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	})

	//* OPEN ROUTES
    router.POST("/signup", controllers.SignUp)
    router.POST("/login", controllers.Login)

	router.GET("/auth/:provider/callback", controllers.GoogleAuthCallbackfunc)
	router.GET("/auth/:provider", controllers.OAuthProvider)
	router.GET("/logout/:provider", controllers.OAuthLogout)
	router.GET("/user", controllers.GetAllUsers)
	router.GET("/user/:userID", controllers.GetSingleUser)

	//* TEST ROUTES
	router.GET("/validate", middleware.RequireAuth, controllers.Validate)

	//* PRIVATE ROUTES
	router.POST("/todo", middleware.RequireAuth, controllers.AddTodo)
	router.GET("/todo", middleware.RequireAuth, controllers.GetAllTodos)
	router.PUT("/todo/:todo_id", middleware.RequireAuth, controllers.UpdateTodo)
	router.PUT("/toggleTodo/:todo_id", middleware.RequireAuth, controllers.ToggleTodoState)
	router.DELETE("/todo/:todo_id", middleware.RequireAuth, controllers.DeleteTodo)
	router.Run()
}

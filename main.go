package main

import (
	"context"
	"os"
	"study-bean/controllers"
	"study-bean/initializers"
	"study-bean/middleware"

	"github.com/gin-gonic/gin"
)


func init() {
    initializers.LoadEnvVariables()
	initializers.ConnectToDB()
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
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

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


    router.POST("/signup", controllers.SignUp)
    router.POST("/login", controllers.Login)
	router.GET("/validate", middleware.RequireAuth, controllers.Validate)
    router.GET("/user", controllers.GetAllUsers)

	router.POST("/todo", middleware.RequireAuth, controllers.AddTodo)

	router.Run()
}

package main

import (
	"net/http"
	"study-bean/controllers"
	"study-bean/initializers"
	"study-bean/middlware"
	"study-bean/models"

	"github.com/gin-gonic/gin"
)


func init() {
    initializers.LoadEnvVariables()
}

func addUser(c *gin.Context) {
	var newUser models.User

	if err := c.BindJSON(&newUser); err != nil {
		return
	}

	models.Users = append(models.Users, newUser)
	c.JSON(http.StatusCreated, newUser)
}

// getAlbums responds with the list of all albums as JSON.
func getUsers(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, models.Users)
}



// func getUser(context *gin.Context) {
//     id := context.Param("id")

//     user, err := getUserById(id)

//     if err != nil {
//         context.JSON(http.StatusNotFound, gin.H{"message": "User Not Found"})
//         return
//     }

//     context.JSON(http.StatusOK, user)
// }

// func updatePassword(context *gin.Context) {
//     id := context.Param("id")

//     user, err := getUserById(id)

//     if err != nil {
//         context.JSON(http.StatusNotFound, gin.H{"message": "User Not Found"})
//         return
//     }

//     var body struct {
//         Password string `json:"password" binding:"required"`
//     }

//     if err := context.ShouldBindJSON(&body); err != nil {
//         context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
//     }

//     user.Password = body.Password
//     context.JSON(http.StatusOK, user)

// }

func main() {
    router := gin.Default()

    router.POST("/signup", controllers.SignUp)
    router.POST("/login", controllers.Login)
    router.GET("/validate", middlware.RequireAuth,controllers.Validate)


    router.GET("/user", getUsers)
	router.POST("/user", addUser)
    // router.GET("/user/:id", getUser)
    // router.PATCH("/user/:id", updatePassword)
    router.Run()
}
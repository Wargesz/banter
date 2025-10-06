package main

import (
	"banter/controllers"
	"banter/initialisers"
	"banter/middleware"
	"os"

	"github.com/gin-gonic/gin"
)

func init() {
    initialisers.LoadEnvVars()
    initialisers.DbConnect()
    initialisers.SyncDatabase()
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.tmpl")
    r.Static("/static", "./static")
    routes(r)
	r.Run(":" + os.Getenv("PORT"))
}

func routes(r *gin.Engine) {
	r.GET("/", middleware.Auth, controllers.Posts)
	r.GET("/login", controllers.LoginPage)
	r.POST("/login", controllers.LoginUser)
    r.GET("/signup", controllers.SignUpPage)
    r.POST("/signup", controllers.SignUpUser)
}

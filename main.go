package main

import (
	"banter/controllers"
	"banter/initialisers"
	"banter/middleware"
	"fmt"
	"os"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	initialisers.LoadEnvVars()
	initialisers.DbConnect()
	initialisers.SyncDatabase()
}

func main() {
	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
	})
	r.LoadHTMLGlob("templates/*.tmpl")
	r.Static("/static", "./static")
	routes(r)
	r.Run(":" + os.Getenv("PORT"))
}

func routes(r *gin.Engine) {
	r.GET("/", middleware.Auth, controllers.Posts)
	r.POST("/", middleware.Auth, controllers.AddPost)
	r.GET("/login", controllers.LoginPage)
	r.POST("/login", controllers.LoginUser)
	r.GET("/signup", controllers.SignUpPage)
	r.POST("/signup", controllers.SignUpUser)
    r.DELETE("/remove/:id", middleware.Auth, controllers.RemovePost)
	r.POST("/edit", middleware.Auth, controllers.EditPost)
    r.GET("/profile", middleware.Auth, controllers.Profile)
    r.DELETE("/profile", middleware.Auth, controllers.RemoveAccount)
}

func formatAsDate(t time.Time) string {
	weekday := t.Weekday().String()
	year, month, day := t.Date()
	hour := t.Hour()
	minute := t.Minute()
	return fmt.Sprintf("%s, %d/%02d/%02d %02d:%02d", weekday, year, month, day, hour, minute)
}

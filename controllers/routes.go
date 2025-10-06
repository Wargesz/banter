package controllers

import (
	"banter/initialisers"
	"banter/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Posts(c *gin.Context) {
	c.HTML(http.StatusOK, "templates/posts.tmpl", gin.H{})
}

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "templates/login.tmpl", gin.H{})
}

func LoginUser(c *gin.Context) {
	username := c.PostForm("username")
	if username == "" {
		c.HTML(http.StatusBadRequest, "templates/login.tmpl", gin.H{
			"errorMessage": "no username provided",
		})
		return
	}
	password := c.PostForm("password")
	if password == "" {
		c.HTML(http.StatusBadRequest, "templates/login.tmpl", gin.H{
			"errorMessage": "no password provided",
		})
		return
	}
	var user models.User
	initialisers.DB.First(&user, "username = ?", username)
	if user.ID == 0 {
		c.HTML(http.StatusBadRequest, "templates/login.tmpl", gin.H{
			"errorMessage": "invalid username",
		})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		c.HTML(http.StatusBadRequest, "templates/login.tmpl", gin.H{
			"errorMessage": "failed to hash password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})
	tokenString, err := token.SignedString(os.Getenv("SECRET"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "templates/login.tmpl", gin.H{
			"errorMessage": "failed to create cookie",
		})
		return
	}
    c.SetSameSite(http.SameSiteLaxMode)
    c.SetCookie("Authorization", tokenString, 3600 * 24 * 7, "", "", false, true)

	c.Redirect(http.StatusFound, "/")
}

func SignUpPage(c *gin.Context) {
	c.HTML(http.StatusOK, "templates/signup.tmpl", gin.H{})
}

func SignUpUser(c *gin.Context) {
	username := c.PostForm("username")
	if username == "" {
		c.HTML(http.StatusBadRequest, "templates/signup.tmpl", gin.H{
			"errorMessage": "no username provided",
		})
		return
	}
	password := c.PostForm("password")
	if password == "" {
		c.HTML(http.StatusBadRequest, "templates/signup.tmpl", gin.H{
			"errorMessage": "no password provided",
		})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		c.HTML(http.StatusBadRequest, "templates/signup.tmpl", gin.H{
			"errorMessage": "failed to hash password",
		})
		return
	}
	user := models.User{Username: username, Password: string(hash)}
	result := initialisers.DB.Create(&user)
	if result.Error != nil {
		c.HTML(http.StatusBadRequest, "templates/signup.tmpl", gin.H{
			"errorMessage": "failed to create user",
        })
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})
	tokenString, err := token.SignedString(os.Getenv("SECRET"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "templates/signup.tmpl", gin.H{
			"errorMessage": "failed to create cookie",
		})
		return
	}
    c.SetSameSite(http.SameSiteLaxMode)
    c.SetCookie("Authorization", tokenString, 3600 * 24 * 7, "", "", false, true)
	c.Redirect(http.StatusFound, "/")
}

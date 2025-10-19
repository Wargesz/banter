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

type PostID struct {
	Id int `json:"id" binding:"required"`
}

type EditedPost struct {
	Id      int    `json:"id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func Posts(c *gin.Context) {
	var posts []models.Post
	initialisers.DB.Preload("User").Find(&posts)
	user, _ := c.Get("user")
	c.HTML(http.StatusOK, "templates/posts.tmpl", gin.H{
		"posts": posts,
		"user":  user,
	})
}

func AddPost(c *gin.Context) {
	user, _ := c.Get("user")
	title := c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "no title provided",
		})
		return
	}
	content := c.PostForm("content")
	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "no title provided",
		})
		return
	}
	var post models.Post
	post.Title = title
	post.Content = content
	post.UserID = user.(models.User).ID
	result := initialisers.DB.Create(&post)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "failed to create post",
		})
		return
	}
	c.Redirect(http.StatusFound, "/")
}

func EditPost(c *gin.Context) {
	var editedPost EditedPost
	err := c.ShouldBindJSON(&editedPost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": err.Error(),
		})
		return
	}
	result := initialisers.DB.Model(&models.Post{}).
		Where("id = ?", editedPost.Id).Update("content", editedPost.Content)
    if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": result.Error,
		})
		return
    }
	c.JSON(http.StatusOK, gin.H{
		"redirect": "/",
	})
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
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.HTML(http.StatusBadRequest, "templates/login.tmpl", gin.H{
			"errorMessage": "failed to create cookie",
		})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*7, "", "", false, true)

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
	user := models.User{Username: username, Password: string(hash),
		ProfilePicture: "static/picture.png"}
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
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.HTML(http.StatusBadRequest, "templates/signup.tmpl", gin.H{
			"errorMessage": "failed to create cookie",
		})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*7, "", "", false, true)
	c.Redirect(http.StatusFound, "/")
}

func RemovePost(c *gin.Context) {
	var postId PostID
	err := c.ShouldBindJSON(&postId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	result := initialisers.DB.Delete(&models.Post{}, postId.Id)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": result.Error,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"redirect": "/",
	})
}

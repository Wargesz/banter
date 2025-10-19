package controllers

import (
	"banter/initialisers"
	"banter/models"
	"net/http"

	"github.com/gin-gonic/gin"
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

func RemovePost(c *gin.Context) {
	user, _ := c.Get("user")
	postId := c.Param("id")
    var post models.Post
    result := initialisers.DB.Where("id == ?", postId).First(&post)
    if result.Error != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "reason": "no post",
        })
        return
    }
    if post.UserID != user.(models.User).ID {
        c.JSON(http.StatusUnauthorized, gin.H{
            "reason": "cannot delete other user posts",
        })
        return
    }
    result = initialisers.DB.Where("id == ?", postId).Delete(&models.Post{}) 
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

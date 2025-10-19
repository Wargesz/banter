package controllers

import (
	"banter/initialisers"
	"banter/models"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

var allowedExtensions = []string{"png", "jpeg", "jpg"}

func Profile(c *gin.Context) {
	user, _ := c.Get("user")
	//show own profile view
	queryUsername := c.Query("username")
	if queryUsername == "" || user.(models.User).Username == queryUsername {
		c.HTML(http.StatusOK, "templates/edit-profile.tmpl", gin.H{
			"user":      user,
			"postCount": getPostCount(user.(models.User).ID),
		})
		return
	}
	//show public profile view
	var queryUser models.User
	results := initialisers.DB.First(&queryUser, "username == ?", queryUsername)
	if results.Error != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}
	c.HTML(http.StatusOK, "templates/view-profile.tmpl", gin.H{
		"user":      queryUser,
		"postCount": getPostCount(queryUser.ID),
	})
}

func getPostCount(userId uint) int64 {
	var count int64
	initialisers.DB.Preload("User").Model(&models.Post{}).
		Where("user_id == ?", userId).Count(&count)
	return count
}

func RemoveAccount(c *gin.Context) {
	user, _ := c.Get("user")
	var u models.User
	initialisers.DB.First(&u, "id == ?", user.(models.User).ID)
	result := initialisers.DB.Select(clause.Associations).Delete(&u)
	if user.(models.User).ProfilePicture != "static/users/default.png" {
		err := os.Remove("./" + user.(models.User).ProfilePicture)
		if err != nil {
			println("could not remove profile picture")
		}
	}
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": "could not delete user",
		})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"reason": "delete success",
	})
}

func ChangePicture(c *gin.Context) {
	u, _ := c.Get("user")
	form, _ := c.MultipartForm()
	file := form.File["file"][0]
	fileType := strings.Split(file.Filename, ".")[1]
	if !slices.Contains(allowedExtensions, fileType) {
		c.HTML(http.StatusBadRequest, "templates/edit-profile.tmpl", gin.H{
			"user":          u.(models.User),
			"postCount":     getPostCount(u.(models.User).ID),
			"statusMessage": "forbidden filetype",
		})
		return
	}
	path := "static/users/user-" + u.(models.User).Username + "." + fileType
	err := c.SaveUploadedFile(file, path)
	if err != nil {
		c.HTML(http.StatusBadRequest, "templates/edit-profile.tmpl", gin.H{
			"user":          u.(models.User),
			"postCount":     getPostCount(u.(models.User).ID),
			"statusMessage": "upload failed",
		})
		return
	}
	var user models.User
	initialisers.DB.Where("id == ?", u.(models.User).ID).First(&user)
	user.ProfilePicture = path
	initialisers.DB.Save(&user)
	c.HTML(http.StatusBadRequest, "templates/edit-profile.tmpl", gin.H{
		"user":          user,
		"postCount":     getPostCount(user.ID),
		"statusMessage": "upload success",
	})
}

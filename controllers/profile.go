package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Profile(c *gin.Context) {
    c.HTML(http.StatusOK, "templates/profile.tmpl", gin.H{

    })
}

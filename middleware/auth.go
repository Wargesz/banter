package middleware

import (
	"banter/initialisers"
	"banter/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const SECRET = "myCoolTokenString"

func Auth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("SECRET")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.Redirect(http.StatusFound, "/login")
			return
		}
        var user models.User
        initialisers.DB.First(&user, "id == ?", claims["sub"])
        if user.ID == 0 {
			c.Redirect(http.StatusFound, "/login")
			return
        }
        c.Set("user", user)
        c.Next()
        return
	} else {
		c.Redirect(http.StatusFound, "/login")
		return
	}
}

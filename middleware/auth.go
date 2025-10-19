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

func Auth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
        abortWithRedirectToLogin(c)
        return
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("SECRET")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		abortWithRedirectToLogin(c)
        return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			abortWithRedirectToLogin(c)
            return
		}
		var user models.User
		result := initialisers.DB.First(&user, "username == ?", claims["sub"])
		if result.Error != nil {
			abortWithRedirectToLogin(c)
            return
		}
		if user.ID == 0 {
			abortWithRedirectToLogin(c)
            return
		}
		c.Set("user", user)
		return
	}
    abortWithRedirectToLogin(c)
}

func abortWithRedirectToLogin(c *gin.Context) {
    c.Redirect(http.StatusFound, "/login")
	c.Abort()
}

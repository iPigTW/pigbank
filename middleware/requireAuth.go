package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ipigtw/pigbank/initialize"
	"github.com/ipigtw/pigbank/models"
)
type claim struct {
	Sub string
	Exp int64
	jwt.Claims
}
func RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	
	token, _ := jwt.ParseWithClaims(tokenString, &claim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})
	claims := token.Claims.(*claim)
	if time.Now().Unix() > claims.Exp {
		c.AbortWithStatus(http.StatusUnauthorized)
		fmt.Println(claims.Exp)
	}
	var user *models.User
	initialize.DB.First(&user, claims.Sub)
	if user.ID == 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	c.Set("user", user)
	c.Next()
}
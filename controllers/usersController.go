package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ipigtw/pigbank/initialize"
	"github.com/ipigtw/pigbank/models"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	var body struct {
		Name string
		Email string
		Password string
	}
	if c.Bind(&body) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid body",
		})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Failed to generate hash",
		})
		return
	}
	user := models.User{
		Name: body.Name,
		Email: body.Email,
		Password: string(hash),
	}
	result := initialize.DB.Create(&user)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Failed to add user to database",
		})
		return
		
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})

}
func Login(c *gin.Context) {
	var body struct {
		Name string
		Email string
		Password string
	}
	if c.Bind(&body) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid body",
		})
		return
	}
	var user models.User
	initialize.DB.First(&user, "email = ?", body.Email)
	if user.ID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "No user found",
		})
		return
	}
	
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Password Incorrect",
		})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(os.Getenv("SECRET")))
	
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24, "", "", true, true)
	fmt.Println(body)
	fmt.Println(user)
	c.JSON(http.StatusOK, gin.H{
		"message":"Logged in!",
		
	})
}
func Validate(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"message": "Hi"+user.(models.User).Name,
	})
}
func Transfer(c *gin.Context) {
	user, _ := c.Get("user")
	var body struct {
		Target string
		Amount int
	}
	if c.Bind(&body) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid body",
		})
		return
	}
	if body.Amount <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Amount is too low",
		})
		return
	}
	if body.Amount > user.(*models.User).Balance {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Amount must be lower than your balance",
		})
		return
	}
	var target *models.User
	initialize.DB.First(&target, "name = ?", body.Target)
	if target.ID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Target not found",
		})
		return
	}
	target.Balance += body.Amount
	user.(*models.User).Balance -= body.Amount
	initialize.DB.Save(user)
	initialize.DB.Save(&target)
	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})
	
}
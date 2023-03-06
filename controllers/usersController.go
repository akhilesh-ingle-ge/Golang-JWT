package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/akhilesh-ingle-ge/jwt/config"
	"github.com/akhilesh-ingle-ge/jwt/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {

	// step-1:= get the email/password from the request body
	var body struct {
		Email    string
		Password string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read body",
		})
		return
	}

	// step-2:= hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to hash password",
		})
		return
	}

	// step-3:= create a user
	user := models.User{Email: body.Email, Password: string(hash)}
	result := config.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to create user",
		})
		return
	}

	// step-4:= respond
	c.JSON(http.StatusOK, gin.H{
		"Success": "Success",
	})
}

func Login(c *gin.Context) {
	// step-1:= get the email/password from the request body
	var body struct {
		Email    string
		Password string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read body",
		})
		return
	}

	// Step-2:= Look up requested user
	var user models.User

	config.DB.First(&user, "Email = ?", body.Email)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Invalid email or password",
		})
		return
	}

	// Step-3:= Compare sent in password with user hash password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Invalid email or password",
		})
		return
	}

	// step-4:= generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Minute * 5).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to create token",
		})
		return
	}

	// step-5:= send it back
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"Message": "Success",
	})
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"Message": user,
	})
}

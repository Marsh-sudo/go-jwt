package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/marsh-sudo/go-jwt/initializers"
	"github.com/marsh-sudo/go-jwt/models"
	"golang.org/x/crypto/bcrypt"
)



func SignUp(c *gin.Context){

	var body struct {
		Email string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"failed to create body",
	})
	return

	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"failed to hash password",
		})
		return
	}

	// create the user
	user := models.User{Email: body.Email, Password: string(hash)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"failed to create user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}



func Login(c *gin.Context) {
	// get email and passs
	var body struct {
		Email string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"failed to create body",
	})
	return

	}


	// look up user
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"invlid email and password",
		})
		return
	}

	//compare sent in pass
	err := bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(body.Password))

	if err != nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"invalid password or email",
		})
		return
	}

	//generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"subject": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"failed to create token",
		})
		return
	}
	
	

	// send it back
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization",tokenString,3600 * 24 * 30, "","",false,true)
	c.JSON(http.StatusOK,gin.H{})
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"message":user,
	})
}
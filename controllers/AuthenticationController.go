package controllers

import (
	"net/http"
	"rmalan/go/mig-assignment/auth"
	"rmalan/go/mig-assignment/config"
	"rmalan/go/mig-assignment/models"
	"strings"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterUser(context *gin.Context) {
	var user models.User

	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if err := user.HashPassword(user.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	record := config.Instance.Create(&user)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusCreated, gin.H{"user_id": user.ID, "email": user.Email, "username": user.Username})
}

func Login(context *gin.Context) {
	var request LoginRequest
	var user models.User

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// check if email exists and password is correct
	record := config.Instance.Where("username = ?", request.Username).First(&user)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}

	credentialError := user.CheckPassword(request.Password)
	if credentialError != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		context.Abort()
		return
	}

	tokenString, err := auth.GenerateJWT(user.ID, user.Username, user.Email)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	authentication := models.Authentication{Token: tokenString}
	recordAuth := config.Instance.Create(&authentication)
	if recordAuth.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": recordAuth.Error.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func Logout(context *gin.Context) {
	authHeader := context.GetHeader("Authorization")
	bearerToken := strings.Split(authHeader, " ")

	var authentication models.Authentication

	config.Instance.Where("token = ?", bearerToken[1]).Delete(&authentication)

	context.JSON(http.StatusOK, gin.H{"message": "success"})
}

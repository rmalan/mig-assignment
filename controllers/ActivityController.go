package controllers

import (
	"net/http"
	"rmalan/go/mig-assignment/config"
	"rmalan/go/mig-assignment/models"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type StoreActivityRequest struct {
	Date     string `json:"date" binding:"required"`
	Start    string `json:"start" binding:"required"`
	End      string `json:"end" binding:"required"`
	Activity string `json:"activity" binding:"required"`
}

type UpdateActivityRequest struct {
	Date     string `json:"date"`
	Start    string `json:"start"`
	End      string `json:"end"`
	Activity string `json:"activity"`
}

func ActivityByUser(context *gin.Context) {
	authHeader := context.GetHeader("Authorization")
	bearerToken := strings.Split(authHeader, " ")

	token, err := jwt.ParseWithClaims(
		bearerToken[1],
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)

	if err != nil {
		return
	}

	userId := token.Claims.(*JWTClaim).ID

	var activities []models.Activity
	if err := config.Instance.Where("user_id = ?", userId).Find(&activities).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	context.JSON(http.StatusOK, activities)
}

func ActivityByUserByDate(context *gin.Context) {
	authHeader := context.GetHeader("Authorization")
	bearerToken := strings.Split(authHeader, " ")

	token, err := jwt.ParseWithClaims(
		bearerToken[1],
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)

	if err != nil {
		return
	}

	userId := token.Claims.(*JWTClaim).ID

	var activities []models.Activity
	if err := config.Instance.Where("user_id = ?", userId).Where("date = ?", context.Param("date")).Find(&activities).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	context.JSON(http.StatusOK, activities)
}

func Store(context *gin.Context) {
	authHeader := context.GetHeader("Authorization")
	bearerToken := strings.Split(authHeader, " ")

	token, err := jwt.ParseWithClaims(
		bearerToken[1],
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)

	if err != nil {
		return
	}

	userId := token.Claims.(*JWTClaim).ID

	var request StoreActivityRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	var attendance models.Attendance
	if err := config.Instance.Where("user_id = ?", userId).Where("date = ?", request.Date).Not("check_in = ?", "").First(&attendance).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "user has not checked in!"})
		return
	}

	activity := models.Activity{UserId: userId, Date: request.Date, Start: request.Start, End: request.End, Activity: request.Activity}
	record := config.Instance.Create(&activity)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusCreated, gin.H{"date": request.Date, "start": request.Start, "end": request.End, "activity": request.Activity})
}

func Update(context *gin.Context) {
	authHeader := context.GetHeader("Authorization")
	bearerToken := strings.Split(authHeader, " ")

	token, err := jwt.ParseWithClaims(
		bearerToken[1],
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)

	if err != nil {
		return
	}

	userId := token.Claims.(*JWTClaim).ID

	var request UpdateActivityRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	var activity models.Activity
	if err := config.Instance.Where("id = ?", context.Param("id")).Where("user_id = ?", userId).First(&activity).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "record not found!"})
		return
	}

	var attendance models.Attendance
	if err := config.Instance.Where("user_id = ?", userId).Where("date = ?", request.Date).Not("check_in = ?", "").First(&attendance).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "user has not checked in!"})
		return
	}

	record := config.Instance.Model(&activity).Updates(models.Activity{UserId: userId, Date: request.Date, Start: request.Start, End: request.End, Activity: request.Activity})
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, request)
}

func Destroy(context *gin.Context) {
	var activity models.Activity
	if err := config.Instance.Where("id = ?", context.Param("id")).First(&activity).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "record not found!"})
		return
	}

	config.Instance.Delete(&activity)

	context.JSON(http.StatusOK, gin.H{"message": "success"})
}

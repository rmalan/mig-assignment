package controllers

import (
	"net/http"
	"rmalan/go/mig-assignment/config"
	"rmalan/go/mig-assignment/models"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("supersecretkey")

type JWTClaim struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

type CheckInRequest struct {
	Date    string `json:"date" binding:"required"`
	CheckIn string `json:"check_in" binding:"required"`
}

type CheckOutRequest struct {
	Date     string `json:"date" binding:"required"`
	CheckOut string `json:"check_out" binding:"required"`
}

func CheckIn(context *gin.Context) {
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

	var request CheckInRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	var attendance models.Attendance
	if err := config.Instance.Where("user_id = ?", userId).Where("date = ?", request.Date).First(&attendance).Error; err != nil {
		checkIn := models.Attendance{UserId: userId, Date: request.Date, CheckIn: request.CheckIn}
		record := config.Instance.Create(&checkIn)
		if record.Error != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
			context.Abort()
			return
		}

		context.JSON(http.StatusCreated, gin.H{"date": request.Date, "check_in": request.CheckIn})
		return
	}

	checkIn := models.Attendance{CheckIn: request.CheckIn}
	record := config.Instance.Model(&attendance).Updates(checkIn)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, gin.H{"date": request.Date, "check_in": request.CheckIn})
}

func CheckOut(context *gin.Context) {
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

	var request CheckOutRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	var attendance models.Attendance
	if err := config.Instance.Where("user_id = ?", userId).Where("date = ?", request.Date).First(&attendance).Error; err != nil {
		checkOut := models.Attendance{UserId: userId, Date: request.Date, CheckOut: request.CheckOut}
		record := config.Instance.Create(&checkOut)
		if record.Error != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
			context.Abort()
			return
		}

		context.JSON(http.StatusCreated, gin.H{"date": request.Date, "check_out": request.CheckOut})
		return
	}

	checkOut := models.Attendance{CheckOut: request.CheckOut}
	record := config.Instance.Model(&attendance).Updates(checkOut)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, gin.H{"date": request.Date, "check_out": request.CheckOut})
}

func AttendanceByUser(context *gin.Context) {
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

	var attendances []models.Attendance
	if err := config.Instance.Where("user_id = ?", userId).Find(&attendances).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	context.JSON(http.StatusOK, attendances)
}

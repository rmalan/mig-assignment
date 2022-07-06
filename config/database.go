package config

import (
	"log"
	"rmalan/go/mig-assignment/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Instance *gorm.DB
var dbError error

// func DbConnect(connectionString string) {
// 	Instance, dbError = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
func DbConnect() {
	Instance, dbError = gorm.Open(sqlite.Open("mig_assignment.sqlite"), &gorm.Config{})
	if dbError != nil {
		log.Fatal(dbError)
		panic("Cannot connect to DB")
	}

	log.Println("Connected to Database!")
}

func DbMigrate() {
	Instance.AutoMigrate(&models.Activity{})
	Instance.AutoMigrate(&models.Attendance{})
	Instance.AutoMigrate(&models.User{})

	log.Println("Database Migration Completed!")
}

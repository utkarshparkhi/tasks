package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"tasks/models"
	"tasks/views"
)

var db *gorm.DB
var err error

const (
	dsn = "root:root@tcp(127.0.0.1:3306)/tasks?charset=utf8mb4&parseTime=True&loc=Local"
)

func main() {
	router := gin.Default()
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&models.Task{})
	router.POST("/create", views.CreateTask)
	router.GET("/get/:id", views.GetTask)
	router.GET("/get", views.GetAll)
	router.POST("/update/", views.Update)
	router.GET("/audit/:id", views.Audit)
	router.DELETE("/delete/:id", views.Delete)
	router.Run("localhost:8080")
}

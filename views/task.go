package views

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"tasks/models"
	"time"
)

const (
	oneWeek = 7 * 24 * time.Hour
	dsn     = "root:root@tcp(127.0.0.1:3306)/tasks?charset=utf8mb4&parseTime=True&loc=Local"
)

func CreateTask(c *gin.Context) {
	var newTask *models.Task
	var intermediateTask models.APITask
	if err := c.BindJSON(&intermediateTask); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	newTask = intermediateTask.GetTask()
	if !models.IsValidStatus(newTask.Status) {
		newTask.Status = models.PENDING
	}
	if models.IsValidETA(newTask.ETA) {
		newTask.ETA = time.Now().Add(oneWeek)
	}
	newTask.UpdatedAt = time.Now()
	create(newTask, c)
}

func GetTask(c *gin.Context) {

	var key, _ = strconv.Atoi(c.Param("id"))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err == nil {
		var task models.APITask
		db.Model(models.Task{}).Select([]string{"Title", "Status", "ETA"}).Last(&task, key)
		c.IndentedJSON(http.StatusOK, &task)
	}
}
func GetAll(c *gin.Context) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err == nil {
		var tasks []models.APITask
		db.Model(models.Task{}).Select([]string{"ID", "Title", "Status", "ETA"}).Find(&tasks)
		c.IndentedJSON(http.StatusOK, &tasks)
	}
}
func Update(c *gin.Context) {
	var newTask *models.Task
	var intermediateTask models.APITask
	if err := c.BindJSON(&intermediateTask); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	newTask = intermediateTask.GetTask()
	if !models.IsValidStatus(newTask.Status) {
		newTask.Status = models.PENDING
	}
	newTask.UpdatedAt = time.Now()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err == nil {
		var count int64
		db.Model(models.Task{}).Where("id=?", newTask.ID).Count(&count)
		if count < 1 {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}
	if newTask.ID == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	create(newTask, c)
}
func Audit(c *gin.Context) {
	id := c.Param("id")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	var tasks []models.Task
	if err == nil {
		db.Model(models.Task{}).Find(&tasks, id)
		var diffs []models.Difference
		for i := 1; i < len(tasks); i++ {
			diffs = append(diffs, *models.GetDiff(&tasks[i-1], &tasks[i]))
		}
		c.IndentedJSON(http.StatusOK, diffs)
	}

}
func Delete(c *gin.Context) {
	id := c.Param("id")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err == nil {
		db.Delete(models.Task{}, id)
		c.Status(http.StatusOK)
	}
}
func create(task *models.Task, c *gin.Context) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err == nil {
		db.Create(task)
		c.IndentedJSON(http.StatusCreated, task)
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

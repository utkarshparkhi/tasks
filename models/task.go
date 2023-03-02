package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"reflect"
	"time"
)

type Status string

const (
	PENDING     Status = "Pending"
	IN_PROGRESS Status = "In Progress"
	IN_REVIEW   Status = "In Review"
	COMPLETE    Status = "Complete"
)

type Task struct {
	gorm.Model `json:",omitempty"`
	ID         uint      `gorm:"primaryKey" json:"ID,omitempty"`
	Title      string    `json:"title,omitempty"`
	Status     Status    `json:"status,omitempty" gorm:"type:enum('PENDING', 'IN_PROGRESS', 'IN_REVIEW', 'COMPLETE')"`
	ETA        time.Time `json:"ETA,omitempty"`
	UpdatedAt  time.Time `gorm:"primaryKey" gorm:"autoCreateTime" json:"updated_at,omitempty"`
}
type APITask struct {
	ID     uint      `json:"ID,omitempty"`
	Title  string    `json:"title"`
	Status Status    `json:"status"`
	ETA    time.Time `json:"ETA"`
}
type Difference struct {
	Task1         *Task `json:"oldTask"`
	Task2         *Task `json:"updatedTask"`
	UpdatedFields []string
}

func GetDiff(task1 *Task, task2 *Task) *Difference {
	diff := new(Difference)
	diff.Task1 = task1
	diff.Task2 = task2
	val1 := reflect.ValueOf(task1)
	val2 := reflect.ValueOf(task2)
	for i := 0; i < val1.Elem().NumField(); i++ {
		fmt.Println(val1.Elem().Type().Field(i).Name)
		fmt.Println(val1.Elem().Field(i), val2.Elem().Field(i))
		if val1.Elem().Field(i).Interface() != val2.Elem().Field(i).Interface() && val1.Elem().Type().Field(i).Name != "Model" {
			diff.UpdatedFields = append(diff.UpdatedFields, val1.Elem().Type().Field(i).Name)
		}
	}
	return diff
}

func IsValidStatus(status Status) bool {
	return status == PENDING || status == IN_PROGRESS || status == IN_REVIEW || status == COMPLETE

}

func IsValidETA(time2 time.Time) bool {
	return time2.Before(time.Now())
}

func (apiTask *APITask) GetTask() *Task {
	return &Task{
		ID:     apiTask.ID,
		Title:  apiTask.Title,
		Status: apiTask.Status,
		ETA:    apiTask.ETA,
	}
}

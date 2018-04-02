package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Model struct {
	ID        int `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time	`json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}

type Employee struct {
	Model
	Name     string `json:"name" gorm:"type:varchar(255);not null"`
	Company  string `json:"company" gorm:"type:varchar(255)"`
	Email    string `json:"email" gorm:"type:varchar(255);not null;unique"`
	Password string `json:"password" gorm:"type:varchar(255);not null"`
}

func FindEmployees(db *gorm.DB) (employees []Employee, err error) {
	err = db.Find(&employees).Error
	return
}

func (e *Employee) Create(db *gorm.DB) (err error) {
	err = db.Create(e).Error
	return
}

func (e *Employee) Find(db *gorm.DB) (err error) {
	err = db.First(e).Error
	return
}

func (e *Employee) Update(db *gorm.DB) (err error) {
	err = db.Model(e).Update(e).Error
	return
}

func (e *Employee) Delete(db *gorm.DB) (err error) {
	err = db.Delete(e).Error
	return
}
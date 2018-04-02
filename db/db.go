package db

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jinzhu/gorm"
	"echo-sample/models"
	"github.com/jinzhu/configor"
	"fmt"
)

var config = struct {
	DBName     string `default:"echo"`
	User     string `default:"postgres"`
	Host     string `default:"localhost"`
	Password string `default:"password" env:"DBPassword"`
	Port     string `default:"5433"`
}{}

func New() (db *gorm.DB) {

	configor.Load(&config)

	args := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		config.Host,
		config.Port,
		config.User,
		config.DBName,
		config.Password)

	db, err := gorm.Open("postgres", args)

	if err != nil {
		panic(err)
	}

	db.LogMode(true)

	autoMigrate(db)

	return
}

func autoMigrate(db *gorm.DB) {
	db.AutoMigrate(&models.Employee{})
}

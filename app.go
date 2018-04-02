package main

import (
	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
	"echo-sample/db"
	"github.com/jinzhu/configor"
)

type app struct {
	*echo.Echo
	db *gorm.DB
}

var config = struct {
	Host string `default:"localhost"`
	Port string `default:"1323"`
}{}

func main() {

	// Load configuration
	configor.Load(&config)

	app := &app{echo.New(), db.New()}

	app.Debug = true

	// Routing setup
	app.initRouter()

	defer app.db.Close()

	// Wrap db pointer into echo.context as a middleware
	app.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			context.Set("db", app.db)
			return next(context)
		}
	})

	// launch
	app.Logger.Fatal(app.Start(config.Host + ":" + config.Port))
}

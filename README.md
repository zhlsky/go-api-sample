# API sample by using [echo](https://github.com/labstack/echo) and [gorm](https://github.com/jinzhu/gorm)
## 1.Start
### 1.1 launch
```
# Go dependency management tool https://golang.github.io/dep/
dep ensure 


# launch DB container
docker-compose up

# Create DB
docker exec echo-db createdb -U postgres echo

# launch App
go run *.go
```
### 1.2 send a request
```
# use httpie
http POST localhost:1323/api/v1/employee name=zhl email=zhanghl@yahoo.co.jp company=echo password=password
```

```
# Response

HTTP/1.1 200 OK
Content-Length: 243
Content-Type: application/json; charset=UTF-8
Date: Thu, 29 Mar 2018 09:51:14 GMT

{
    "company": "echo",
    "created_at": "2018-03-29T18:51:14.984260399+09:00",
    "deleted_at": null,
    "email": "zhanghl@yahoo.co.jp",
    "id": 1,
    "name": "zhl",
    "password": "password",
    "updated_at": "2018-03-29T18:51:14.984260399+09:00"
}
```

## 2.See app.go

Use [echo](https://github.com/labstack/echo) and [gorm](https://github.com/jinzhu/gorm)

```go
// define a struct which has an echo pointer and a DB pointer
type app struct {
	*echo.Echo
	db *gorm.DB
}
```

##### Use [configor](https://github.com/jinzhu/configor) for configuration.[configor](https://github.com/jinzhu/configor) is Golang Configuration tool that support YAML, JSON, TOML, Shell Environment by writing struct tags.
```
# define host and port for server
var config = struct {
	Host string `default:"localhost"`
	Port string `default:"1323"`
}{}
```

##### Instead using global variable, wrap a DB pointer into context as a middleware of echo so that we can decoupling from top and it's easy to test.
``` go
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
```

## 3.See router
##### Make a route group and register routes for HTTP method.Routes can be registered by specifying HTTP method, path and a matching handler.
```go
func (app *app) initRouter() {
	v1 := app.Group("/api/v1")
	{
		v1.GET("/employee", echo.HandlerFunc(GetEmployees))
		v1.GET("/employee/:id", echo.HandlerFunc(GetEmployee))
		v1.POST("/employee", echo.HandlerFunc(CreateEmployee))
		v1.PATCH("/employee/:id", echo.HandlerFunc(UpdateEmployee))
		v1.DELETE("/employee/:id", echo.HandlerFunc(DeleteEmployee))
	}
}
```

## 4.See handler
##### echo wraps HTTP request and response into context. So you can get request parameter from context, put a HTTP status code and your data into JSON response.

```go
func GetEmployee(context echo.Context) error {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		echo.NewHTTPError(http.StatusBadRequest, "Employee Id must be int")
	}

	employee := &Employee{Model: Model{ID: id}}
	if err := employee.Find(context.Get("db").(*gorm.DB)); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Employee does not exist.")
	}

	return context.JSON(http.StatusOK, employee)
}
```

## 5. See model
##### Extract the common parts of all model struct such as id, created_at and so on.The struct tag of json means response JSON key.The struct tag of gorm looks like SQL and in fact it really works as DDL.
```go
type Model struct {
	ID        int `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time	`json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" grom:"index"`
}
```
##### Inherit common parts and write json and gorm struct tags to define your response and your data type. 
```go
type Employee struct {
	Model
	Name     string `json:"name" gorm:"type:varchar(255);not null"`
	Company  string `json:"company" gorm:"type:varchar(255)"`
	Email    string `json:"email" gorm:"type:varchar(255);not null;unique"`
	Password string `json:"password" gorm:"type:varchar(255);not null"`
}
```
##### Simple CURD. For Detail check [gorm docs](http://gorm.io/docs/).It also support advanced topics like Raw SQL, transaction, migration and so on.
```go
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
```

## 6. See db.go
##### The first import line means that we use PostgreSQL dialect and initialize it. Your can change Mysql or SQLite.
```go
import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jinzhu/gorm"
	"echo-sample/models"
	"github.com/jinzhu/configor"
	"fmt"
)
```
##### DB configuration by using [configor](https://github.com/jinzhu/configor).
```go
var config = struct {
	DBName     string `default:"echo"`
	User     string `default:"postgres"`
	Host     string `default:"localhost"`
	Password string `default:"password" env:"DBPassword"`
	Port     string `default:"5433"`
}{}
```
##### Load the DB configuration and create a DB connection. Shut down by throw a panic if DB connection failed. Then make set the DB logger on and auto migrate your data.
```go
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


```
#####  Just pass the model(DB table) to the AutoMigration function to ensure auto migration.
```go
func autoMigrate(db *gorm.DB) {
	db.AutoMigrate(&models.Employee{})
}
```

## 7. Add more in future...
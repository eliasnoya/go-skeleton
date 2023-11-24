package main

import (
	"database/sql"
	"fmt"
	"os"
	"tuples/infrastructure/app"
	"tuples/infrastructure/middlewares"
	"tuples/modules/tenant"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *sql.DB
var orm *gorm.DB

func main() {

	BootEnv()
	BootDbms()

	// Application container all the basic app ports shared between all requests
	app := &app.Container{Db: db, Orm: orm, Gin: gin.Default()}

	// Global middlewares
	middlewares.SetupGlobalMiddlewares(app)

	// TENANT ENDPOINTS
	tenant.SetupServices(app)

	app.Listen()
}

func BootEnv() {
	godotenv.Load()
}

func BootDbms() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE"))

	sqlConnection, sqlConnectionError := sql.Open("mysql", dsn)

	if sqlConnectionError != nil {
		panic(sqlConnectionError.Error())
	}

	db = sqlConnection

	var gormLogLevel logger.LogLevel = logger.Info

	if os.Getenv("ENV") == "PROD" {
		gormLogLevel = logger.Silent
	}

	ormConnection, ormConnectionError := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlConnection,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(gormLogLevel),
	})

	if ormConnectionError != nil {
		panic(ormConnectionError.Error())
	}

	orm = ormConnection
}

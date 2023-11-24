package main

import (
	"database/sql"
	"fmt"
	"os"
	"tuples/modules/tenant"

	"strings"
	"text/template"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/gertd/go-pluralize"
)

// instance different connection and diferrent enviorment for CLI
var db *sql.DB
var orm *gorm.DB

func main() {
	godotenv.Load()
	CliBootDbms()

	entrypoint := os.Args[1]

	switch entrypoint {
	case "make:model":
		makeModel(os.Args[2])
	case "migrate":
		migrate()
	}

	os.Exit(0)
}

func migrate() {
	orm.AutoMigrate(&tenant.Tenant{})
}

func makeModel(name string) {

	pluralize := pluralize.NewClient()

	singular := pluralize.Singular(name)
	plural := pluralize.Plural(name)

	PkgName := strings.ToLower(singular)
	ModelVar := PkgName[0]

	data := struct {
		PkgName   string
		ModelName string
		ModelVar  string
		TableName string
	}{
		PkgName:   PkgName,
		ModelName: strings.Title(singular),
		ModelVar:  string(ModelVar), // first letter lower
		TableName: strings.ToLower(plural),
	}

	tmpl, parseStubError := template.ParseFiles("./cli/stubs/model.tmpl")
	if parseStubError != nil {
		fmt.Println("Error parsing template:", parseStubError)
		return
	}

	createDirError := os.MkdirAll("./modules/"+PkgName, os.ModePerm)
	if createDirError != nil {
		fmt.Println("Error creating directory:", createDirError)
		return
	}

	// Create model file
	outputFile, createFileError := os.Create("./modules/" + PkgName + "/entity.go")
	if createFileError != nil {
		fmt.Println("Error creating output file:", outputFile)
		return
	}
	defer outputFile.Close()

	renderTemplateError := tmpl.Execute(outputFile, data)
	if renderTemplateError != nil {
		fmt.Println("Error executing template:", renderTemplateError)
	}

	fmt.Println("Model " + PkgName + " created...")
}

func CliBootDbms() {
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

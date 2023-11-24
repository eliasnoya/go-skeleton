package app

import (
	"database/sql"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Container struct {
	Db  *sql.DB
	Orm *gorm.DB
	Gin *gin.Engine
}

func (c *Container) Listen() {
	host := os.Getenv("APP_HOST")
	port := os.Getenv("APP_PORT")
	c.Gin.Run(host + ":" + port)
}

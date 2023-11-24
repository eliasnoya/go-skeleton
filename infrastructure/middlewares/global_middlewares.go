package middlewares

import (
	"net/http"
	"tuples/infrastructure/app"
	"tuples/modules/tenant"

	helmet "github.com/danielkov/gin-helmet"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type GlobalMiddlewares struct {
	App *app.Container
}

func SetupGlobalMiddlewares(container *app.Container) {
	gm := &GlobalMiddlewares{App: container}
	// cors setup
	gm.UseCors()
	// compress responses
	container.Gin.Use(gzip.Gzip(gzip.BestSpeed))
	// basic security headers
	container.Gin.Use(helmet.Default())
	// bussiness logic middlewares start here
	container.Gin.Use(gm.DetectTenant)
	// recover for panics
	container.Gin.Use(gin.Recovery())
}

// Setup cors, separo func por legibilidad
func (am *GlobalMiddlewares) UseCors() {
	config := cors.DefaultConfig()
	config.AllowWildcard = true
	config.AllowOrigins = []string{
		"http://*.tuples.local",
		"http://*.tuples.local:3000",
		"http://test.tuples.local:3000/",
	}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Cache-Control", "X-Requested-With", "X-Tenant"}
	config.AllowCredentials = true
	am.App.Gin.Use(cors.New(config))
}

// Detect tenant (X-Tenant Header) or fail -> all requests
func (am *GlobalMiddlewares) DetectTenant(c *gin.Context) {
	tenantDomain := c.GetHeader("X-Tenant")

	var tenant tenant.Tenant
	result := am.App.Orm.First(&tenant, "domain = ?", tenantDomain)

	// tenant not found? Shutdown request
	if result.Error != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": true, "message": "Tenant not detected"})
		c.Abort()
		return
	}

	// Load tenant to gin context
	c.Set("tenant", tenant)
	c.Next()
}

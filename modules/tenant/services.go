package tenant

import (
	"net/http"
	"tuples/helpers"
	"tuples/infrastructure/app"
	"tuples/lib"

	"github.com/gin-gonic/gin"
)

type TenantServices struct {
	App *app.Container
}

func SetupServices(container *app.Container) {
	ts := &TenantServices{App: container}

	tenantGroup := container.Gin.Group("/tenant")
	tenantGroup.GET("/setup", ts.GetSettings)
	tenantGroup.GET("/setup/step", ts.GetStep)
	tenantGroup.POST("/setup/step", ts.PostStep)

	tenantGroup.PATCH("/setup/theme", ts.PatchThemeSettings)
}

func (s *TenantServices) GetSettings(c *gin.Context) {

	tenant := c.Keys["tenant"].(Tenant) // from DetectTenant middleware

	lightColors := helpers.MergeStructs(DefaultColorScheme("light"), NewColorSchemeFromJson(tenant.LightSetup))
	darkColors := helpers.MergeStructs(DefaultColorScheme("dark"), NewColorSchemeFromJson(tenant.DarkSetup))

	response := Settings{
		Theme: Theme{
			Logo:    tenant.Logo,
			Default: tenant.DefaultTheme,
			Light:   lightColors.(ColorScheme),
			Dark:    darkColors.(ColorScheme),
		},
	}

	c.JSON(200, response)
}

func (s *TenantServices) GetStep(c *gin.Context) {
	tenant := c.Keys["tenant"].(Tenant)
	c.JSON(200, gin.H{"error": false, "step": tenant.SetupStep})
}

func (s *TenantServices) PostStep(c *gin.Context) {

	var request SetupStepRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.Step > 0 {
		tenant := c.Keys["tenant"].(Tenant) // from DetectTenant middleware
		s.App.Orm.Model(tenant).Update("setup_step", request.Step)

		c.JSON(200, gin.H{"error": false, "message": "Step updated"})
	}

	c.JSON(200, gin.H{"error": true, "message": "No step before 0"})
}

func (s *TenantServices) PatchThemeSettings(c *gin.Context) {

	var setupRequest SetupThemeRequest

	if valid := lib.IsValidRequest(c, &setupRequest); !valid {
		return
	}

	tenant := c.Keys["tenant"].(Tenant)
	tenant.Logo = setupRequest.Logo
	tenant.DefaultTheme = setupRequest.DefaultTheme
	tenant.LightSetup = setupRequest.Light.AsJson()
	tenant.DarkSetup = setupRequest.Dark.AsJson()
	tenant.SetupStep = setupRequest.NextStep

	s.App.Orm.Save(&tenant)

	c.JSON(201, gin.H{"error": false, "message": "Theme Saved"})
}

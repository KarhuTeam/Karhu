package front

import (
	"github.com/gin-gonic/gin"
	"github.com/karhuteam/karhu/models"
	"github.com/karhuteam/karhu/web"
	"net/http"
)

type ConfigurationController struct {
}

func NewConfigurationController(s *web.Server) *ConfigurationController {

	ctl := &ConfigurationController{}

	// 1 - Add a configuration
	s.GET("/application/configuration/:app_id", ctl.getAddConfigurationAction)
	s.POST("/application/configuration/:app_id", ctl.postAddConfigurationAction)
	s.GET("/application/configuration/:app_id/:config_id", ctl.getEditConfigurationAction)
	s.POST("/application/configuration/:app_id/:config_id", ctl.postEditConfigurationAction)

	return ctl
}

func (ctl *ConfigurationController) getApplication(c *gin.Context, id string) *models.Application {

	application, err := models.ApplicationMapper.FetchOne(id)

	// Error 500
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error_500.html", map[string]interface{}{
			"error": err,
		})
		return nil
	}

	// Error 404
	if application == nil {
		c.HTML(http.StatusNotFound, "error_404.html", map[string]interface{}{
			"title": "Application not found",
			"text":  "Application not found... It's not my fault",
		})
		return nil
	}

	return application
}

func (ctl *ConfigurationController) getConfig(c *gin.Context, app *models.Application, id string) *models.Config {

	config, err := models.ConfigMapper.FetchOne(app, id)

	// Error 500
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error_500.html", map[string]interface{}{
			"error": err,
		})
		return nil
	}

	// Error 404
	if config == nil {
		c.HTML(http.StatusNotFound, "error_404.html", map[string]interface{}{
			"title": "Config not found",
			"text":  "Config not found... It's not my fault",
		})
		return nil
	}

	return config
}

/**
 * 1 - add a configuration
 */
func (ctl *ConfigurationController) postAddConfigurationAction(c *gin.Context) {

	appId := c.Param("app_id")

	application := ctl.getApplication(c, appId)
	if application == nil {
		return
	}

	var form models.ConfigCreateForm
	if err := c.Bind(&form); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Workaround for bool field, because http form only send string
	// The value is "on" for true, and "" for false
	form.Enabled = c.PostForm("enabled") == "on"

	if err := form.Validate(); err != nil {
		c.HTML(http.StatusOK, "configuration_add.html", map[string]interface{}{
			"errors":      err.Errors,
			"application": application,
			"form":        form,
		})
		return
	}

	config := models.ConfigMapper.Create(application, &form)

	if err := models.ConfigMapper.Save(config); err != nil {
		panic(err)
	}

	c.Redirect(http.StatusFound, "/application/show/"+application.Name)
}

func (ctl *ConfigurationController) getAddConfigurationAction(c *gin.Context) {

	appId := c.Param("app_id")

	application := ctl.getApplication(c, appId)
	if application == nil {
		return
	}

	c.HTML(http.StatusOK, "configuration_add.html", map[string]interface{}{
		"application": application,
	})
}

func (ctl *ConfigurationController) getEditConfigurationAction(c *gin.Context) {

	appId := c.Param("app_id")
	configId := c.Param("config_id")

	application := ctl.getApplication(c, appId)
	if application == nil {
		return
	}

	config := ctl.getConfig(c, application, configId)
	if config == nil {
		return
	}

	// Hydrate the form
	var form models.ConfigUpdateForm
	form.Hydrate(config)

	c.HTML(http.StatusOK, "configuration_edit.html", map[string]interface{}{
		"application": application,
		"form":        form,
	})
}

func (ctl *ConfigurationController) postEditConfigurationAction(c *gin.Context) {

	appId := c.Param("app_id")
	configId := c.Param("config_id")

	application := ctl.getApplication(c, appId)
	if application == nil {
		return
	}

	config := ctl.getConfig(c, application, configId)
	if config == nil {
		return
	}

	var form models.ConfigUpdateForm
	if err := c.Bind(&form); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Workaround for bool field, because http form only send string
	// The value is "on" for true, and "" for false
	form.Enabled = c.PostForm("enabled") == "on"

	if err := form.Validate(); err != nil {
		c.HTML(http.StatusOK, "configuration_edit.html", map[string]interface{}{
			"errors":      err.Errors,
			"application": application,
			"form":        form,
		})
		return
	}

	config.Update(&form)

	if err := models.ConfigMapper.Update(config); err != nil {
		panic(err)
	}

	c.Redirect(http.StatusFound, "/application/show/"+application.Name)
}

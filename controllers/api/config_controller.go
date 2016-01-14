package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gotoolz/errors"
	"github.com/karhuteam/karhu/models"
	"net/http"
)

type ConfigController struct {
}

func NewConfigController(s *gin.RouterGroup) *ConfigController {

	ctl := &ConfigController{}

	s.POST("/apps/:id/configs", ctl.postConfig)
	s.PUT("/apps/:id/configs/:config_id", ctl.putConfig)
	s.GET("/apps/:id/configs", ctl.getConfigList)
	s.GET("/apps/:id/configs/:config_id", ctl.getConfig)

	return ctl
}

func (cc *ConfigController) getApp(c *gin.Context) (*models.Application, error) {

	id := c.Param("id")

	app, err := models.ApplicationMapper.FetchOne(id)
	if err != nil {
		panic(err)
	}

	if app == nil {
		return nil, errors.New(errors.Error{
			Label: "invalid_application",
			Field: "id",
			Text:  "Invalid application ID in URL",
		})
	}

	return app, nil
}

func (cc *ConfigController) postConfig(c *gin.Context) {

	app, err := cc.getApp(c)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	var form models.ConfigCreateForm
	if err := c.Bind(&form); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := form.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	conf := models.ConfigMapper.Create(app, &form)

	if err := models.ConfigMapper.Save(conf); err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, conf)
}

func (cc *ConfigController) putConfig(c *gin.Context) {

	app, err := cc.getApp(c)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	configId := c.Param("config_id")

	config, err := models.ConfigMapper.FetchOne(app, configId)
	if err != nil {
		panic(err)
	}

	if config == nil {
		c.JSON(http.StatusNotFound, config)
		return
	}

	var form models.ConfigUpdateForm
	if err := c.Bind(&form); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := form.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	config.Update(&form)

	if err := models.ConfigMapper.Update(config); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, config)
}

func (cc *ConfigController) getConfigList(c *gin.Context) {

	app, err := cc.getApp(c)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	configs, err := models.ConfigMapper.FetchAll(app)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, configs)
}

func (cc *ConfigController) getConfig(c *gin.Context) {

	app, err := cc.getApp(c)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	configId := c.Param("config_id")

	config, err := models.ConfigMapper.FetchOne(app, configId)
	if err != nil {
		panic(err)
	}

	if config == nil {
		c.JSON(http.StatusNotFound, config)
		return
	}

	c.JSON(http.StatusOK, config)
}

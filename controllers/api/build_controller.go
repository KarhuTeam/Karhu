package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gotoolz/errors"
	"github.com/karhuteam/karhu/models"
	"net/http"
)

type BuildController struct {
}

func NewBuildController(s *gin.RouterGroup) *BuildController {

	ctl := &BuildController{}

	s.POST("/apps/:id/env/:env/builds", ctl.postBuild)
	s.GET("/apps/:id/env/:env/builds", ctl.getBuildList)

	return ctl
}

func (pc *BuildController) getAppEnv(c *gin.Context) (*models.Application, *models.Environment, error) {

	id := c.Param("id")

	app, err := models.ApplicationMapper.FetchOne(id)
	if err != nil {
		panic(err)
	}

	if app == nil {
		return nil, nil, errors.New(errors.Error{
			Label: "invalid_application",
			Field: "id",
			Text:  "Invalid application ID in URL",
		})
	}

	envId := c.Param("env")

	env, err := models.EnvironmentMapper.FetchOne(app, envId)
	if err != nil {
		panic(err)
	}

	if env == nil {
		return nil, nil, errors.New(errors.Error{
			Label: "invalid_environment",
			Field: "id",
			Text:  "Invalid environment ID in URL",
		})
	}

	return app, env, nil
}

func (pc *BuildController) postBuild(c *gin.Context) {

	var form models.BuildCreateForm
	if err := c.Bind(&form); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := form.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	_, env, err := pc.getAppEnv(c)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	build := models.BuildMapper.Create(env, &form)

	if err := models.BuildMapper.Save(build); err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, build)
}

func (pc *BuildController) getBuildList(c *gin.Context) {

	_, env, err := pc.getAppEnv(c)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	builds, err := models.BuildMapper.FetchAll(env)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, builds)
}

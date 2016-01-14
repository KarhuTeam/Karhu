package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gotoolz/errors"
	"github.com/karhuteam/karhu/models"
	"net/http"
)

type DeploymentController struct {
}

func NewDeploymentController(s *gin.RouterGroup) *DeploymentController {

	ctl := &DeploymentController{}

	s.GET("/apps/:id/deploy/:deploy_id", ctl.getDeployment)
	s.GET("/apps/:id/deploy", ctl.getDeploymentList)

	return ctl
}

func (dc *DeploymentController) getApp(c *gin.Context) (*models.Application, error) {

	id := c.Param("id")

	app, err := models.ApplicationMapper.FetchOne(id)
	if err != nil {
		panic(err)
	}

	if app == nil {
		return nil, errors.New(errors.Error{
			Label: "invalid_application",
			Field: "app_id",
			Text:  "Invalid application ID in URL",
		})
	}

	return app, nil
}

func (dc *DeploymentController) getDeployment(c *gin.Context) {

	app, err := dc.getApp(c)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
	}

	deployId := c.Param("deploy_id")

	depl, err := models.DeploymentMapper.FetchOne(app, deployId)
	if err != nil {
		panic(err)
	}

	if depl == nil {
		c.JSON(http.StatusNotFound, depl)
		return
	}

	c.JSON(http.StatusOK, depl)
}

func (dc *DeploymentController) getDeploymentList(c *gin.Context) {

	app, err := dc.getApp(c)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	deploys, err := models.DeploymentMapper.FetchAll(app)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, deploys)
}

package api

import (
	"github.com/gin-gonic/gin"
	"github.com/karhuteam/karhu/models"
	"net/http"
)

type ApplicationController struct {
}

func NewApplicationController(s *gin.RouterGroup) *ApplicationController {

	ctl := &ApplicationController{}

	s.POST("/apps", ctl.postApplication)
	s.GET("/apps", ctl.getApplicationList)
	s.GET("/apps/:id", ctl.getApplication)
	s.PUT("/apps/:id", ctl.putApplication)
	s.DELETE("/apps/:id", ctl.deleteApplication)

	return ctl
}

func (pc *ApplicationController) postApplication(c *gin.Context) {

	var form models.ApplicationCreateForm
	if err := c.Bind(&form); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := form.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	application := models.ApplicationMapper.Create(&form)

	if err := models.ApplicationMapper.Save(application); err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, application)
}

func (pc *ApplicationController) getApplicationList(c *gin.Context) {

	applications, err := models.ApplicationMapper.FetchAll()
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, applications)
}

func (pc *ApplicationController) getApplication(c *gin.Context) {

	id := c.Param("id")

	application, err := models.ApplicationMapper.FetchOne(id)
	if err != nil {
		panic(err)
	}

	if application == nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, application)
}

func (pc *ApplicationController) putApplication(c *gin.Context) {

	id := c.Param("id")

	application, err := models.ApplicationMapper.FetchOne(id)
	if err != nil {
		panic(err)
	}

	if application == nil {
		c.JSON(http.StatusNotFound, application)
		return
	}

	var form models.ApplicationUpdateForm
	if err := c.Bind(&form); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := form.Validate(application); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	application.Update(&form)

	if err := models.ApplicationMapper.Update(application); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, application)
}

func (pc *ApplicationController) deleteApplication(c *gin.Context) {

	id := c.Param("id")

	application, err := models.ApplicationMapper.FetchOne(id)
	if err != nil {
		panic(err)
	}

	if err := models.ApplicationMapper.Delete(application); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, nil)
}

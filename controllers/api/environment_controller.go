package api

import (
	"github.com/gin-gonic/gin"
	"github.com/karhuteam/karhu/models"
	"github.com/gotoolz/errors"
	"log"
	"net/http"
)

type EnviromnentController struct {
}

func NewEnvironmentController(s *gin.RouterGroup) *EnviromnentController {

	ctl := &EnviromnentController{}

	s.POST("/apps/:id/env", ctl.postEnvironment)
	// s.GET("/apps/:id/env", ctl.getApplicationList)
	// s.GET("/apps/:id/env/:env", ctl.getApplication)
	// s.PUT("/apps/:id/env/:env", ctl.putApplication)
	// s.DELETE("/apps/:id/env/:env", ctl.deleteApplication)

	return ctl
}

func (ec *EnviromnentController) getApp(c *gin.Context) (*models.Application, error) {

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

func (ec *EnviromnentController) postEnvironment(c *gin.Context) {

	var form models.EnvironmentCreateForm
	if err := c.Bind(&form); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := form.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	app, err := ec.getApp(c)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	env := models.EnvironmentMapper.Create(app, &form)

	log.Println("env:", env)

	if err := models.EnvironmentMapper.Save(env); err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, env)
}

// func (pc *EnviromnentController) getApplicationList(c *gin.Context) {
//
// 	applications, err := models.ApplicationMapper.FetchAll()
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	c.JSON(http.StatusOK, applications)
// }
//
// func (pc *EnviromnentController) getApplication(c *gin.Context) {
//
// 	id := c.Param("id")
//
// 	application, err := models.ApplicationMapper.FetchOne(id)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	if application == nil {
// 		c.JSON(http.StatusNotFound, nil)
// 		return
// 	}
//
// 	c.JSON(http.StatusOK, application)
// }
//
// func (pc *EnviromnentController) putApplication(c *gin.Context) {
//
// 	var form models.ApplicationUpdateForm
// 	if err := c.Bind(&form); err != nil {
// 		c.AbortWithStatus(http.StatusBadRequest)
// 		return
// 	}
//
// 	if err := form.Validate(); err != nil {
// 		c.JSON(http.StatusBadRequest, err)
// 		return
// 	}
//
// 	id := c.Param("id")
//
// 	application, err := models.ApplicationMapper.FetchOne(id)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	if application == nil {
// 		c.JSON(http.StatusNotFound, application)
// 		return
// 	}
//
// 	application.Update(&form)
//
// 	if err := models.ApplicationMapper.Update(application); err != nil {
// 		panic(err)
// 	}
//
// 	c.JSON(http.StatusOK, application)
// }
//
// func (pc *EnviromnentController) deleteApplication(c *gin.Context) {
//
// 	id := c.Param("id")
//
// 	application, err := models.ApplicationMapper.FetchOne(id)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	if err := models.ApplicationMapper.Delete(application); err != nil {
// 		panic(err)
// 	}
//
// 	c.JSON(http.StatusOK, nil)
// }

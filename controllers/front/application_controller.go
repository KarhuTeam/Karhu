package front

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gotoolz/errors"
	"github.com/karhuteam/karhu/models"
	"github.com/karhuteam/karhu/web"
	"net/http"
	"strconv"
)

type ApplicationController struct {
}

func NewApplicationController(s *web.Server) *ApplicationController {

	ctl := &ApplicationController{}

	// 1 - Applications list
	s.GET("/", ctl.getApplicationsAction)
	// 2 - Show an application
	s.GET("/application/show/:id", ctl.getApplicationAction)
	// 3 - Add an application
	s.GET("/application/add", ctl.getAddApplicationAction)
	s.POST("/application/add", ctl.postAddApplicationAction)
	// 4 - Edit an application
	s.GET("/application/edit/:id", ctl.getEditApplicationAction)
	s.POST("/application/edit/:id", ctl.postEditApplicationAction)
	// 5 - Delete an application
	s.GET("/application/delete/:id", ctl.getDeleteApplicationAction)

	return ctl
}

func (ctl *ApplicationController) getApplication(c *gin.Context, id string) *models.Application {

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

/**
 * 1 - Applications list
 */
func (ctl *ApplicationController) getApplicationsAction(c *gin.Context) {

	tag := c.DefaultQuery("tag", "")

	applications, err := models.ApplicationMapper.FetchAllByTag(tag)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error_500.html", map[string]interface{}{
			"error": err,
		})
		return
	}

	tags, err := models.ApplicationMapper.FetchAllTags()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error_500.html", map[string]interface{}{
			"error": err,
		})
		return
	}

	c.HTML(http.StatusOK, "application_list.html", map[string]interface{}{
		"tags":         tags,
		"applications": applications,
	})
}

/**
 * 2 - Show an application
 */
func (ctl *ApplicationController) getApplicationAction(c *gin.Context) {

	id := c.Param("id")

	// Get the application
	application := ctl.getApplication(c, id)
	if application == nil {
		return
	}

	builds, err := models.BuildMapper.FetchAll(application)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error_500.html", map[string]interface{}{
			"error": err,
		})
		return
	}

	configs, err := models.ConfigMapper.FetchAll(application)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error_500.html", map[string]interface{}{
			"error": err,
		})
		return
	}

	logfiles, err := models.LogfileMapper.FetchAll(application)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error_500.html", map[string]interface{}{
			"error": err,
		})
		return
	}

	if application.Type == models.APPLICATION_TYPE_SERVICE {

		c.HTML(http.StatusOK, "service_show.html", map[string]interface{}{
			"application": application,
			"configs":     configs,
			"builds":      builds,
			"logfiles":    logfiles,
		})

	} else {

		deployments, err := models.DeploymentMapper.FetchAll(application)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error_500.html", map[string]interface{}{
				"error": err,
			})
			return
		}

		// Limit array size
		count, _ := strconv.Atoi(c.DefaultQuery("count", "10"))
		if count <= 0 {
			count = 10
		}
		if len(builds) > count {
			builds = builds[:count]
		}
		if len(deployments) > count {
			deployments = deployments[:count]
		}

		c.HTML(http.StatusOK, "application_show.html", map[string]interface{}{
			"application": application,
			"builds":      builds,
			"deployments": deployments,
			"configs":     configs,
			"logfiles":    logfiles,
		})
	}
}

/**
 * 3 - Add an application
 */
func (ctl *ApplicationController) getAddApplicationAction(c *gin.Context) {

	switch c.DefaultQuery("type", "application") {
	case models.APPLICATION_TYPE_SERVICE:
		c.HTML(http.StatusOK, "service_add.html", nil)
	// case models.APPLICATION_TYPE_APP:
	default:
		c.HTML(http.StatusOK, "application_add.html", nil)
	}
}

func (ctl *ApplicationController) postAddApplicationAction(c *gin.Context) {

	var form models.ApplicationCreateForm
	if err := c.Bind(&form); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := form.Validate(); err != nil {
		fmt.Println(err.Errors)
		c.HTML(http.StatusOK, "application_add.html", map[string]interface{}{
			"errors": err.Errors,
			"form":   form,
		})
		return
	}

	application := models.ApplicationMapper.Create(&form)

	if form.Type == models.APPLICATION_TYPE_SERVICE {

		if len(form.Packages) == 0 {
			c.HTML(http.StatusOK, "service_add.html", map[string]interface{}{
				"errors": errors.New(errors.Error{
					Label: "invalid_packages",
					Field: "packages",
					Text:  "Invalid packages count, min 1",
				}).Errors,
				"form": form,
			})
			return
		}

		build := models.BuildMapper.CreateService(application, form.Packages)

		if err := models.BuildMapper.Save(build); err != nil {
			panic(err)
		}

	}

	if err := models.ApplicationMapper.Save(application); err != nil {
		panic(err)
	}

	c.Redirect(http.StatusFound, "/")
}

/**
 * 4 - Edit an application
 */

func (ctl *ApplicationController) getEditApplicationAction(c *gin.Context) {

	id := c.Param("id")

	// Get the application
	application := ctl.getApplication(c, id)
	if application == nil {
		return
	}

	// Hydrate the form
	var form models.ApplicationUpdateForm
	form.Hydrate(application)

	switch c.DefaultQuery("type", "application") {
	case models.APPLICATION_TYPE_SERVICE:
		c.HTML(http.StatusOK, "service_edit.html", map[string]interface{}{
			"form": form,
		})
	// case models.APPLICATION_TYPE_APP:
	default:
		c.HTML(http.StatusOK, "application_edit.html", map[string]interface{}{
			"form": form,
		})
	}
}

func (ctl *ApplicationController) postEditApplicationAction(c *gin.Context) {

	id := c.Param("id")

	// Get the application
	application := ctl.getApplication(c, id)
	if application == nil {
		return
	}

	// Hydrate the form with the request values
	var form models.ApplicationUpdateForm
	if err := c.Bind(&form); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Check the form values
	if err := form.Validate(application); err != nil {
		fmt.Println(err.Errors)
		c.HTML(http.StatusOK, "application_edit.html", map[string]interface{}{
			"errors": err.Errors,
			"form":   form,
		})
		return
	}

	// Update the application
	application.Update(&form)

	if application.Type == models.APPLICATION_TYPE_SERVICE {

		if len(form.Packages) == 0 {
			c.HTML(http.StatusOK, "service_edit.html", map[string]interface{}{
				"errors": errors.New(errors.Error{
					Label: "invalid_packages",
					Field: "packages",
					Text:  "Invalid packages count, min 1",
				}).Errors,
				"form": form,
			})
			return
		}

		build, err := models.BuildMapper.FetchLast(application)
		if err != nil {
			panic(err)
		}

		if build == nil {
			panic("No last build for service: " + application.Name)
		}

		build.RuntimeCfg.Dependencies = form.Packages

		if err := models.BuildMapper.Update(build); err != nil {
			panic(err)
		}

	}

	// Save the application
	if err := models.ApplicationMapper.Update(application); err != nil {
		panic(err)
	}

	c.Redirect(http.StatusFound, "/application/show/"+application.Id.Hex())
}

/**
 * 5 - Delete an application
 */

func (ctl *ApplicationController) getDeleteApplicationAction(c *gin.Context) {

	id := c.Param("id")

	// Get the application
	application := ctl.getApplication(c, id)
	if application == nil {
		return
	}

	// Delete the application
	models.ApplicationMapper.Delete(application)

	c.Redirect(http.StatusFound, "/")
}

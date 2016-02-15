package front

import (
	"github.com/gin-gonic/gin"
	"github.com/karhuteam/karhu/models"
	"github.com/karhuteam/karhu/web"
	"net/http"
)

type LogfileController struct {
}

func NewLogfileController(s *web.Server) *LogfileController {

	ctl := &LogfileController{}

	// 1 - Add a logfile
	s.GET("/application/logfile/:app_id", ctl.getAddLogfileAction)
	s.POST("/application/logfile/:app_id", ctl.postAddLogfileAction)
	s.GET("/application/logfile/:app_id/:logfile_id", ctl.getEditLogfileAction)
	s.POST("/application/logfile/:app_id/:logfile_id", ctl.postEditLogfileAction)

	return ctl
}

func (ctl *LogfileController) getApplication(c *gin.Context, id string) *models.Application {

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

func (ctl *LogfileController) getLogfile(c *gin.Context, app *models.Application, id string) *models.Logfile {

	logfile, err := models.LogfileMapper.FetchOne(app, id)

	// Error 500
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error_500.html", map[string]interface{}{
			"error": err,
		})
		return nil
	}

	// Error 404
	if logfile == nil {
		c.HTML(http.StatusNotFound, "error_404.html", map[string]interface{}{
			"title": "Logfile not found",
			"text":  "Logfile not found... It's not my fault",
		})
		return nil
	}

	return logfile
}

/**
 * 1 - add a logfile
 */
func (ctl *LogfileController) postAddLogfileAction(c *gin.Context) {

	appId := c.Param("app_id")

	application := ctl.getApplication(c, appId)
	if application == nil {
		return
	}

	var form models.LogfileCreateForm
	if err := c.Bind(&form); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Workaround for bool field, because http form only send string
	// The value is "on" for true, and "" for false
	form.Enabled = c.PostForm("enabled") == "on"

	if err := form.Validate(); err != nil {
		c.HTML(http.StatusOK, "logfile_add.html", map[string]interface{}{
			"errors":      err.Errors,
			"application": application,
			"form":        form,
		})
		return
	}

	logfile := models.LogfileMapper.Create(application, &form)

	if err := models.LogfileMapper.Save(logfile); err != nil {
		panic(err)
	}

	c.Redirect(http.StatusFound, "/application/show/"+application.Name)
}

func (ctl *LogfileController) getAddLogfileAction(c *gin.Context) {

	appId := c.Param("app_id")

	application := ctl.getApplication(c, appId)
	if application == nil {
		return
	}

	var form models.ConfigCreateForm
	form.Notify = "restart:" + application.Name
	form.Enabled = true

	c.HTML(http.StatusOK, "logfile_add.html", map[string]interface{}{
		"application": application,
		"form":        form,
	})
}

func (ctl *LogfileController) getEditLogfileAction(c *gin.Context) {

	appId := c.Param("app_id")
	logfileId := c.Param("logfile_id")

	application := ctl.getApplication(c, appId)
	if application == nil {
		return
	}

	logfile := ctl.getLogfile(c, application, logfileId)
	if logfile == nil {
		return
	}

	// Hydrate the form
	var form models.LogfileUpdateForm
	form.Hydrate(logfile)

	c.HTML(http.StatusOK, "logfile_edit.html", map[string]interface{}{
		"application": application,
		"form":        form,
	})
}

func (ctl *LogfileController) postEditLogfileAction(c *gin.Context) {

	appId := c.Param("app_id")
	logfileId := c.Param("logfile_id")

	application := ctl.getApplication(c, appId)
	if application == nil {
		return
	}

	logfile := ctl.getLogfile(c, application, logfileId)
	if logfile == nil {
		return
	}

	var form models.LogfileUpdateForm
	if err := c.Bind(&form); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Workaround for bool field, because http form only send string
	// The value is "on" for true, and "" for false
	form.Enabled = c.PostForm("enabled") == "on"

	if err := form.Validate(); err != nil {
		c.HTML(http.StatusOK, "logfile_edit.html", map[string]interface{}{
			"errors":      err.Errors,
			"application": application,
			"form":        form,
		})
		return
	}

	logfile.Update(&form)

	if err := models.LogfileMapper.Update(logfile); err != nil {
		panic(err)
	}

	c.Redirect(http.StatusFound, "/application/show/"+application.Name)
}

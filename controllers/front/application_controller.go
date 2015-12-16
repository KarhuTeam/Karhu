package front

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/karhuteam/karhu/models"
	"github.com/karhuteam/karhu/web"
	"net/http"
)

type ApplicationController struct {
}

func NewApplicationController(s *web.Server) *ApplicationController {

	ctl := &ApplicationController{}

	s.GET("/", ctl.getApplicationsAction)
	s.GET("/application/show/:id", ctl.getApplicationAction)
	s.GET("/application/add", ctl.getAddApplicationAction)
	s.POST("/application/add", ctl.postAddApplicationAction)

	return ctl
}

/**
 * Add Application
 */
func (ctl *ApplicationController) getAddApplicationAction(c *gin.Context) {

	c.HTML(http.StatusOK, "add.html", nil)
}

func (ctl *ApplicationController) postAddApplicationAction(c *gin.Context) {

	var form models.ApplicationCreateForm
	if err := c.Bind(&form); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := form.Validate(); err != nil {
		fmt.Println(err.Errors)
		c.HTML(http.StatusOK, "add.html", map[string]interface{}{
			"errors": err.Errors,
			"form":   form,
		})
		return
	}

	application := models.ApplicationMapper.Create(&form)

	if err := models.ApplicationMapper.Save(application); err != nil {
		panic(err)
	}

	c.Redirect(http.StatusMovedPermanently, "/")
}

/**
 * Get2 Applications
 */
func (ctl *ApplicationController) getApplicationsAction(c *gin.Context) {

	applications, _ := models.ApplicationMapper.FetchAll()

	c.HTML(http.StatusOK, "list.html", map[string]interface{}{
		"applications": applications,
	})
}

/**
 * Get Application
 */
func (ctl *ApplicationController) getApplicationAction(c *gin.Context) {

	id := c.Param("id")

	application, err := models.ApplicationMapper.FetchOne(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", map[string]interface{}{
			"error": err,
		})
		return
	}

	if application == nil {
		c.HTML(http.StatusNotFound, "404.html", map[string]interface{}{
			"text": "Application not found",
		})
		return
	}

	c.HTML(http.StatusOK, "show.html", map[string]interface{}{
		"application": application,
	})
}

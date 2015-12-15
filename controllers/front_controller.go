package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/karhuteam/karhu/models"
	"github.com/karhuteam/karhu/web"
	"net/http"
)

type FrontController struct {
}

func NewFrontController(s *web.Server) *FrontController {

	ctl := &FrontController{}

	s.GET("/", ctl.getApplicationsAction)
	s.GET("/application/add", ctl.getAddApplicationAction)
	s.POST("/application/add", ctl.postAddApplicationAction)

	return ctl
}

func (ctl *FrontController) postAddApplicationAction(c *gin.Context) {

	var form models.ApplicationCreateForm
	if err := c.Bind(&form); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := form.Validate(); err != nil {
		// TODO
		//c.JSON(http.StatusBadRequest, err)
		fmt.Println(form)
		c.AbortWithError(21321312, err)
		return
	}

	application := models.ApplicationMapper.Create(&form)

	if err := models.ApplicationMapper.Save(application); err != nil {
		panic(err)
	}

	c.Redirect(http.StatusMovedPermanently, "/")
}

func (ctl *FrontController) getAddApplicationAction(c *gin.Context) {

	c.HTML(http.StatusOK, "add.html", nil)
}

func (ctl *FrontController) getApplicationsAction(c *gin.Context) {

	applications, _ := models.ApplicationMapper.FetchAll()

	c.HTML(http.StatusOK, "list.html", map[string]interface{}{
		"applications": applications,
	})
}

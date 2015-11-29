package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/karhuteam/karhu/web"
	"net/http"
)

type HomeController struct {
}

func NewHomeController(s *web.Server) *HomeController {

	ctl := &HomeController{}

	s.GET("/", ctl.getHome)

	return ctl
}

func (ctl *HomeController) getHome(c *gin.Context) {

	c.HTML(http.StatusOK, "home.html", nil)
}

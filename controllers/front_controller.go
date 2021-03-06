package controllers

import (
	"github.com/karhuteam/karhu/controllers/front"
	"github.com/karhuteam/karhu/web"
)

type FrontController struct {
}

func NewFrontController(s *web.Server) *FrontController {

	ctl := &FrontController{}

	front.NewMonitoringController(s)
	front.NewApplicationController(s)
	front.NewNodeController(s)
	front.NewDeploymentController(s)
	front.NewConfigurationController(s)
	front.NewLogController(s)
	front.NewLogfileController(s)
	front.NewAlertController(s)

	return ctl
}

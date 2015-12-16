package controllers

import (
	"github.com/karhuteam/karhu/controllers/front"
	"github.com/karhuteam/karhu/web"
)

type FrontController struct {
}

func NewFrontController(s *web.Server) *FrontController {

	ctl := &FrontController{}

	front.NewApplicationController(s)

	return ctl
}

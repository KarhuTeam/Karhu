package controllers

import (
	"github.com/karhuteam/karhu/controllers/api"
	"github.com/karhuteam/karhu/web"
)

type APIController struct {
}

func NewAPIController(s *web.Server) *APIController {

	ctl := &APIController{}

	apiGroup := s.Group("/api")

	api.NewApplicationController(apiGroup)
	api.NewEnvironmentController(apiGroup)
	api.NewBuildController(apiGroup)
	api.NewScriptController(apiGroup)

	return ctl
}

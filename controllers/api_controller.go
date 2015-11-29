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

	api.NewProjectController(apiGroup)

	return ctl
}

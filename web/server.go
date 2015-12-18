package web

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	*gin.Engine
}

func NewServer() *Server {

	s := &Server{gin.Default()}

	s.LoadHTMLGlob("views/**")

	s.Static("/static", "static")

	return s
}

func (s *Server) Run(bind string) {

	s.Engine.Run(bind)
}

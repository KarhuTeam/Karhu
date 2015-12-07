package main

import (
	"github.com/gotoolz/env"
	"github.com/karhuteam/karhu/controllers"
	_ "github.com/karhuteam/karhu/models" // For mgo
	"github.com/karhuteam/karhu/web"
)

func main() {

	s := web.NewServer()
	controllers.NewHomeController(s)
	controllers.NewAPIController(s)

	s.Run(env.GetDefault("LISTEN_ADDR", ":8080")) // listen and serve
}

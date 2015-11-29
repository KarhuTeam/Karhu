package main

import (
	"github.com/karhuteam/karhu/controllers"
	_ "github.com/karhuteam/karhu/models" // For mgo
	"github.com/karhuteam/karhu/web"
	"github.com/wayt/goenv"
)

func main() {

	s := web.NewServer()
	controllers.NewHomeController(s)
	controllers.NewAPIController(s)

	s.Run(goenv.GetDefault("LISTEN_ADDR", ":8080")) // listen and serve
}

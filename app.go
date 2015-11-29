package main

import (
	"flag"
	// "fmt"
	"github.com/karhuteam/karhu/controllers"
	"github.com/karhuteam/karhu/web"
)

var bind = flag.String("bind", ":8080", "Kahru bind address")

func main() {

	flag.Parse()

	s := web.NewServer()

	controllers.NewHomeController(s)

	s.Run(*bind) // listen and serve
}

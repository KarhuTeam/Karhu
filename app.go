package main

import (
	"github.com/gotoolz/env"
	"github.com/karhuteam/karhu/alerts"
	"github.com/karhuteam/karhu/controllers"
	_ "github.com/karhuteam/karhu/models"         // For mgo
	_ "github.com/karhuteam/karhu/ressources/ssh" // For ssh key-pair
	"github.com/karhuteam/karhu/web"
	"math/rand"
	"time"
)

func main() {

	rand.Seed(time.Now().UnixNano())

	s := web.NewServer()
	controllers.NewFrontController(s)
	controllers.NewAPIController(s)

	go alerts.Run()

	s.Run(env.GetDefault("LISTEN_ADDR", ":8080")) // listen and serve
}

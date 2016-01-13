package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type deploymentMapper struct{}

var DeploymentMapper = &deploymentMapper{}

const deploymentCollection = "deployment"

func init() {
	col := C(deploymentCollection)
	defer col.Database.Session.Close()

	// App Index
	col.EnsureIndex(mgo.Index{
		Key:        []string{"application_id"},
		Unique:     false,
		DropDups:   false,
		Background: true, // See notes.
		Sparse:     true,
	})

	col.EnsureIndex(mgo.Index{
		Key:        []string{"build_id"},
		Unique:     false,
		DropDups:   false,
		Background: true, // See notes.
		Sparse:     true,
	})
}

type Deployment struct {
	Id            bson.ObjectId `json:"id" bson:"_id"`
	ApplicationId bson.ObjectId `json:"-" bson:"application_id"`
	Application   *Application  `json:"application" bson:"-"`
	BuildId       bson.ObjectId `json:"-" bson:"build_id"`
	Build         *Build        `json:"build" bson:"-"`
	CreatedAt     time.Time     `json:"created_at" bson:"created_at"`
}

type Deployments []*Deployment

func (dm *deploymentMapper) Create(app *Application, build *Build) *Deployment {

	return &Deployment{
		Id:            bson.NewObjectId(),
		ApplicationId: app.Id,
		Application:   app,
		BuildId:       build.Id,
		Build:         build,
		CreatedAt:     time.Now(),
	}
}

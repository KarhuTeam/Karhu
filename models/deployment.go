package models

import (
	"github.com/gotoolz/errors"
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

type DeploymentStatus string

const (
	STATUS_NEW     DeploymentStatus = "new"
	STATUS_RUNNING                  = "running"
	STATUS_DONE                     = "done"
	STATUS_ERROR                    = "error"
)

type Deployment struct {
	Id            bson.ObjectId    `json:"id" bson:"_id"`
	ApplicationId bson.ObjectId    `json:"-" bson:"application_id"`
	Application   *Application     `json:"-" bson:"-"`
	BuildId       bson.ObjectId    `json:"-" bson:"build_id"`
	Build         *Build           `json:"-" bson:"-"`
	TmpPath       string           `json:"-" bson:"-"`
	CreatedAt     time.Time        `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at" bson:"updated_at"`
	Status        DeploymentStatus `json:"status" bson:"status"`
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
		UpdatedAt:     time.Now(),
		Status:        STATUS_NEW,
	}
}

func (dm *deploymentMapper) FetchOne(app *Application, deployId string) (*Deployment, error) {

	if !bson.IsObjectIdHex(deployId) {

		return nil, errors.New(errors.Error{
			Label: "invalid_deploy_id",
			Field: "deploy_id",
			Text:  "Invalid deploy id hex",
		})
	}

	col := C(deploymentCollection)
	defer col.Database.Session.Close()

	deploy := new(Deployment)
	if err := col.Find(bson.M{"application_id": app.Id, "_id": bson.ObjectIdHex(deployId)}).One(deploy); err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	return deploy, nil
}

func (dm *deploymentMapper) FetchAll(app *Application) (Deployments, error) {

	col := C(deploymentCollection)
	defer col.Database.Session.Close()

	var deploys Deployments
	if err := col.Find(bson.M{"application_id": app.Id}).Sort("-created_at").All(&deploys); err != nil {
		return nil, err
	}

	return deploys, nil
}

func (dm *deploymentMapper) Save(d *Deployment) error {

	col := C(deploymentCollection)
	defer col.Database.Session.Close()

	return col.Insert(d)
}

func (dm *deploymentMapper) Update(d *Deployment) error {

	col := C(deploymentCollection)
	defer col.Database.Session.Close()

	d.UpdatedAt = time.Now()

	return col.UpdateId(d.Id, d)
}

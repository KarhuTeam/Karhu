package models

import (
	"github.com/gotoolz/validator"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type buildMapper struct{}

var BuildMapper = &buildMapper{}

const buildCollection = "build"

func init() {
	col := C(buildCollection)
	defer col.Database.Session.Close()

	// App Index
	col.EnsureIndex(mgo.Index{
		Key:        []string{"application_id"},
		Unique:     false,
		DropDups:   false,
		Background: true, // See notes.
		Sparse:     true,
	})

	// Env Index
	col.EnsureIndex(mgo.Index{
		Key:        []string{"environment_id"},
		Unique:     false,
		DropDups:   false,
		Background: true, // See notes.
		Sparse:     true,
	})
}

type Build struct {
	Id            bson.ObjectId          `json:"id" bson:"_id"`
	ApplicationId bson.ObjectId          `json:"application_id" bson:"application_id"`
	EnvironmentId bson.ObjectId          `json:"environment_id" bson:"environment_id"`
	Version       string                 `json:"version" bson:"version"`
	CommitHash    string                 `json:"commit_hash" bson:"commit_id"`
	CommitUrl     string                 `json:"commit_url" bson:"commit_url"`
	Tags          []string               `json:"tags" bson:"tags"`
	Vars          map[string]interface{} `json:"vars" bson:"vars"`
	CreatedAt     time.Time              `json:"created_at" bson:"created_at"`
}

type Builds []*Build

// Build creation form
type BuildCreateForm struct {
	Version    string                 `json:"version" valid:"ascii,required"`
	CommitHash string                 `json:"commit_hash" valid:"hexadecimal,required"`
	CommitUrl  string                 `json:"commit_url" valid:"url,required"`
	Tags       []string               `json:"tags" valid:"-"`
	Vars       map[string]interface{} `json:"vars" valid:"-"`
}

// Validator for build creation
func (f BuildCreateForm) Validate() error {
	return validator.Validate(&f)
}

func (bm *buildMapper) Create(e *Environment, f *BuildCreateForm) *Build {

	return &Build{
		Id:            bson.NewObjectId(),
		EnvironmentId: e.Id,
		Version:       f.Version,
		CommitHash:    f.CommitHash,
		CommitUrl:     f.CommitUrl,
		Tags:          f.Tags,
		Vars:          f.Vars,
		CreatedAt:     time.Now(),
	}
}

func (bm *buildMapper) Save(b *Build) error {

	col := C(buildCollection)
	defer col.Database.Session.Close()

	return col.Insert(b)
}

// func (bm *buildMapper) Update(b *Build) error {
//
// 	col := C(buildCollection)
// 	defer col.Database.Session.Close()
//
// 	b.UpdatedAt = time.Now()
//
// 	return col.UpdateId(b.Id, b)
// }

func (bm *buildMapper) Delete(b *Build) error {

	col := C(buildCollection)
	defer col.Database.Session.Close()

	return col.RemoveId(b.Id)
}

func (bm *buildMapper) FetchAll(e *Environment) (Builds, error) {

	col := C(buildCollection)
	defer col.Database.Session.Close()

	var builds Builds
	if err := col.Find(bson.M{"environment_id": e.Id}).Sort("-created_at").All(&builds); err != nil {
		return nil, err
	}

	return builds, nil
}

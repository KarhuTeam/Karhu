package models

import (
	"github.com/wayt/govalidator"
	// "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type buildMapper struct{}

var BuildMapper = &buildMapper{}

const buildCollection = "build"

type Build struct {
	Id         bson.ObjectId          `json:"id" bson:"_id"`
	ProjectId  bson.ObjectId          `json:"project_id" bson:"project_id"`
	Version    string                 `json:"version" bson:"version"`
	CommitHash string                 `json:"commit_hash" bson:"commit_id"`
	CommitUrl  string                 `json:"commit_url" bson:"commit_url"`
	Tags       []string               `json:"tags" bson:"tags"`
	Vars       map[string]interface{} `json:"vars" bson:"vars"`
	CreatedAt  time.Time              `json:"created_at" bson:"created_at"`
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
	return govalidator.Validate(&f)
}

func (bm *buildMapper) Create(p *Project, f *BuildCreateForm) *Build {

	return &Build{
		Id:         bson.NewObjectId(),
		ProjectId:  p.Id,
		Version:    f.Version,
		CommitHash: f.CommitHash,
		CommitUrl:  f.CommitUrl,
		Tags:       f.Tags,
		Vars:       f.Vars,
		CreatedAt:  time.Now(),
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

func (bm *buildMapper) FetchAll(p *Project) (Builds, error) {

	col := C(buildCollection)
	defer col.Database.Session.Close()

	var builds Builds
	if err := col.Find(bson.M{"project_id": p.Id}).Sort("-created_at").All(&builds); err != nil {
		return nil, err
	}

	return builds, nil
}

package models

import (
	"github.com/gotoolz/validator"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type environmentMapper struct{}

var EnvironmentMapper = &environmentMapper{}

const environmentCollection = "environment"

func init() {
	col := C(environmentCollection)
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
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   false,
		Background: true, // See notes.
		Sparse:     true,
	})
}

type Environment struct {
	Id             bson.ObjectId          `json:"id" bson:"_id"`
	ApplicationId  bson.ObjectId          `json:"application_id" bson:"application_id"`
	Name           string                 `json:"name" bson:"name"`
	Description    string                 `json:"description" bson:"description"`
	CurrentBuildId bson.ObjectId          `json:"current_build_id,omitempty" bson:"current_build_id,omitempty"`
	Tags           []string               `json:"tags" bson:"tags"` // Tags are used for environment search
	Vars           map[string]interface{} `json:"vars" bson:"vars"` // Vars are set in env when deploying a environment
	CreatedAt      time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at" bson:"updated_at"`
}

func (p *Environment) Update(f *EnvironmentUpdateForm) {

	p.Name = f.Name
	p.Description = f.Description
	p.Tags = f.Tags
	p.Vars = f.Vars
}

type Environments []*Environment

// Environment creation form
type EnvironmentCreateForm struct {
	Name        string                 `json:"name" valid:"ascii,required"`
	Description string                 `json:"description" valid:"ascii"`
	Tags        []string               `json:"tags" valid:"-"`
	Vars        map[string]interface{} `json:"vars" valid:"-"`
}

// Validator for environment creation
func (f EnvironmentCreateForm) Validate() error {
	return validator.Validate(&f)
}

// Environment update form
type EnvironmentUpdateForm struct {
	Name        string                 `json:"name" valid:"ascii,required"`
	Description string                 `json:"description" valid:"ascii"`
	Tags        []string               `json:"tags" valid:"-"`
	Vars        map[string]interface{} `json:"vars" valid:"-"`
}

// Validator for environment update
func (f EnvironmentUpdateForm) Validate() error {
	return validator.Validate(&f)
}

func (pm *environmentMapper) Create(a *Application, f *EnvironmentCreateForm) *Environment {

	return &Environment{
		Id:            bson.NewObjectId(),
		ApplicationId: a.Id,
		Name:          f.Name,
		Description:   f.Description,
		Tags:          f.Tags,
		Vars:          f.Vars,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func (pm *environmentMapper) Save(p *Environment) error {

	col := C(environmentCollection)
	defer col.Database.Session.Close()

	return col.Insert(p)
}

func (pm *environmentMapper) Update(p *Environment) error {

	col := C(environmentCollection)
	defer col.Database.Session.Close()

	p.UpdatedAt = time.Now()

	return col.UpdateId(p.Id, p)
}

func (pm *environmentMapper) Delete(p *Environment) error {

	col := C(environmentCollection)
	defer col.Database.Session.Close()

	return col.RemoveId(p.Id)
}

func (pm *environmentMapper) FetchAll(a *Application) (Environments, error) {

	col := C(environmentCollection)
	defer col.Database.Session.Close()

	var environments Environments
	if err := col.Find(bson.M{"application_id": a.Id}).All(&environments); err != nil {
		return nil, err
	}

	return environments, nil
}

func (pm *environmentMapper) FetchOne(a *Application, id string) (*Environment, error) {

	col := C(environmentCollection)
	defer col.Database.Session.Close()

	environment := new(Environment)
	if err := col.Find(bson.M{"_id": bson.ObjectIdHex(id), "application_id": a.Id}).One(environment); err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	return environment, nil
}

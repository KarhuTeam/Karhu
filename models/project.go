package models

import (
	"github.com/wayt/govalidator"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type projectMapper struct{}

var ProjectMapper = &projectMapper{}

const projectCollection = "project"

type Project struct {
	Id             bson.ObjectId          `json:"id" bson:"_id"`
	Name           string                 `json:"name" bson:"name"`
	Description    string                 `json:"description" bson:"description"`
	CurrentBuildId bson.ObjectId          `json:"current_build_id" bson:"current_build_id"`
	Tags           []string               `json:"tags" bson:"tags"` // Tags are used for project search
	Vars           map[string]interface{} `json:"vars" bson:"vars"` // Vars are set in env when deploying a project
	CreatedAt      time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at" bson:"updated_at"`
}

func (p *Project) Update(f *ProjectUpdateForm) {

	p.Name = f.Name
	p.Description = f.Description
	p.Tags = f.Tags
	p.Vars = f.Vars
}

type Projects []*Project

// Project creation form
type ProjectCreateForm struct {
	Name        string                 `json:"name" valid:"ascii,required"`
	Description string                 `json:"description" valid:"ascii"`
	Tags        []string               `json:"tags" valid:"-"`
	Vars        map[string]interface{} `json:"vars" valid:"-"`
}

// Validator for project creation
func (f ProjectCreateForm) Validate() error {
	return govalidator.Validate(&f)
}

// Project update form
type ProjectUpdateForm struct {
	Name        string                 `json:"name" valid:"ascii,required"`
	Description string                 `json:"description" valid:"ascii"`
	Tags        []string               `json:"tags" valid:"-"`
	Vars        map[string]interface{} `json:"vars" valid:"-"`
}

// Validator for project update
func (f ProjectUpdateForm) Validate() error {
	return govalidator.Validate(&f)
}

func (pm *projectMapper) Create(f *ProjectCreateForm) *Project {

	return &Project{
		Id:          bson.NewObjectId(),
		Name:        f.Name,
		Description: f.Description,
		Tags:        f.Tags,
		Vars:        f.Vars,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (pm *projectMapper) Save(p *Project) error {

	col := C(projectCollection)
	defer col.Database.Session.Close()

	return col.Insert(p)
}

func (pm *projectMapper) Update(p *Project) error {

	col := C(projectCollection)
	defer col.Database.Session.Close()

	p.UpdatedAt = time.Now()

	return col.UpdateId(p.Id, p)
}

func (pm *projectMapper) Delete(p *Project) error {

	col := C(projectCollection)
	defer col.Database.Session.Close()

	return col.RemoveId(p.Id)
}

func (pm *projectMapper) FetchAll() (Projects, error) {

	col := C(projectCollection)
	defer col.Database.Session.Close()

	var projects Projects
	if err := col.Find(nil).All(&projects); err != nil {
		return nil, err
	}

	return projects, nil
}

func (pm *projectMapper) FetchOne(id string) (*Project, error) {

	col := C(projectCollection)
	defer col.Database.Session.Close()

	project := new(Project)
	if err := col.FindId(bson.ObjectIdHex(id)).One(project); err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	return project, nil
}

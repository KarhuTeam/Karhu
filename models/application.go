package models

import (
	goerrors "errors"
	"github.com/gotoolz/errors"
	"github.com/gotoolz/validator"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type applicationMapper struct{}

var ApplicationMapper = &applicationMapper{}

const applicationCollection = "application"

type Application struct {
	Id          bson.ObjectId     `json:"id" bson:"_id"`
	Name        string            `json:"name" bson:"name"`
	Description string            `json:"description" bson:"description"`
	Tags        []string          `json:"tags" bson:"tags"` // Tags are used for application search
	Vars        map[string]string `json:"vars" bson:"vars"` // Vars are set in env when deploying a application
	CreatedAt   time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" bson:"updated_at"`
}

func (p *Application) Update(f *ApplicationUpdateForm) {

	p.Name = f.Name
	p.Description = f.Description
	p.Tags = f.Tags
	p.Vars = f.Vars
}

type Applications []*Application

// Application creation form
type ApplicationCreateForm struct {
	Name        string            `form:"name" json:"name" valid:"ascii,required"`
	Description string            `form:"description" json:"description" valid:"ascii"`
	Tags        []string          `form:"tags[]" json:"tags" valid:"-"`
	Vars        map[string]string `form:"-" json:"vars" valid:"-"`
	VarKeys     []string          `form:"varKeys[]" json:"-" valid:"-"`
	VarValues   []string          `form:"varValues[]" json:"-" valid:"-"`
}

// Validator for application creation
func (f ApplicationCreateForm) Validate() *errors.Errors {
	return validator.Validate(&f)
}

// Application update form
type ApplicationUpdateForm struct {
	Name        string            `json:"name" valid:"ascii,required"`
	Description string            `json:"description" valid:"ascii"`
	Tags        []string          `json:"tags" valid:"-"`
	Vars        map[string]string `json:"vars" valid:"-"`
}

// Validator for application update
func (f ApplicationUpdateForm) Validate() error {
	return validator.Validate(&f)
}

func (pm *applicationMapper) Create(f *ApplicationCreateForm) *Application {

	if f.Vars == nil {
		f.Vars = make(map[string]string)
	}
	for i := 0; i < len(f.VarKeys); i++ {
		f.Vars[f.VarKeys[i]] = f.VarValues[i]
	}

	return &Application{
		Id:          bson.NewObjectId(),
		Name:        f.Name,
		Description: f.Description,
		Tags:        f.Tags,
		Vars:        f.Vars,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (pm *applicationMapper) Save(p *Application) error {

	col := C(applicationCollection)
	defer col.Database.Session.Close()

	return col.Insert(p)
}

func (pm *applicationMapper) Update(p *Application) error {

	col := C(applicationCollection)
	defer col.Database.Session.Close()

	p.UpdatedAt = time.Now()

	return col.UpdateId(p.Id, p)
}

func (pm *applicationMapper) Delete(p *Application) error {

	col := C(applicationCollection)
	defer col.Database.Session.Close()

	return col.RemoveId(p.Id)
}

func (pm *applicationMapper) FetchAll() (Applications, error) {

	col := C(applicationCollection)
	defer col.Database.Session.Close()

	var applications Applications
	if err := col.Find(nil).All(&applications); err != nil {
		return nil, err
	}

	return applications, nil
}

func (pm *applicationMapper) FetchOne(id string) (*Application, error) {

	if !bson.IsObjectIdHex(id) {
		return nil, goerrors.New("Invalid id")
	}

	col := C(applicationCollection)
	defer col.Database.Session.Close()

	application := new(Application)
	if err := col.FindId(bson.ObjectIdHex(id)).One(application); err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	return application, nil
}

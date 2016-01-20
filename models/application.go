package models

import (
	"github.com/asaskevich/govalidator"
	"github.com/gotoolz/errors"
	"github.com/gotoolz/validator"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"regexp"
	"time"
)

type applicationMapper struct{}

var ApplicationMapper = &applicationMapper{}

const applicationCollection = "application"

var slugRegexp = regexp.MustCompile(`^[0-9a-z\-]+$`)

// Slug name validator
func init() {
	govalidator.TagMap["slug"] = govalidator.Validator(func(str string) bool {
		return slugRegexp.MatchString(str)
	})

	col := C(applicationCollection)
	defer col.Database.Session.Close()

	// App Index
	col.EnsureIndex(mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true, // See notes.
		Sparse:     true,
	})
}

type Application struct {
	Id          bson.ObjectId `json:"id" bson:"_id"`
	Name        string        `json:"name" bson:"name"` // Slug name
	Description string        `json:"description" bson:"description"`
	Tags        []string      `json:"tags" bson:"tags"` // Tags are used for application search
	CreatedAt   time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" bson:"updated_at"`
}

func (p *Application) Update(f *ApplicationUpdateForm) {

	p.Name = f.Name
	p.Description = f.Description
	p.Tags = f.Tags
}

type Applications []*Application

// Application creation form
type ApplicationCreateForm struct {
	Name        string   `form:"name" json:"name" valid:"slug,required"`
	Description string   `form:"description" json:"description" valid:"ascii"`
	Tags        []string `form:"tags[]" json:"tags" valid:"-"`
}

// Validator for application creation
func (f ApplicationCreateForm) Validate() *errors.Errors {
	if errs := validator.Validate(&f); errs != nil {
		return errs
	}

	app, err := ApplicationMapper.FetchOne(f.Name)
	if err != nil {
		panic(err)
	}

	if app != nil {
		return errors.New(errors.Error{
			Label: "duplicate_name",
			Field: "name",
			Text:  "Duplicate application name: " + f.Name,
		})
	}

	return nil
}

// Application update form
type ApplicationUpdateForm struct {
	Name        string   `form:"name" json:"name" valid:"slug,required"`
	Description string   `form:"description" json:"description" valid:"ascii"`
	Tags        []string `form:"tags[]" json:"tags" valid:"-"`
}

func (f *ApplicationUpdateForm) Hydrate(a *Application) {
	f.Name = a.Name
	f.Description = a.Description
	f.Tags = a.Tags
}

// Validator for application update
func (f ApplicationUpdateForm) Validate(app *Application) *errors.Errors {
	if errs := validator.Validate(&f); errs != nil {
		return errs
	}

	if f.Name != app.Name { // App identifier change
		app, err := ApplicationMapper.FetchOne(f.Name)
		if err != nil {
			panic(err)
		}

		if app != nil {
			return errors.New(errors.Error{
				Label: "duplicate_name",
				Field: "name",
				Text:  "Duplicate application name: " + f.Name,
			})
		}

	}

	return nil
}

func (pm *applicationMapper) Create(f *ApplicationCreateForm) *Application {

	return &Application{
		Id:          bson.NewObjectId(),
		Name:        f.Name,
		Description: f.Description,
		Tags:        f.Tags,
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

func (pm *applicationMapper) FetchOne(idOrSlug string) (*Application, error) {

	col := C(applicationCollection)
	defer col.Database.Session.Close()

	application := new(Application)
	if err := col.Find(findIdOrSlugQuery(idOrSlug)).One(application); err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	return application, nil
}

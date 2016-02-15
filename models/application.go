package models

import (
	"github.com/asaskevich/govalidator"
	"github.com/gotoolz/errors"
	"github.com/gotoolz/validator"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"regexp"
	"strings"
	"time"
)

type applicationMapper struct{}

var ApplicationMapper = &applicationMapper{}

const applicationCollection = "application"
const (
	APPLICATION_TYPE_APP     string = "app"
	APPLICATION_TYPE_SERVICE        = "service"
)

var slugRegexp = regexp.MustCompile(`^[0-9a-z\-]+$`)
var appTypes = []string{APPLICATION_TYPE_APP, APPLICATION_TYPE_SERVICE}

// Slug name validator
func init() {
	govalidator.TagMap["slug"] = govalidator.Validator(func(str string) bool {
		return slugRegexp.MatchString(str)
	})

	govalidator.TagMap["app_type"] = govalidator.Validator(func(str string) bool {

		for _, typ := range appTypes {
			if typ == str {
				return true
			}
		}

		return false
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
	Id          bson.ObjectId  `json:"id" bson:"_id"`
	Name        string         `json:"name" bson:"name"` // Slug name
	Type        string         `json:"type" bson:"type"`
	Description string         `json:"description" bson:"description"`
	Tags        []string       `json:"tags" bson:"tags"` // Tags are used for application search
	DepsIds     []string       `json:"-" bson:"deps"`
	Deps        []*Application `json:"deps" bson:"-"`
	CreatedAt   time.Time      `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" bson:"updated_at"`
}

func (p *Application) Update(f *ApplicationUpdateForm) {

	p.Name = f.Name
	p.Description = f.Description
	p.Tags = f.Tags

	var deps []*Application

	for _, d := range f.Deps {
		app, err := ApplicationMapper.FetchOne(d)
		if err != nil {
			panic(err)
		}

		deps = append(deps, app)
	}

	p.Deps = deps
}

type Applications []*Application

// Tags filter
type TagsFilter []string

func (tf TagsFilter) HasTag(key string) bool {
	for _, tag := range tf {
		if key == tag {
			return true
		}
	}
	return false
}

func (tf TagsFilter) Query(key string) string {

	tags := make([]string, 0)

	for _, tag := range tf {
		if key != tag && tag != "" {
			tags = append(tags, tag)
		}
	}

	if !tf.HasTag(key) {
		tags = append(tags, key)
	}

	return "?tags=" + strings.Join(tags, ",")
}

// Application creation form
type ApplicationCreateForm struct {
	Name        string   `form:"name" json:"name" valid:"slug,required"`
	Type        string   `form:"type" json:"type" valid:"app_type,required"`
	Description string   `form:"description" json:"description" valid:"ascii"`
	Tags        []string `form:"tags[]" json:"tags" valid:"-"`
	Deps        []string `form:"deps[]" json:"deps" valid:"-"` // Inter app deps

	// Services
	Packages []string `form:"packages[]" json:"packages" valid:"-"` // for service
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

	for _, d := range f.Deps {
		depApp, err := ApplicationMapper.FetchOne(d)
		if err != nil {
			panic(err)
		}

		if depApp == nil {
			return errors.New(errors.Error{
				Label: "invalid_dep",
				Field: "deps",
				Text:  "Invalid dependence: " + d,
			})
		}
	}

	return nil
}

// Application update form
type ApplicationUpdateForm struct {
	Name        string   `form:"name" json:"name" valid:"slug,required"`
	Description string   `form:"description" json:"description" valid:"ascii"`
	Tags        []string `form:"tags[]" json:"tags" valid:"-"`
	Deps        []string `form:"deps[]" json:"deps" valid:"-"` // Inter app deps

	// Services
	Packages []string `form:"packages[]" json:"packages" valid:"-"` // for service
}

func (f *ApplicationUpdateForm) Hydrate(a *Application) {
	f.Name = a.Name
	f.Description = a.Description
	f.Tags = a.Tags

	for _, d := range a.Deps {
		f.Deps = append(f.Deps, d.Name)
	}

	if a.Type == APPLICATION_TYPE_SERVICE {

		build, err := BuildMapper.FetchLast(a)
		if err != nil {
			panic(err)
		}

		if build == nil {
			panic("no build for application: " + a.Name)
		}

		f.Packages = build.RuntimeCfg.Dependencies
	}
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

	for _, d := range f.Deps {
		depApp, err := ApplicationMapper.FetchOne(d)
		if err != nil {
			panic(err)
		}

		if depApp == nil {
			return errors.New(errors.Error{
				Label: "invalid_dep",
				Field: "deps",
				Text:  "Invalid dependence: " + d,
			})
		}
	}

	return nil
}

func (am *applicationMapper) Create(f *ApplicationCreateForm) *Application {

	var deps []*Application

	for _, d := range f.Deps {
		app, err := am.FetchOne(d)
		if err != nil {
			panic(err)
		}

		deps = append(deps, app)
	}

	return &Application{
		Id:          bson.NewObjectId(),
		Name:        f.Name,
		Type:        f.Type,
		Description: f.Description,
		Tags:        f.Tags,
		Deps:        deps,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (pm *applicationMapper) Save(p *Application) error {

	col := C(applicationCollection)
	defer col.Database.Session.Close()

	p.DepsIds = nil
	for _, app := range p.Deps {

		p.DepsIds = append(p.DepsIds, app.Id.Hex())
	}

	return col.Insert(p)
}

func (pm *applicationMapper) Update(p *Application) error {

	col := C(applicationCollection)
	defer col.Database.Session.Close()

	p.UpdatedAt = time.Now()

	p.DepsIds = nil
	for _, app := range p.Deps {

		p.DepsIds = append(p.DepsIds, app.Id.Hex())
	}

	return col.UpdateId(p.Id, p)
}

func (pm *applicationMapper) Delete(p *Application) error {

	col := C(applicationCollection)
	defer col.Database.Session.Close()

	return col.RemoveId(p.Id)
}

func (am *applicationMapper) FetchAll() (Applications, error) {

	col := C(applicationCollection)
	defer col.Database.Session.Close()

	var applications Applications
	if err := col.Find(nil).All(&applications); err != nil {
		return nil, err
	}

	for _, application := range applications {
		for _, dep := range application.DepsIds {
			app, err := am.FetchOne(dep)
			if err != nil {
				panic(err)
			}

			application.Deps = append(application.Deps, app)
		}
	}

	return applications, nil
}

func (am *applicationMapper) FetchAllByTag(tags []string) (Applications, error) {

	col := C(applicationCollection)
	defer col.Database.Session.Close()

	var applications Applications
	var query *mgo.Query

	if len(tags) > 0 {
		query = col.Find(bson.M{"tags": bson.M{"$all": tags}})
	} else {
		query = col.Find(nil)
	}

	if err := query.All(&applications); err != nil {
		return nil, err
	}

	for _, application := range applications {
		for _, dep := range application.DepsIds {
			app, err := am.FetchOne(dep)
			if err != nil {
				panic(err)
			}

			application.Deps = append(application.Deps, app)
		}
	}

	return applications, nil
}

func (am *applicationMapper) FetchAllTags() ([]string, error) {

	col := C(applicationCollection)
	defer col.Database.Session.Close()

	var result []string

	if err := col.Find(nil).Distinct("tags", &result); err != nil {
		return result, err
	}

	return result, nil
}

func (am *applicationMapper) FetchOne(idOrSlug string) (*Application, error) {

	col := C(applicationCollection)
	defer col.Database.Session.Close()

	application := new(Application)
	if err := col.Find(findIdOrSlugQuery(idOrSlug)).One(application); err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	for _, dep := range application.DepsIds {
		app, err := am.FetchOne(dep)
		if err != nil {
			panic(err)
		}

		application.Deps = append(application.Deps, app)
	}

	return application, nil
}

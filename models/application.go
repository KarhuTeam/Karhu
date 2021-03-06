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
	APPLICATION_TYPE_DOCKER         = "docker"
)

var slugRegexp = regexp.MustCompile(`^[0-9a-z\-]+$`)
var appTypes = []string{APPLICATION_TYPE_APP, APPLICATION_TYPE_SERVICE, APPLICATION_TYPE_DOCKER}

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
	Id                  bson.ObjectId  `json:"id" bson:"_id"`
	Name                string         `json:"name" bson:"name"` // Slug name
	Type                string         `json:"type" bson:"type"`
	Description         string         `json:"description" bson:"description"`
	Tags                []string       `json:"tags" bson:"tags"` // Tags are used for application search
	DepsIds             []string       `json:"-" bson:"deps"`
	CurrentDeploymentId string         `json:"-" bson:"current_deployment_id"`
	Deps                []*Application `json:"deps" bson:"-"`
	CreatedAt           time.Time      `json:"created_at" bson:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at" bson:"updated_at"`
}

func (a *Application) Update(f *ApplicationUpdateForm) {

	a.Name = f.Name
	a.Description = f.Description
	a.Tags = f.Tags

	var deps []*Application

	for _, d := range f.Deps {
		app, err := ApplicationMapper.FetchOne(d)
		if err != nil {
			panic(err)
		}

		deps = append(deps, app)
	}

	a.Deps = deps
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
	ApplicationDockerForm
	Name        string   `form:"name" json:"name" valid:"slug,required"`
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

type ApplicationDockerForm struct {
	// Docker
	Image            string   `form:"image" json:"image" valid:"-"`
	Pull             string   `form:"pull" json:"pull" valid:"-"`
	PortsHost        []string `form:"ports-host[]" json:"ports-host" valid:"-"`
	PortsContainer   []string `form:"ports-container[]" json:"ports-container" valid:"-"`
	PortsProto       []string `form:"ports-proto[]" json:"ports-proto" valid:"-"`
	VolumesHost      []string `form:"volumes-host[]" json:"volumes-host" valid:"-"`
	VolumesContainer []string `form:"volumes-container[]" json:"volumes-container" valid:"-"`
	LinksContainer   []string `form:"links-container[]" json:"links-container" valid:"-"`
	LinksAlias       []string `form:"links-alias[]" json:"links-alias" valid:"-"`
	EnvKey           []string `form:"env-key[]" json:"env-key" valid:"-"`
	EnvValue         []string `form:"env-value[]" json:"env-value" valid:"-"`
	AutoRestart      string   `form:"restart" json:"restart" valid:"-"`
}

// Application update form
type ApplicationUpdateForm struct {
	ApplicationDockerForm
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

		f.Packages = build.RuntimeCfg.Dependencies.ToString()
	} else if a.Type == APPLICATION_TYPE_DOCKER {

		build, err := BuildMapper.FetchLast(a)
		if err != nil {
			panic(err)
		}

		if build == nil {
			panic("no build for application: " + a.Name)
		}

		f.Image = build.RuntimeCfg.Docker.Image
		f.Pull = ""
		if build.RuntimeCfg.Docker.Pull == "always" {
			f.Pull = "on"
		}
		for _, port := range build.RuntimeCfg.Docker.Ports {
			cfg := strings.SplitN(port, ":", 2)
			f.PortsHost = append(f.PortsHost, cfg[0])
			p := strings.SplitN(cfg[1], "/", 2)
			f.PortsContainer = append(f.PortsContainer, p[0])
			f.PortsProto = append(f.PortsProto, p[1])
		}

		for _, volume := range build.RuntimeCfg.Docker.Volumes {
			cfg := strings.SplitN(volume, ":", 2)
			f.VolumesHost = append(f.VolumesHost, cfg[0])
			f.VolumesContainer = append(f.VolumesContainer, cfg[1])
		}
		for _, link := range build.RuntimeCfg.Docker.Links {
			cfg := strings.SplitN(link, ":", 2)
			f.LinksContainer = append(f.LinksContainer, cfg[0])
			f.LinksAlias = append(f.LinksAlias, cfg[1])
		}
		for _, env := range build.RuntimeCfg.Docker.Env {
			cfg := strings.SplitN(env, ":", 2)
			f.EnvKey = append(f.EnvKey, cfg[0])
			f.EnvValue = append(f.EnvValue, strings.TrimSpace(cfg[1]))
		}

		f.AutoRestart = build.RuntimeCfg.Docker.Restart
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

	typ := APPLICATION_TYPE_APP
	if len(f.Packages) > 0 {
		typ = APPLICATION_TYPE_SERVICE
	} else if len(f.Image) > 0 {
		typ = APPLICATION_TYPE_DOCKER
	}

	return &Application{
		Id:          bson.NewObjectId(),
		Name:        f.Name,
		Type:        typ,
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

func (a *Application) Deployment() (*Deployment, error) {

	if a.CurrentDeploymentId == "" {
		return nil, nil
	}

	return DeploymentMapper.FetchOne(a, a.CurrentDeploymentId)
}

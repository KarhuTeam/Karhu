package models

import (
	"github.com/gotoolz/errors"
	"github.com/gotoolz/validator"
	// "github.com/karhuteam/karhu/ressources/application"
	// "github.com/karhuteam/karhu/ressources/file"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	// "io/ioutil"
	// "mime/multipart"
	// "net/http"
	"time"
)

type configMapper struct{}

var ConfigMapper = &configMapper{}

const configCollection = "config"

func init() {
	col := C(configCollection)
	defer col.Database.Session.Close()

	// App Index
	col.EnsureIndex(mgo.Index{
		Key:        []string{"application_id"},
		Unique:     false,
		DropDups:   false,
		Background: true, // See notes.
		Sparse:     true,
	})
}

type Config struct {
	Id            bson.ObjectId `json:"id" bson:"_id"`
	ApplicationId bson.ObjectId `json:"application_id" bson:"application_id"`
	Path          string        `json:"path" bson:"path"`
	Content       string        `json:"content" bson:"content"`
	Enabled       bool          `json:"enabled" bson:"enabled"`
	CreatedAt     time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at" bson:"updated_at"`
}

type Configs []*Config

// Application creation form
type ConfigCreateForm struct {
	Path    string `form:"path" json:"path" valid:"ascii,required"`
	Content string `form:"content" json:"content" valid:"ascii,required"`
	Enabled bool   `form:"enabled" json:"enabled" valid:"-"`
}

// Validator for application creation
func (f ConfigCreateForm) Validate() *errors.Errors {
	return validator.Validate(&f)
}

// Application creation form
type ConfigUpdateForm struct {
	Path    string `form:"path" json:"path" valid:"ascii,required"`
	Content string `form:"content" json:"content" valid:"ascii,required"`
	Enabled bool   `form:"enabled" json:"enabled" valid:"-"`
}

// Validator for application creation
func (f ConfigUpdateForm) Validate() *errors.Errors {
	return validator.Validate(&f)
}

func (f *ConfigUpdateForm) Hydrate(c *Config) {
	f.Path = c.Path
	f.Content = c.Content
	f.Enabled = c.Enabled
}

func (c *Config) Update(form *ConfigUpdateForm) {
	c.Path = form.Path
	c.Content = form.Content
	c.Enabled = form.Enabled
}

func (cm *configMapper) Create(app *Application, form *ConfigCreateForm) *Config {

	return &Config{
		Id:            bson.NewObjectId(),
		ApplicationId: app.Id,
		Path:          form.Path,
		Content:       form.Content,
		Enabled:       form.Enabled,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func (cm *configMapper) Save(c *Config) error {

	col := C(configCollection)
	defer col.Database.Session.Close()

	return col.Insert(c)
}

func (cm *configMapper) Update(c *Config) error {

	col := C(configCollection)
	defer col.Database.Session.Close()

	c.UpdatedAt = time.Now()

	return col.UpdateId(c.Id, c)
}

func (bm *configMapper) Delete(c *Config) error {

	col := C(configCollection)
	defer col.Database.Session.Close()

	return col.RemoveId(c.Id)
}

func (cm *configMapper) FetchAll(app *Application) (Configs, error) {

	col := C(configCollection)
	defer col.Database.Session.Close()

	var configs Configs
	if err := col.Find(bson.M{"application_id": app.Id}).Sort("-created_at").All(&configs); err != nil {
		return nil, err
	}

	return configs, nil
}

func (cm *configMapper) FetchAllEnabled(app *Application) (Configs, error) {

	col := C(configCollection)
	defer col.Database.Session.Close()

	var configs Configs
	if err := col.Find(bson.M{"application_id": app.Id, "enabled": true}).All(&configs); err != nil {
		return nil, err
	}

	return configs, nil
}

func (cm *configMapper) FetchOne(app *Application, configId string) (*Config, error) {

	if !bson.IsObjectIdHex(configId) {

		return nil, errors.New(errors.Error{
			Label: "invalid_config_id",
			Field: "config_id",
			Text:  "Invalid config id hex",
		})
	}

	col := C(configCollection)
	defer col.Database.Session.Close()

	config := new(Config)
	if err := col.Find(bson.M{"application_id": app.Id, "_id": bson.ObjectIdHex(configId)}).One(config); err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	return config, nil
}

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

type logfileMapper struct{}

var LogfileMapper = &logfileMapper{}

const logfileCollection = "logfile"

func init() {

	col := C(logfileCollection)
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

type Logfile struct {
	Id            bson.ObjectId `json:"id" bson:"_id"`
	ApplicationId bson.ObjectId `json:"application_id" bson:"application_id"`
	Path          string        `json:"path" bson:"path"`
	Enabled       bool          `json:"enabled" bson:"enabled"`
	CreatedAt     time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at" bson:"updated_at"`
}

type Logfiles []*Logfile

// Application creation form
type LogfileCreateForm struct {
	Path    string `form:"path" json:"path" valid:"ascii,required"`
	Enabled bool   `form:"enabled" json:"enabled" valid:"-"`
}

// Validator for application creation
func (f LogfileCreateForm) Validate() *errors.Errors {
	if errs := validator.Validate(&f); errs != nil {
		return errs
	}

	return nil
}

// Application creation form
type LogfileUpdateForm struct {
	Path    string `form:"path" json:"path" valid:"ascii,required"`
	Enabled bool   `form:"enabled" json:"enabled" valid:"-"`
}

// Validator for application creation
func (f LogfileUpdateForm) Validate() *errors.Errors {

	if errs := validator.Validate(&f); errs != nil {
		return errs
	}

	return nil
}

func (f *LogfileUpdateForm) Hydrate(lf *Logfile) {
	f.Path = lf.Path
	f.Enabled = lf.Enabled
}

func (lf *Logfile) Update(form *LogfileUpdateForm) {
	lf.Path = form.Path
	lf.Enabled = form.Enabled
}

func (lfm *logfileMapper) Create(app *Application, form *LogfileCreateForm) *Logfile {

	lf := &Logfile{
		Id:            bson.NewObjectId(),
		ApplicationId: app.Id,
		Path:          form.Path,
		Enabled:       form.Enabled,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return lf
}

func (lfm *logfileMapper) Save(lf *Logfile) error {

	col := C(logfileCollection)
	defer col.Database.Session.Close()

	if err := col.Insert(lf); err != nil {
		return err
	}

	if err := LogstashRefreshApplicationsFilters(); err != nil {
		return err
	}

	return nil
}

func (lfm *logfileMapper) Update(lf *Logfile) error {

	col := C(logfileCollection)
	defer col.Database.Session.Close()

	lf.UpdatedAt = time.Now()

	if err := col.UpdateId(lf.Id, lf); err != nil {
		return err
	}

	if err := LogstashRefreshApplicationsFilters(); err != nil {
		return err
	}
	return nil
}

func (lfm *logfileMapper) Delete(lf *Logfile) error {

	col := C(logfileCollection)
	defer col.Database.Session.Close()

	if err := col.RemoveId(lf.Id); err != nil {
		return err
	}

	if err := LogstashRefreshApplicationsFilters(); err != nil {
		return err
	}
	return nil
}

func (lfm *logfileMapper) FetchAll(app *Application) (Logfiles, error) {

	col := C(logfileCollection)
	defer col.Database.Session.Close()

	var logfiles Logfiles
	if err := col.Find(bson.M{"application_id": app.Id}).Sort("-created_at").All(&logfiles); err != nil {
		return nil, err
	}

	return logfiles, nil
}

func (lfm *logfileMapper) FetchAllEnabled(app *Application) (Logfiles, error) {

	col := C(logfileCollection)
	defer col.Database.Session.Close()

	var logfiles Logfiles
	if err := col.Find(bson.M{"application_id": app.Id, "enabled": true}).All(&logfiles); err != nil {
		return nil, err
	}

	return logfiles, nil
}

func (lfm *logfileMapper) FetchOne(app *Application, logfileId string) (*Logfile, error) {

	if !bson.IsObjectIdHex(logfileId) {

		return nil, errors.New(errors.Error{
			Label: "invalid_logfile_id",
			Field: "logfile_id",
			Text:  "Invalid logfile id hex",
		})
	}

	col := C(logfileCollection)
	defer col.Database.Session.Close()

	logfile := new(Logfile)
	if err := col.Find(bson.M{"application_id": app.Id, "_id": bson.ObjectIdHex(logfileId)}).One(logfile); err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	return logfile, nil
}

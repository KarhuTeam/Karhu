package models

import (
	"github.com/gotoolz/errors"
	// "github.com/gotoolz/validator"
	"github.com/karhuteam/karhu/ressources/file"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"mime/multipart"
	"net/http"
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
}

type Build struct {
	Id            bson.ObjectId `json:"id" bson:"_id"`
	ApplicationId bson.ObjectId `json:"application_id" bson:"application_id"`
	CommitHash    string        `json:"commit_hash" bson:"commit_id"`
	FilePath      string        `json:"file_path" bson:"file_path"`
	CreatedAt     time.Time     `json:"created_at" bson:"created_at"`
}

type Builds []*Build

func (b *Build) AttachFile(f multipart.File) error {

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return errors.New(errors.Error{
			Label: "internal_error",
			Field: "file",
			Text:  err.Error(),
		})
	}

	if contentType := http.DetectContentType(data); contentType != "application/zip" {
		return errors.New(errors.Error{
			Label: "internal_error",
			Field: "file",
			Text:  "Bad content-type, want application/zip have " + contentType,
		})
	}

	b.FilePath, err = file.Store("builds", b.Id.Hex()+".zip", data)
	return err
}

// // Build creation form
// type BuildCreateForm struct {
// 	Version    string                 `json:"version" valid:"ascii,required"`
// 	CommitHash string                 `json:"commit_hash" valid:"hexadecimal,required"`
// 	CommitUrl  string                 `json:"commit_url" valid:"url,required"`
// 	Tags       []string               `json:"tags" valid:"-"`
// 	Vars       map[string]interface{} `json:"vars" valid:"-"`
// }
//
// // Validator for build creation
// func (f BuildCreateForm) Validate() error {
// 	return validator.Validate(&f)
// }

func (bm *buildMapper) Create(app *Application, commitHash string) *Build {

	return &Build{
		Id:            bson.NewObjectId(),
		ApplicationId: app.Id,
		CommitHash:    commitHash,
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

func (bm *buildMapper) FetchAll(app *Application) (Builds, error) {

	col := C(buildCollection)
	defer col.Database.Session.Close()

	var builds Builds
	if err := col.Find(bson.M{"application_id": app.Id}).Sort("-created_at").All(&builds); err != nil {
		return nil, err
	}

	return builds, nil
}

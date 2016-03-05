package models

import (
	"github.com/asaskevich/govalidator"
	"github.com/gotoolz/errors"
	"github.com/gotoolz/validator"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type accessMapper struct{}

var AccessMapper = &accessMapper{}

const accessCollection = "access"

func init() {

	col := C(accessCollection)
	defer col.Database.Session.Close()

	// Name Index
	col.EnsureIndex(mgo.Index{
		Key:        []string{"type"},
		Unique:     true,
		DropDups:   true,
		Background: true, // See notes.
		Sparse:     true,
	})

	govalidator.TagMap["access_provider"] = govalidator.Validator(func(str string) bool {

		if str == "ec2" || str == "do" {
			return true
		}
		return false
	})
}

type Access struct {
	Id         bson.ObjectId `json:"id" bson:"_id"`
	Type       string        `json:"type" bson:"type"`
	AccessKey  string        `json:"access_key" bson:"access_key"`
	PrivateKey string        `json:"-" bson:"private_key"`
	CreatedAt  time.Time     `json:"created_at" bson:"created_at"`
}

type Accesses []*Access

// AlertGroup creation form
type AccessCreateForm struct {
	Type       string `json:"type" form:"type" valid:"access_provider,required"`
	AccessKey  string `json:"access_key" form:"access_key" valid:"required"`
	PrivateKey string `json:"private_key" form:"private_key" valid:"-"`
}

// Validator for application creation
func (f AccessCreateForm) Validate() *errors.Errors {
	if errs := validator.Validate(&f); errs != nil {
		return errs
	}

	a, err := AccessMapper.FetchOne(f.Type)
	if err != nil {
		panic(err)
	}

	if a != nil {
		return errors.New(errors.Error{
			Label: "duplicate_type",
			Field: "type",
			Text:  "Duplicate access type: " + f.Type,
		})
	}

	return nil
}

func (am *accessMapper) Create(f *AccessCreateForm) *Access {

	return &Access{
		Id:         bson.NewObjectId(),
		Type:       f.Type,
		AccessKey:  f.AccessKey,
		PrivateKey: f.PrivateKey,
		CreatedAt:  time.Now(),
	}
}

func (pm *accessMapper) Save(a *Access) error {

	col := C(accessCollection)
	defer col.Database.Session.Close()

	return col.Insert(a)
}

// func (pm *accessMapper) Update(a *Access) error {
//
// 	col := C(accessCollection)
// 	defer col.Database.Session.Close()
//
// 	ag.UpdatedAt = time.Now()
//
// 	return col.UpdateId(ag.Id, ag)
// }

func (pm *accessMapper) Delete(a *Access) error {

	col := C(accessCollection)
	defer col.Database.Session.Close()

	return col.RemoveId(a.Id)
}

func (am *accessMapper) FetchOne(typ string) (*Access, error) {

	col := C(accessCollection)
	defer col.Database.Session.Close()

	a := new(Access)
	if err := col.Find(bson.M{"type": typ}).One(a); err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	return a, nil
}

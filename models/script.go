package models

import (
	"github.com/wayt/govalidator"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func init() {
	col := C(scriptCollection)
	defer col.Database.Session.Close()

	// Project Id Index
	col.EnsureIndex(mgo.Index{
		Key:        []string{"project_id"},
		Unique:     false,
		DropDups:   false,
		Background: true, // See notes.
		Sparse:     true,
	})

	// Name Index
	col.EnsureIndex(mgo.Index{
		Key:        []string{"name"},
		Unique:     false,
		DropDups:   false,
		Background: true, // See notes.
		Sparse:     true,
	})
}

type scriptMapper struct{}

var ScriptMapper = &scriptMapper{}

const scriptCollection = "script"

type Script struct {
	Id        bson.ObjectId `json:"id" bson:"_id"`
	ProjectId bson.ObjectId `json:"project_id" bson:"project_id"`
	Name      string        `json:"name" bson:"name"`
	Content   string        `json:"content" bson:"content"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}

func (s *Script) Update(f *ScriptUpdateForm) {
	s.Content = f.Content
}

type Scripts []*Script

// Script creation form
type ScriptCreateForm struct {
	Name    string `json:"name" valid:"ascii,required"`
	Content string `json:"content" valid:"required"`
}

// Validator for script creation
func (f ScriptCreateForm) Validate() error {
	return govalidator.Validate(&f)
}

// Script update form
type ScriptUpdateForm struct {
	Content string `json:"content" valid:"required"`
}

// Validator for script update
func (f ScriptUpdateForm) Validate() error {
	return govalidator.Validate(&f)
}

func (sm *scriptMapper) Create(p *Project, f *ScriptCreateForm) *Script {

	return &Script{
		Id:        bson.NewObjectId(),
		ProjectId: p.Id,
		Name:      f.Name,
		Content:   f.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (sm *scriptMapper) Save(s *Script) error {

	col := C(scriptCollection)
	defer col.Database.Session.Close()

	return col.Insert(s)
}

func (sm *scriptMapper) Update(s *Script) error {

	col := C(scriptCollection)
	defer col.Database.Session.Close()

	s.UpdatedAt = time.Now()

	return col.UpdateId(s.Id, s)
}

func (sm *scriptMapper) Delete(s *Script) error {

	col := C(scriptCollection)
	defer col.Database.Session.Close()

	return col.RemoveId(s.Id)
}

func (sm *scriptMapper) FetchAll(p *Project) (Scripts, error) {

	col := C(scriptCollection)
	defer col.Database.Session.Close()

	var scripts Scripts
	if err := col.Find(bson.M{"project_id": p.Id}).All(&scripts); err != nil {
		return nil, err
	}

	return scripts, nil
}

func (sm *scriptMapper) FetchOne(p *Project, id string) (*Script, error) {

	if !bson.IsObjectIdHex(id) {
		return sm.FetchOneByName(p, id)
	}

	col := C(scriptCollection)
	defer col.Database.Session.Close()

	script := new(Script)
	if err := col.Find(bson.M{"_id": bson.ObjectIdHex(id), "project_id": p.Id}).One(script); err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	return script, nil
}

func (sm *scriptMapper) FetchOneByName(p *Project, name string) (*Script, error) {

	col := C(scriptCollection)
	defer col.Database.Session.Close()

	script := new(Script)
	if err := col.Find(bson.M{"name": name, "project_id": p.Id}).One(script); err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	return script, nil
}

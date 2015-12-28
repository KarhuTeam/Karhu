package models

//
// import (
// 	"github.com/gotoolz/validator"
// 	"gopkg.in/mgo.v2"
// 	"gopkg.in/mgo.v2/bson"
// 	"time"
// )
//
// func init() {
// 	col := C(scriptCollection)
// 	defer col.Database.Session.Close()
//
// 	// App Index
// 	col.EnsureIndex(mgo.Index{
// 		Key:        []string{"application_id"},
// 		Unique:     false,
// 		DropDups:   false,
// 		Background: true, // See notes.
// 		Sparse:     true,
// 	})
//
// 	// Name Index
// 	col.EnsureIndex(mgo.Index{
// 		Key:        []string{"environment_id", "name"},
// 		Unique:     true,
// 		DropDups:   false,
// 		Background: true, // See notes.
// 		Sparse:     true,
// 	})
// }
//
// type scriptMapper struct{}
//
// var ScriptMapper = &scriptMapper{}
//
// const scriptCollection = "script"
//
// type Script struct {
// 	Id            bson.ObjectId `json:"id" bson:"_id"`
// 	ApplicationId bson.ObjectId `json:"application_id" bson:"application_id"`
// 	EnvironmentId bson.ObjectId `json:"environment_id" bson:"environment_id"`
// 	Name          string        `json:"name" bson:"name"`
// 	Content       string        `json:"content" bson:"content"`
// 	CreatedAt     time.Time     `json:"created_at" bson:"created_at"`
// 	UpdatedAt     time.Time     `json:"updated_at" bson:"updated_at"`
// }
//
// func (s *Script) Update(f *ScriptUpdateForm) {
// 	s.Content = f.Content
// }
//
// type Scripts []*Script
//
// // Script creation form
// type ScriptCreateForm struct {
// 	Name    string `json:"name" valid:"ascii,required"`
// 	Content string `json:"content" valid:"required"`
// }
//
// // Validator for script creation
// func (f ScriptCreateForm) Validate() error {
// 	return validator.Validate(&f)
// }
//
// // Script update form
// type ScriptUpdateForm struct {
// 	Content string `json:"content" valid:"required"`
// }
//
// // Validator for script update
// func (f ScriptUpdateForm) Validate() error {
// 	return validator.Validate(&f)
// }
//
// func (sm *scriptMapper) Create(e *Environment, f *ScriptCreateForm) *Script {
//
// 	return &Script{
// 		Id:            bson.NewObjectId(),
// 		EnvironmentId: e.Id,
// 		Name:          f.Name,
// 		Content:       f.Content,
// 		CreatedAt:     time.Now(),
// 		UpdatedAt:     time.Now(),
// 	}
// }
//
// func (sm *scriptMapper) Save(s *Script) error {
//
// 	col := C(scriptCollection)
// 	defer col.Database.Session.Close()
//
// 	return col.Insert(s)
// }
//
// func (sm *scriptMapper) Update(s *Script) error {
//
// 	col := C(scriptCollection)
// 	defer col.Database.Session.Close()
//
// 	s.UpdatedAt = time.Now()
//
// 	return col.UpdateId(s.Id, s)
// }
//
// func (sm *scriptMapper) Delete(s *Script) error {
//
// 	col := C(scriptCollection)
// 	defer col.Database.Session.Close()
//
// 	return col.RemoveId(s.Id)
// }
//
// func (sm *scriptMapper) FetchAll(e *Environment) (Scripts, error) {
//
// 	col := C(scriptCollection)
// 	defer col.Database.Session.Close()
//
// 	var scripts Scripts
// 	if err := col.Find(bson.M{"environment_id": e.Id}).All(&scripts); err != nil {
// 		return nil, err
// 	}
//
// 	return scripts, nil
// }
//
// func (sm *scriptMapper) FetchOne(e *Environment, id string) (*Script, error) {
//
// 	if !bson.IsObjectIdHex(id) {
// 		return sm.FetchOneByName(e, id)
// 	}
//
// 	col := C(scriptCollection)
// 	defer col.Database.Session.Close()
//
// 	script := new(Script)
// 	if err := col.Find(bson.M{"_id": bson.ObjectIdHex(id), "environment_id": e.Id}).One(script); err != nil {
// 		if err == mgo.ErrNotFound {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}
//
// 	return script, nil
// }
//
// func (sm *scriptMapper) FetchOneByName(e *Environment, name string) (*Script, error) {
//
// 	col := C(scriptCollection)
// 	defer col.Database.Session.Close()
//
// 	script := new(Script)
// 	if err := col.Find(bson.M{"name": name, "environment_id": e.Id}).One(script); err != nil {
// 		if err == mgo.ErrNotFound {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}
//
// 	return script, nil
// }

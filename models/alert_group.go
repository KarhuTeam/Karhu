package models

import (
	"github.com/gotoolz/errors"
	"github.com/gotoolz/validator"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type alertGroupMapper struct{}

var AlertGroupMapper = &alertGroupMapper{}

const alertGroupCollection = "alert_group"

func init() {

	col := C(alertGroupCollection)
	defer col.Database.Session.Close()

	// Name Index
	col.EnsureIndex(mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true, // See notes.
		Sparse:     true,
	})
}

type AlertGroupMethod struct {
	Type  string `json:"type" bson:"type"`
	Value string `json:"value" bson:"value"`
}

type AlertGroupMethods []AlertGroupMethod

type AlertGroup struct {
	Id          bson.ObjectId     `json:"id" bson:"_id"`
	Name        string            `json:"name" bson:"name"` // Slug name
	Description string            `json:"description" bson:"description"`
	Methods     AlertGroupMethods `json:"methods" bson:"methods"`
	CreatedAt   time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" bson:"updated_at"`
}

type AlertGroups []*AlertGroup

// AlertGroup creation form
type AlertGroupCreateForm struct {
	Name        string   `form:"name" json:"name" valid:"slug,required"`
	Description string   `form:"description" json:"description" valid:"ascii"`
	MethodType  []string `form:"method-type[]" json:"-" valid:"-"`
	MethodValue []string `form:"method-value[]" json:"-" valid:"-"`
}

// Validator for application creation
func (f AlertGroupCreateForm) Validate() *errors.Errors {
	if errs := validator.Validate(&f); errs != nil {
		return errs
	}

	ag, err := AlertGroupMapper.FetchOne(f.Name)
	if err != nil {
		panic(err)
	}

	if ag != nil {
		return errors.New(errors.Error{
			Label: "duplicate_name",
			Field: "name",
			Text:  "Duplicate alert group name: " + f.Name,
		})
	}

	if len(f.MethodType) == 0 {
		return errors.New(errors.Error{
			Label: "invalid_methods",
			Field: "method-type[],method-value[]",
			Text:  "Invalid method type-value",
		})
	}

	if len(f.MethodType) != len(f.MethodValue) {
		return errors.New(errors.Error{
			Label: "invalid_methods",
			Field: "method-type[],method-value[]",
			Text:  "Invalid method type-value",
		})
	}

	for _, value := range f.MethodValue {
		if value == "" {
			return errors.New(errors.Error{
				Label: "invalid_method_value",
				Field: "method-value[]",
				Text:  "Invalid method value",
			})
		}
	}

	return nil
}

// AlertGroup update form
type AlertGroupUpdateForm struct {
	Name        string   `form:"name" json:"name" valid:"slug,required"`
	Description string   `form:"description" json:"description" valid:"ascii"`
	MethodType  []string `form:"method-type[]" json:"-" valid:"-"`
	MethodValue []string `form:"method-value[]" json:"-" valid:"-"`
}

// Validator for AlertGroup update
func (f AlertGroupUpdateForm) Validate(ag *AlertGroup) *errors.Errors {
	if errs := validator.Validate(&f); errs != nil {
		return errs
	}

	if f.Name != ag.Name { // App identifier change
		group, err := AlertGroupMapper.FetchOne(f.Name)
		if err != nil {
			panic(err)
		}

		if group != nil {
			return errors.New(errors.Error{
				Label: "duplicate_name",
				Field: "name",
				Text:  "Duplicate application name: " + f.Name,
			})
		}
	}

	if len(f.MethodType) == 0 {
		return errors.New(errors.Error{
			Label: "invalid_methods",
			Field: "method-type[],method-value[]",
			Text:  "Invalid method type-value",
		})
	}

	if len(f.MethodType) != len(f.MethodValue) {
		return errors.New(errors.Error{
			Label: "invalid_methods",
			Field: "method-type[],method-value[]",
			Text:  "Invalid method type-value",
		})
	}

	for _, value := range f.MethodValue {
		if value == "" {
			return errors.New(errors.Error{
				Label: "invalid_method_value",
				Field: "method-value[]",
				Text:  "Invalid method value",
			})
		}
	}

	return nil
}

func (ag *AlertGroup) Update(f *AlertGroupUpdateForm) {

	ag.Name = f.Name
	ag.Description = f.Description

	var methods AlertGroupMethods

	for i := range f.MethodType {
		methods = append(methods, AlertGroupMethod{
			Type:  f.MethodType[i],
			Value: f.MethodValue[i],
		})
	}

	ag.Methods = methods
}

func (f *AlertGroupUpdateForm) Hydrate(ag *AlertGroup) {
	f.Name = ag.Name
	f.Description = ag.Description

	for _, method := range ag.Methods {
		f.MethodType = append(f.MethodType, method.Type)
		f.MethodValue = append(f.MethodValue, method.Value)
	}
}

func (agm *alertGroupMapper) Create(f *AlertGroupCreateForm) *AlertGroup {

	var methods AlertGroupMethods

	for i := range f.MethodType {
		methods = append(methods, AlertGroupMethod{
			Type:  f.MethodType[i],
			Value: f.MethodValue[i],
		})
	}

	return &AlertGroup{
		Id:          bson.NewObjectId(),
		Name:        f.Name,
		Description: f.Description,
		Methods:     methods,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (pm *alertGroupMapper) Save(ag *AlertGroup) error {

	col := C(alertGroupCollection)
	defer col.Database.Session.Close()

	return col.Insert(ag)
}

func (pm *alertGroupMapper) Update(ag *AlertGroup) error {

	col := C(alertGroupCollection)
	defer col.Database.Session.Close()

	ag.UpdatedAt = time.Now()

	return col.UpdateId(ag.Id, ag)
}

func (pm *alertGroupMapper) Delete(ag *AlertGroup) error {

	col := C(alertGroupCollection)
	defer col.Database.Session.Close()

	return col.RemoveId(ag.Id)
}

func (am *alertGroupMapper) FetchAll() (AlertGroups, error) {

	col := C(alertGroupCollection)
	defer col.Database.Session.Close()

	var alertGroups AlertGroups
	if err := col.Find(nil).All(&alertGroups); err != nil {
		return nil, err
	}

	return alertGroups, nil
}

func (am *alertGroupMapper) FetchOne(idOrSlug string) (*AlertGroup, error) {

	col := C(alertGroupCollection)
	defer col.Database.Session.Close()

	ag := new(AlertGroup)
	if err := col.Find(findIdOrSlugQuery(idOrSlug)).One(ag); err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	return ag, nil
}

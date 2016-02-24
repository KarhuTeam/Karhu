package models

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	STATUS_OPEN        string = "open"
	STATUS_ACKNOWLEDGE        = "acknowledge"
	STATUS_CLOSED             = "closed"
)

var AlertStatus = []string{STATUS_OPEN, STATUS_ACKNOWLEDGE, STATUS_CLOSED}

type AlertMessage struct {
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	Text      string    `json:"text" bson:"text"`
}

type Alert struct {
	Id       bson.ObjectId  `json:"id" bson:"_id"`
	Name     string         `json:"name" bson:"name"`
	Status   string         `json:"status" bson:"status"`
	PolicyId bson.ObjectId  `json:"policy_id" bson:"policy_id"`
	Messages []AlertMessage `json:"messages" bson:"messages"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type Alerts []*Alert

func (a *Alert) AddMessage(text string) {
	a.Messages = append(a.Messages, AlertMessage{
		CreatedAt: time.Now(),
		Text:      text,
	})
}

type alertMapper struct{}

var AlertMapper = &alertMapper{}

const alertCollection = "alert"

func (am *alertMapper) Create(policy *AlertPolicy, reason error) *Alert {

	a := &Alert{
		Id:        bson.NewObjectId(),
		Name:      policy.Name,
		Status:    STATUS_OPEN,
		PolicyId:  policy.Id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	a.AddMessage(reason.Error())
	return a
}

func (am *alertMapper) Save(a *Alert) error {

	col := C(alertCollection)
	defer col.Database.Session.Close()

	return col.Insert(a)
}

func (am *alertMapper) Update(a *Alert) error {

	col := C(alertCollection)
	defer col.Database.Session.Close()

	a.UpdatedAt = time.Now()

	return col.UpdateId(a.Id, a)
}

func (am *alertMapper) Delete(a *Alert) error {

	col := C(alertCollection)
	defer col.Database.Session.Close()

	return col.RemoveId(a.Id)
}

func (am *alertMapper) FetchAll() (Alerts, error) {

	col := C(alertCollection)
	defer col.Database.Session.Close()

	var alerts Alerts
	if err := col.Find(nil).Sort("-created_at").All(&alerts); err != nil {
		return nil, err
	}

	return alerts, nil
}

func (am *alertMapper) FetchAllStatus(status string) (Alerts, error) {

	if status == "all" {
		return am.FetchAll()
	}

	col := C(alertCollection)
	defer col.Database.Session.Close()

	var alerts Alerts
	if err := col.Find(bson.M{"status": status}).Sort("-created_at").All(&alerts); err != nil {
		return nil, err
	}

	return alerts, nil
}

func (am *alertMapper) FetchOneByPolicy(policyId string, status []string) (*Alert, error) {

	if !bson.IsObjectIdHex(policyId) {
		return nil, errors.New("Invalid object Id Hex")
	}

	col := C(alertCollection)
	defer col.Database.Session.Close()

	a := new(Alert)
	if err := col.Find(bson.M{"policy_id": bson.ObjectIdHex(policyId), "status": bson.M{"$in": status}}).One(a); err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	return a, nil
}

func (am *alertMapper) FetchOne(id string) (*Alert, error) {

	if !bson.IsObjectIdHex(id) {
		return nil, errors.New("Invalid object Id Hex")
	}

	col := C(alertCollection)
	defer col.Database.Session.Close()

	a := new(Alert)
	if err := col.FindId(bson.ObjectIdHex(id)).One(a); err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	return a, nil
}

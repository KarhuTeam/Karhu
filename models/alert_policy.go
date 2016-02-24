package models

import (
	goerr "errors"
	"github.com/asaskevich/govalidator"
	"github.com/gotoolz/errors"
	"github.com/gotoolz/validator"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

var hostnameRegexp = regexp.MustCompile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)

type alertPolicyMapper struct{}

var AlertPolicyMapper = &alertPolicyMapper{}

const alertPolicyCollection = "alert_policy"

func init() {

	govalidator.TagMap["cond_type"] = govalidator.Validator(func(str string) bool {
		for _, c := range CondTypes {
			if str == c {
				return true
			}
		}
		return false
	})

	govalidator.TagMap["http_proto"] = govalidator.Validator(func(str string) bool {
		for _, p := range []string{"http", "https"} {
			if str == p {
				return true
			}
		}
		return false
	})

	govalidator.TagMap["http_path"] = govalidator.Validator(func(str string) bool {

		if _, err := url.Parse(str); err != nil {
			return false
		}
		return true
	})

	govalidator.TagMap["http_hostname"] = govalidator.Validator(func(str string) bool {

		return hostnameRegexp.MatchString(str)
	})

	govalidator.TagMap["target_type"] = govalidator.Validator(func(str string) bool {
		for _, c := range TargetTypes {
			if str == c {
				return true
			}
		}
		return false
	})

	// col := C(alertPolicyCollection)
	// defer col.Database.Session.Close()

	// // Name Index
	// col.EnsureIndex(mgo.Index{
	// 	Key:        []string{"name"},
	// 	Unique:     true,
	// 	DropDups:   true,
	// 	Background: true, // See notes.
	// 	Sparse:     true,
	// })
}

var CondTypes = []string{"cond-http", "cond-metric"}
var TargetTypes = []string{"target-url", "target-tag", "target-node", "target-all"}

type AlertPolicy struct {
	Id          bson.ObjectId `json:"id" bson:"_id"`
	Name        string        `json:"name" bson:"name"` // Slug name
	Description string        `json:"description" bson:"description"`
	AlertGroups []string      `json:"alert_groups" bson:"alert_groups"`
	CondType    string        `json:"cond_type" bson:"cond_type"`
	Cond        interface{}   `json:"cond" bson:"cond"`
	TargetType  string        `json:"target_type" bson:"target_type"`
	Target      interface{}   `json:"target" bson:"target"`
	Interval    time.Duration `json:"interval" bson:"interval"`
	CreatedAt   time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" bson:"updated_at"`
	NextAt      time.Time     `json:"next_at" bson:"next_at"`
}

type AlertPolicies []*AlertPolicy

type ConditionTypeHttp struct {
	Proto    string `json:"proto" bson:"proto"`
	Path     string `json:"path" bson:"path"`
	Hostname string `json:"hostname" bson:"hostname"`
	Port     int    `json:"port" bson:"port"`
	Status   int    `json:"status" bson:"status"`
}

type TargetTypeUrl struct {
	Host string `json:"host" bson:"host"`
}

type TargetTypeTag struct {
	Tags []string `json:"tags" bson:"tags"`
}

type TargetTypeNode struct {
	Nodes []string `json:"nodes" bson:"nodes"`
}

// AlertGroup creation form
type AlertPolicyForm struct {
	Name        string   `form:"name" json:"name" valid:"ascii,required"`
	Description string   `form:"description" json:"description" valid:"ascii"`
	AlertGroups []string `form:"alert-groups[]" json:"alert_groups" valid:"-"`

	CondType string `form:"cond-type" json:"cond_type" valid:"cond_type,required"`

	// CondType: cond-http
	CondHttpProto    string `form:"cond-http-proto" json:"cond_http_proto" valid:"http_proto"`
	CondHttpPath     string `form:"cond-http-path" json:"cond_http_path" valid:"http_path"`
	CondHttpHostname string `form:"cond-http-hostname" json:"cond_http_hostname" valid:"http_hostname"`
	CondHttpPort     string `form:"cond-http-port" json:"cond_http_port" valid:"port"`
	CondHttpStatus   string `form:"cond-http-status" json:"cond_http_status" valid:"int"`

	TargetType string `form:"target-type" json:"target_type" valid:"target_type,required"`

	// TargetType: target-url
	TargetUrlHost string `form:"target-url-host" json:"target_url_host" valid:"url"`

	// TargetType: target-tag
	TargetTagTags []string `form:"target-tag-tags[]" json:"target_tag_tags" valid:"-"`

	// TargetType: target-node
	TargetNodeNodes []string `form:"target-node-nodes[]" json:"target_node_nodes" valid:"-"`

	Interval string `form:"interval" json:"interval" valid:"int,required"`
}

// Validator for application creation
func (f AlertPolicyForm) Validate() *errors.Errors {
	if errs := validator.Validate(&f); errs != nil {
		return errs
	}

	if len(f.AlertGroups) == 0 {
		return errors.New(errors.Error{
			Label: "unknown_alert_groups",
			Field: "alert_groups",
			Text:  "Missing alert groups",
		})
	}

	for _, agName := range f.AlertGroups {

		ag, err := AlertGroupMapper.FetchOne(agName)
		if err != nil {
			panic(err)
		}

		if ag == nil {
			return errors.New(errors.Error{
				Label: "unknown_alert_groups",
				Field: "alert_groups",
				Text:  "Unknown alert group: " + agName,
			})
		}
	}

	return nil
}

func (f *AlertPolicyForm) Hydrate(ap *AlertPolicy) {

	f.Name = ap.Name
	f.Description = ap.Description
	f.AlertGroups = ap.AlertGroups
	f.Interval = strconv.Itoa(int(ap.Interval / time.Second))
	f.TargetType = ap.TargetType
	switch f.TargetType {
	case "target-url":
		f.TargetUrlHost = ap.Target.(bson.M)["host"].(string)
	case "target-tag":
		for _, tag := range ap.Target.(bson.M)["tags"].([]interface{}) {
			f.TargetTagTags = append(f.TargetTagTags, tag.(string))
		}
	case "target-node":
		for _, node := range ap.Target.(bson.M)["nodes"].([]interface{}) {
			f.TargetNodeNodes = append(f.TargetNodeNodes, node.(string))
		}
	case "target-all":
		// clap
	}
	f.CondType = ap.CondType
	switch f.CondType {
	case "cond-http":
		f.CondHttpHostname = ap.Cond.(bson.M)["hostname"].(string)
		f.CondHttpPath = ap.Cond.(bson.M)["path"].(string)
		f.CondHttpPort = strconv.Itoa(ap.Cond.(bson.M)["port"].(int))
		f.CondHttpProto = ap.Cond.(bson.M)["proto"].(string)
		f.CondHttpStatus = strconv.Itoa(ap.Cond.(bson.M)["status"].(int))
	case "cond-metric":
		// TO BE IMPLEMENTED
	}
}

func (apm *alertPolicyMapper) createCond(f *AlertPolicyForm) interface{} {

	switch f.CondType {
	case "cond-http":
		port, _ := strconv.Atoi(f.CondHttpPort)
		status, _ := strconv.Atoi(f.CondHttpStatus)
		return &ConditionTypeHttp{
			Proto:    f.CondHttpProto,
			Path:     f.CondHttpPath,
			Hostname: f.CondHttpHostname,
			Port:     port,
			Status:   status,
		}
	}

	return nil
}

func (apm *alertPolicyMapper) createTarget(f *AlertPolicyForm) interface{} {

	switch f.TargetType {
	case "target-url":
		return &TargetTypeUrl{
			Host: f.TargetUrlHost,
		}
	case "target-tag":
		return &TargetTypeTag{
			Tags: f.TargetTagTags,
		}
	case "target-node":
		return &TargetTypeNode{
			Nodes: f.TargetNodeNodes,
		}
	case "target-all":
		return nil
	}

	return nil
}

func (apm *alertPolicyMapper) Create(f *AlertPolicyForm) *AlertPolicy {

	interval, _ := strconv.Atoi(f.Interval)

	ap := &AlertPolicy{
		Id:          bson.NewObjectId(),
		Name:        f.Name,
		Description: f.Description,
		AlertGroups: f.AlertGroups,
		CondType:    f.CondType,
		Cond:        apm.createCond(f),
		TargetType:  f.TargetType,
		Target:      apm.createTarget(f),
		Interval:    time.Second * time.Duration(interval),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		NextAt:      time.Now(),
	}

	return ap
}

func (ap *AlertPolicy) Update(f *AlertPolicyForm) {

	interval, _ := strconv.Atoi(f.Interval)

	ap.Name = f.Name
	ap.Description = f.Description
	ap.AlertGroups = f.AlertGroups
	ap.CondType = f.CondType
	ap.Cond = AlertPolicyMapper.createCond(f)
	ap.TargetType = f.TargetType
	ap.Target = AlertPolicyMapper.createTarget(f)
	ap.Interval = time.Second * time.Duration(interval)
	ap.NextAt = time.Now()
}

func (apm *alertPolicyMapper) Save(ap *AlertPolicy) error {

	col := C(alertPolicyCollection)
	defer col.Database.Session.Close()

	return col.Insert(ap)
}

func (apm *alertPolicyMapper) Update(ap *AlertPolicy) error {

	col := C(alertPolicyCollection)
	defer col.Database.Session.Close()

	ap.UpdatedAt = time.Now()

	return col.UpdateId(ap.Id, ap)
}

func (apm *alertPolicyMapper) Delete(ap *AlertPolicy) error {

	col := C(alertPolicyCollection)
	defer col.Database.Session.Close()

	return col.RemoveId(ap.Id)
}

func (apm *alertPolicyMapper) FetchAll() (AlertPolicies, error) {

	col := C(alertPolicyCollection)
	defer col.Database.Session.Close()

	var alertPolicies AlertPolicies
	if err := col.Find(nil).All(&alertPolicies); err != nil {
		return nil, err
	}

	return alertPolicies, nil
}

func (apm *alertPolicyMapper) FetchAllForExecution() (AlertPolicies, error) {

	col := C(alertPolicyCollection)
	defer col.Database.Session.Close()

	now := time.Now()

	var alertPolicies AlertPolicies
	if err := col.Find(bson.M{
		"$or": []bson.M{
			{
				"next_at": bson.M{"$lte": now},
			},
			{
				"next_at": bson.M{"$exists": false},
			},
		},
	}).All(&alertPolicies); err != nil {
		return nil, err
	}

	return alertPolicies, nil
}

func (apm *alertPolicyMapper) FetchOne(id string) (*AlertPolicy, error) {

	if !bson.IsObjectIdHex(id) {
		return nil, goerr.New("Invalid object Id Hex")
	}

	col := C(alertPolicyCollection)
	defer col.Database.Session.Close()

	ap := new(AlertPolicy)
	if err := col.FindId(bson.ObjectIdHex(id)).One(ap); err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	return ap, nil
}

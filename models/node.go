package models

import (
	goerrors "errors"
	"github.com/gotoolz/errors"
	"github.com/gotoolz/validator"
	"github.com/karhuteam/karhu/ressources/ssh"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strconv"
	"time"
)

type nodeMapper struct{}

var NodeMapper = &nodeMapper{}

const nodeCollection = "node"

type Node struct {
	Id          bson.ObjectId `json:"id" bson:"_id"`
	Hostname    string        `json:"hostname" bson:"hostname"`
	IP          string        `json:"ip" bson:"ip"`
	SshPort     int           `json:"ssh_port" bson:"ssh_port"`
	SshUser     string        `json:"ssh_user" bson:"ssh_user"`
	Description string        `json:"description" bson:"description"`
	Tags        []string      `json:"tags" bson:"tags"` // Tags are used for node search
	CreatedAt   time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" bson:"updated_at"`
}

type Nodes []*Node

// Node creation form
type NodeCreateForm struct {
	Hostname    string   `form:"hostname" json:"hostname" valid:"ascii,required"`
	Description string   `form:"description" json:"description" valid:"ascii"`
	IP          string   `form:"ip" json:"ip" valid:"ip,required"`
	SshPort     string   `form:"ssh_port" json:"ssh_port" valid:"int"`
	SshUser     string   `form:"ssh_user" json:"ssh_user" valid:"ascii"`
	Tags        []string `form:"tags[]" json:"tags" valid:"-"`
}

type NodeUpdateForm struct {
	Description string   `form:"description" json:"description" valid:"-"`
	Tags        []string `form:"tags[]" json:"tags" valid:"-"`
}

// Validator for node creation
func (f NodeCreateForm) Validate() *errors.Errors {
	return validator.Validate(&f)
}

func (pm *nodeMapper) Create(f *NodeCreateForm) *Node {

	// Default ssh port
	sshPort, _ := strconv.Atoi(f.SshPort)
	if sshPort == 0 {
		sshPort = 22
	}

	// Default ssh user
	if f.SshUser == "" {
		f.SshUser = "root"
	}

	log.Println("form:", *f)

	return &Node{
		Id:          bson.NewObjectId(),
		Hostname:    f.Hostname,
		Description: f.Description,
		IP:          f.IP,
		SshPort:     sshPort,
		SshUser:     f.SshUser,
		Tags:        f.Tags,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (pm *nodeMapper) Save(n *Node) error {

	col := C(nodeCollection)
	defer col.Database.Session.Close()

	return col.Insert(n)
}

func (pm *nodeMapper) Update(n *Node) error {

	col := C(nodeCollection)
	defer col.Database.Session.Close()

	n.UpdatedAt = time.Now()

	return col.UpdateId(n.Id, n)
}

func (pm *nodeMapper) Delete(n *Node) error {

	col := C(nodeCollection)
	defer col.Database.Session.Close()

	return col.RemoveId(n.Id)
}

func (pm *nodeMapper) FetchAll() (Nodes, error) {

	col := C(nodeCollection)
	defer col.Database.Session.Close()

	var nodes Nodes
	if err := col.Find(nil).All(&nodes); err != nil {
		return nil, err
	}

	return nodes, nil
}

func (pm *nodeMapper) FetchOne(hostname string) (*Node, error) {

	col := C(nodeCollection)
	defer col.Database.Session.Close()

	node := new(Node)
	if err := col.Find(bson.M{"hostname": hostname}).One(node); err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	return node, nil
}

func (pm *nodeMapper) FetchOneById(id string) (*Node, error) {

	if !bson.IsObjectIdHex(id) {
		return nil, goerrors.New("Invalid id")
	}

	col := C(nodeCollection)
	defer col.Database.Session.Close()

	node := new(Node)
	if err := col.FindId(bson.ObjectIdHex(id)).One(node); err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	return node, nil
}

func (pm *nodeMapper) CheckSsh(n *Node) error {

	return ssh.CheckSsh(n.SshUser, n.IP, n.SshPort)
}

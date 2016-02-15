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

func init() {

	go NodeMapper.nodeStatusCheck()
}

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
	Status      string        `json:"status" bson:"status"`
	StatusAt    time.Time     `json:"status_at" bson:"status_at"`
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

	return &Node{
		Id:          bson.NewObjectId(),
		Hostname:    f.Hostname,
		Description: f.Description,
		IP:          f.IP,
		SshPort:     sshPort,
		SshUser:     f.SshUser,
		Tags:        f.Tags,
		Status:      STATUS_NEW,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (pm *nodeMapper) Save(n *Node) error {

	col := C(nodeCollection)
	defer col.Database.Session.Close()

	if err := col.Insert(n); err != nil {
		return err
	}

	if err := LogstashRefreshTagsFilters(); err != nil {
		return err
	}

	return nil
}

func (pm *nodeMapper) internalUpdate(n *Node) error {

	col := C(nodeCollection)
	defer col.Database.Session.Close()

	n.UpdatedAt = time.Now()

	if err := col.UpdateId(n.Id, n); err != nil {
		return err
	}

	return nil
}

func (pm *nodeMapper) Update(n *Node) error {

	col := C(nodeCollection)
	defer col.Database.Session.Close()

	n.UpdatedAt = time.Now()

	if err := col.UpdateId(n.Id, n); err != nil {
		return err
	}

	if err := LogstashRefreshTagsFilters(); err != nil {
		return err
	}

	return nil
}

func (pm *nodeMapper) Delete(n *Node) error {

	col := C(nodeCollection)
	defer col.Database.Session.Close()

	if err := col.RemoveId(n.Id); err != nil {
		return err
	}

	if err := LogstashRefreshTagsFilters(); err != nil {
		return err
	}

	return nil
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

func (pm *nodeMapper) FetchOneIP(ip string) (*Node, error) {

	col := C(nodeCollection)
	defer col.Database.Session.Close()

	node := new(Node)
	if err := col.Find(bson.M{"ip": ip}).One(node); err != nil {
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

func (nm *nodeMapper) FetchAllForApp(app *Application) (Nodes, error) {

	if len(app.Tags) == 0 {
		return nil, goerrors.New("No tags on app")
	}

	col := C(nodeCollection)
	defer col.Database.Session.Close()

	var nodes Nodes
	if err := col.Find(bson.M{"tags": bson.M{"$all": app.Tags}}).All(&nodes); err != nil {
		return nil, err
	}

	return nodes, nil
}

func (nm *nodeMapper) nodeStatusCheck() {

	log.Println("Run NodeStatusCheck routine")

	for {
		time.Sleep(5 * time.Second)

		func() {
			col := C(nodeCollection)
			defer col.Database.Session.Close()

			var nodes Nodes
			if err := col.Find(bson.M{"status_at": bson.M{"$lt": time.Now().Add(-time.Minute)}}).All(&nodes); err != nil {
				log.Println("nodeStatusCheck:", err)
			}

			for _, n := range nodes {

				log.Println("Check node:", *n)
				if err := nm.CheckSsh(n); err != nil {
					log.Println("nodeStatusCheck:", *n, err)
					n.Status = STATUS_ERROR
				} else {
					n.Status = STATUS_DONE
				}

				n.StatusAt = time.Now()
				if err := nm.internalUpdate(n); err != nil {
					log.Println("nodeStatusCheck:", *n, err)
				}
			}
		}()
	}
}

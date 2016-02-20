package models

import (
	"fmt"
	"github.com/gotoolz/errors"
	"github.com/karhuteam/ansible"
	"github.com/karhuteam/karhu/ressources/file"
	"github.com/karhuteam/karhu/ressources/ssh"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"time"
)

type deploymentMapper struct{}

var DeploymentMapper = &deploymentMapper{}

const deploymentCollection = "deployment"

func init() {
	col := C(deploymentCollection)
	defer col.Database.Session.Close()

	// App Index
	col.EnsureIndex(mgo.Index{
		Key:        []string{"application_id"},
		Unique:     false,
		DropDups:   false,
		Background: true, // See notes.
		Sparse:     true,
	})

	col.EnsureIndex(mgo.Index{
		Key:        []string{"build_id"},
		Unique:     false,
		DropDups:   false,
		Background: true, // See notes.
		Sparse:     true,
	})
}

type Deployment struct {
	Id            bson.ObjectId `json:"id" bson:"_id"`
	ApplicationId bson.ObjectId `json:"-" bson:"application_id"`
	Application   *Application  `json:"-" bson:"-"`
	BuildId       bson.ObjectId `json:"-" bson:"build_id"`
	Build         *Build        `json:"-" bson:"-"`
	TmpPath       string        `json:"-" bson:"tmp_path"`
	Logs          string        `json:"logs" bson:"logs"`
	Duration      time.Duration `json:"duration" bson:"duration"`
	CreatedAt     time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at" bson:"updated_at"`
	Status        string        `json:"status" bson:"status"`
}

type Deployments []*Deployment

func (dm *deploymentMapper) Create(app *Application, build *Build) *Deployment {

	return &Deployment{
		Id:            bson.NewObjectId(),
		ApplicationId: app.Id,
		Application:   app,
		BuildId:       build.Id,
		Build:         build,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Status:        STATUS_NEW,
	}
}

func (dm *deploymentMapper) FetchOne(app *Application, deployId string) (*Deployment, error) {

	if !bson.IsObjectIdHex(deployId) {

		return nil, errors.New(errors.Error{
			Label: "invalid_deploy_id",
			Field: "deploy_id",
			Text:  "Invalid deploy id hex",
		})
	}

	col := C(deploymentCollection)
	defer col.Database.Session.Close()

	deploy := new(Deployment)
	if err := col.Find(bson.M{"application_id": app.Id, "_id": bson.ObjectIdHex(deployId)}).One(deploy); err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	return deploy, nil
}

func (dm *deploymentMapper) FetchAll(app *Application) (Deployments, error) {

	col := C(deploymentCollection)
	defer col.Database.Session.Close()

	var deploys Deployments
	if err := col.Find(bson.M{"application_id": app.Id}).Sort("-created_at").All(&deploys); err != nil {
		return nil, err
	}

	return deploys, nil
}

func (dm *deploymentMapper) Save(d *Deployment) error {

	col := C(deploymentCollection)
	defer col.Database.Session.Close()

	return col.Insert(d)
}

func (dm *deploymentMapper) Update(d *Deployment) error {

	col := C(deploymentCollection)
	defer col.Database.Session.Close()

	d.UpdatedAt = time.Now()

	return col.UpdateId(d.Id, d)
}

func (d *Deployment) Run() {

	start := time.Now()

	// catch panic
	defer func() {
		if err := recover(); err != nil {
			log.Println("deployment.run:", err)

			trace := make([]byte, 2048)
			runtime.Stack(trace, true)
			log.Println("deployment.run:", string(trace))

			d.Status = STATUS_ERROR
			d.Duration = time.Since(start)
			if err := DeploymentMapper.Update(d); err != nil {
				log.Println("deployment.run:", err)
			}
		}
	}()

	d.Status = STATUS_RUNNING
	if err := DeploymentMapper.Update(d); err != nil {
		panic(err)
	}

	// Build playbook
	playbook, err := d.Playbook()
	if err != nil {
		panic(err)
	}

	hosts, err := NodeMapper.AnsibleHostsForTags(d.Application.Tags)
	if err != nil {
		panic(err)
	}

	if len(hosts) == 0 {
		d.Logs = "No hosts"
		panic("no hosts")
	}

	a, err := ansible.NewAnsible(playbook, hosts)
	if err != nil {
		panic(err)
	}
	// defer a.Clean()

	a.UseKey(ssh.PrivateKeyPath())

	if err := a.Write(); err != nil {
		panic(err)
	}

	// Extract files
	if err := d.extractArchive(a.Workdir); err != nil {
		panic(err)
	}

	// if err := ioutil.WriteFile(path.Join(tmpPath, "karhu.log"), []byte(fmt.Sprintf("Deployment %s\n", depl.Id.Hex())), 0644); err != nil {
	// 	panic(err)
	// }

	log.Println("ressources/ansible: Run: tmpPath:", a.Workdir)
	d.TmpPath = a.Workdir
	if err := DeploymentMapper.Update(d); err != nil {
		panic(err)
	}

	command := fmt.Sprintf("ansible-playbook -i hosts.ini playbook.yml > karhu.log 2>&1")
	cmd := exec.Command("sh", "-c", fmt.Sprintf("cd %s && %s", a.Workdir, command))

	cmdErr := cmd.Run()

	out, err := ioutil.ReadFile(path.Join(a.Workdir, "/karhu.log"))
	if err != nil {
		panic(err)
	}
	d.Logs = string(out)

	if cmdErr != nil {
		panic(err)
	}

	d.Status = STATUS_DONE
	d.Duration = time.Since(start)
	if err := DeploymentMapper.Update(d); err != nil {
		panic(err)
	}
}

func (d *Deployment) extractArchive(workdir string) error {

	// Fetch zip
	data, err := file.Get(d.Build.FilePath)
	if err != nil {
		return err
	}

	// Write zip file
	zipPath := path.Join(workdir, "karhu.zip")
	if err := ioutil.WriteFile(zipPath, data, 0644); err != nil {
		return err
	}

	destDir := path.Join(workdir, "karhu")
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	// Unzip file
	if err := file.Unzip(zipPath, destDir); err != nil {
		return err
	}

	return nil
}

func (d *Deployment) Playbook() (*ansible.Playbook, error) {

	playbook := ansible.NewPlaybook()

	role := ansible.NewRole("deployment").AddTask(ansible.Task{
		`name`:  `Setup User Group`,
		`group`: fmt.Sprintf(`name=%s system=yes`, d.Build.RuntimeCfg.User),
	}).AddTask(ansible.Task{
		`name`: `Setup User`,
		`user`: fmt.Sprintf(`name=%s group=%s system=yes`, d.Build.RuntimeCfg.User, d.Build.RuntimeCfg.User),
	})
	if d.Build.RuntimeCfg.Workdir != "" {
		role.AddTask(ansible.Task{
			`name`: `Setup Workdir`,
			`file`: fmt.Sprintf(`path=%s state=directory owner=%s group=%s`, d.Build.RuntimeCfg.Workdir, d.Build.RuntimeCfg.User, d.Build.RuntimeCfg.User),
		})
	}

	if d.Build.RuntimeCfg.Binary != nil {
		setupBinaryRole(role, d)
	}
	if len(d.Build.RuntimeCfg.Static) > 0 {
		setupStaticRole(role, d)
	}
	if len(d.Build.RuntimeCfg.Dependencies) > 0 {
		setupDependenciesRole(role, d)
	}

	// Check for application configs
	if configs, err := ConfigMapper.FetchAllEnabled(d.Application); err != nil {
		return nil, err
	} else if len(configs) > 0 {
		setupConfigsRole(role, d, configs)
	}

	playbook.AddRole(role)

	return playbook, nil
}

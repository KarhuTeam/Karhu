package alerts

import (
	"errors"
	// "fmt"
	"github.com/gotoolz/env"
	"github.com/karhuteam/karhu/models"
	"log"
	"os/exec"
	"path"
)

var nagiosPluginsPath = env.GetDefault("NAGIOS_PLUGINS_DIR", "/usr/lib/nagios/plugins")

type CheckNagios struct {
	Policy *models.AlertPolicy
	Plugin string
	Params string
}

func NewCheckNagios(policy *models.AlertPolicy) (*CheckNagios, error) {

	return &CheckNagios{
		Policy: policy,
		Plugin: policy.NagiosPlugin,
		Params: policy.NagiosParams,
	}, nil
}

func (c *CheckNagios) command() string {

	command := path.Join(nagiosPluginsPath, c.Plugin)
	if c.Params != "" {
		command += " " + c.Params
	}

	return command
}

func (c *CheckNagios) HandleTargetKarhu() error {
	log.Println("HandleTargetKarhu")

	cmd := exec.Command("sh", "-c", c.command())
	out, err := cmd.CombinedOutput()

	if err != nil {
		return errors.New(string(out))
	}
	return nil
}

func (c *CheckNagios) HandleTargetNode(node *models.Node) error {

	log.Println("HandleTargetNode")

	out, err := node.ExecSsh(c.command())

	log.Println("out:", string(out))
	if err != nil {
		return errors.New(string(out))
	}
	return nil

}

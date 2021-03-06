package models

import (
	"fmt"
	"github.com/karhuteam/ansible"
	"path"
)

var serviceFile = `##
# KARHU AUTO GENERATED FILE
# SHOULD NO BE MODIFIED ON SERVERS
# WILL BE OVERRIDE AT EACH DEPLOYMENTS
# LOVE <3
##

[Unit]
Description=%s - auto generate systemctl service file

[Service]
#EnvironmentFile=
WorkingDirectory=%s
User=%s
Group=%s

ExecStart=%s
Restart=on-failure

LimitNOFILE=65536

[Install]
WantedBy=multi-user.target`

func setupBinaryRole(role *ansible.Role, d *Deployment) *ansible.Role {

	// Add tasks
	role.AddTask(ansible.Task{
		`name`: `Setup Dest Dirs`,
		`file`: fmt.Sprintf(`path=%s owner=%s group=%s state=directory`, path.Join(d.Build.RuntimeCfg.Workdir, "/bin"), d.Build.RuntimeCfg.Binary.User, d.Build.RuntimeCfg.Binary.User),
	}).AddTask(ansible.Task{
		`name`:   `Setup Bin`,
		`copy`:   fmt.Sprintf(`src={{ ansible_workdir }}%s dest=%s mode=0755 owner=%s group=%s`, path.Join("/karhu/", d.Build.RuntimeCfg.Binary.Bin), path.Join(d.Build.RuntimeCfg.Workdir, "/bin/", d.Build.RuntimeCfg.Binary.Bin), d.Build.RuntimeCfg.Binary.User, d.Build.RuntimeCfg.Binary.User),
		`notify`: fmt.Sprintf(`restarted %s`, d.Application.Name),
	}).AddTask(ansible.Task{
		`name`:   `Setup systemctl Script`,
		`copy`:   fmt.Sprintf(`src=binary.service dest=/lib/systemd/system/%s.service`, d.Application.Name),
		`notify`: `Reload systemctl daemon`,
	})

	role.AddHandler(ansible.Task{
		`name`:  `Reload systemctl daemon`,
		`shell`: `/bin/systemctl daemon-reload`,
	}).AddHandler(ansible.Task{
		`name`:    fmt.Sprintf(`restarted %s`, d.Application.Name),
		`service`: fmt.Sprintf(`name=%s state=restarted enabled=yes`, d.Application.Name),
	})

	role.AddFile(ansible.NewFile("binary.service", []byte(fmt.Sprintf(serviceFile, d.Application.Name, d.Build.RuntimeCfg.Workdir, d.Build.RuntimeCfg.Binary.User, d.Build.RuntimeCfg.Binary.User, path.Join(d.Build.RuntimeCfg.Workdir, "/bin/", d.Build.RuntimeCfg.Binary.Bin)))))

	return role
}

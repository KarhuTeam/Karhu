package models

import (
	"crypto/sha1"
	"fmt"
	"github.com/karhuteam/ansible"
	"path"
)

func setupConfigsRole(role *ansible.Role, d *Deployment, configs Configs) *ansible.Role {

	role.AddTask(ansible.Task{
		`name`: `Setup rsync user`,
		`set_fact`: ansible.Task{
			`rsync_path`: `rsync`,
		},
		`when`: `ansible_ssh_user == "root"`,
	}).AddTask(ansible.Task{
		`name`: `Setup sudo rsync user`,
		`set_fact`: ansible.Task{
			`rsync_path`: `sudo rsync`,
		},
		`when`: `ansible_ssh_user != "root"`,
	}).AddTask(ansible.Task{
		`name`: `Make Sure rsync is installed`,
		`apt`:  `name=rsync state=present update_cache=yes cache_valid_time=86400`,
	})

	for _, conf := range configs {

		destPath := conf.Path
		if len(destPath) > 0 && destPath[0] != '/' {
			destPath = path.Join(d.Build.RuntimeCfg.Workdir, destPath)
		}
		hash := fmt.Sprintf("%x", sha1.Sum([]byte(conf.Path)))

		role.AddTask(ansible.Task{
			"name": "Create Directory",
			"file": fmt.Sprintf(`path=%s state=directory`),
		})

		task := ansible.Task{
			`name`:        `Copy ` + conf.Path,
			`sudo`:        `no`,
			`synchronize`: fmt.Sprintf(`src=%s dest=%s use_ssh_args=yes set_remote_user=yes recursive=yes delete=yes compress=yes mode=push checksum=yes times=no rsync_path={{ rsync_path }}`, hash, destPath),
		}
		if conf.Notify.Service != "" {
			task[`notify`] = fmt.Sprintf("%s %s", conf.Notify.State, conf.Notify.Service)
			role.AddHandler(ansible.Task{
				`name`:    task[`notify`],
				`service`: fmt.Sprintf(`name=%s state=%s`, conf.Notify.Service, conf.Notify.State),
			})
		}

		role.AddTask(task).AddFile(ansible.NewFile(hash, []byte(conf.Content)))
	}

	if d.Build.RuntimeCfg.Workdir != "" {
		role.AddTask(ansible.Task{
			`name`: `Setup Workdir Owner`,
			`file`: fmt.Sprintf(`path=%s recurse=yes group=%s owner=%s`, d.Build.RuntimeCfg.Workdir, d.Build.RuntimeCfg.User, d.Build.RuntimeCfg.User),
		})
	}

	return role
}

package models

import (
	"fmt"
	"github.com/karhuteam/ansible"
	"path"
)

func setupStaticRole(role *ansible.Role, d *Deployment) *ansible.Role {

	var files []map[string]string

	for _, s := range d.Build.RuntimeCfg.Static {

		src := path.Join("/karhu/", s.Src)
		dest := path.Join(d.Build.RuntimeCfg.Workdir, s.Dest)

		// files = append(files, fmt.Sprintf(`{ src: "{{ ansible_workdir }}%s", dest: "%s", mode: "%s", user: "%s" }`, src, dest, s.Mode, s.User))
		files = append(files, map[string]string{
			"src":  fmt.Sprintf("{{ ansible_workdir }}%s", src),
			"dest": dest,
			"mode": s.Mode,
		})
	}

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
	}).AddTask(ansible.Task{
		`name`:        `Copy Files`,
		`sudo`:        `no`,
		`synchronize`: `src={{ item.src }} dest={{ item.dest }} use_ssh_args=yes set_remote_user=yes recursive=yes delete=yes compress=yes mode=push checksum=yes times=no rsync_path={{ rsync_path }}`,
		`with_items`:  files,
	}).AddTask(ansible.Task{
		`name`: `Setup Workdir Owner`,
		`file`: fmt.Sprintf(`path=%s recurse=yes group=%s owner=%s`, d.Build.RuntimeCfg.Workdir, d.Build.RuntimeCfg.User, d.Build.RuntimeCfg.User),
	})

	return role
}

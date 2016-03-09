package models

import (
	"github.com/karhuteam/ansible"
)

func setupDockerRole(role *ansible.Role, d *Deployment) *ansible.Role {

	docker := d.Build.RuntimeCfg.Docker

	// Add tasks
	role.AddTask(ansible.Task{
		`name`:        `Check Docker Install`,
		`command`:     `dpkg-query -l docker-engine`,
		`failed_when`: `False`,
		`register`:    `deb_check`,
	}).AddTask(ansible.Task{
		`name`: `Download Docker Setup`,
		`get_url`: ansible.Task{
			`url`:  `https://get.docker.com`,
			`dest`: `/tmp/docker.sh`,
			`mode`: `0755`,
		},
		`when`: `deb_check.stderr.find('no packages found') != -1`,
	}).AddTask(ansible.Task{
		`name`:       `Setup Deps`,
		`apt`:        `name={{ item }} state=present force=yes update_cache=yes cache_valid_time=86400`,
		`with_items`: []string{"python-pip", "apt-transport-https"},
	}).AddTask(ansible.Task{
		`name`:    `Setup Docker`,
		`command`: `sh /tmp/docker.sh`,
		`when`:    `deb_check.stderr.find('no packages found') != -1`,
	}).AddTask(ansible.Task{
		`name`: `Setup docker-py`,
		`pip`:  `name=docker-py state=present`,
	}).AddTask(ansible.Task{
		`name`: `Launch Container`,
		`docker`: ansible.Task{
			`name`:           docker.Name,
			`image`:          docker.Image,
			`pull`:           docker.Pull,
			`restart_policy`: docker.Restart,
			`state`:          `reloaded`,
			`links`:          docker.Links,
			`ports`:          docker.Ports,
			`volumes`:        docker.Volumes,
			// `env`:            docker.Env,
		},
	})

	return role
}

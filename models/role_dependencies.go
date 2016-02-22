package models

import (
	"github.com/karhuteam/ansible"
)

func setupDependenciesRole(role *ansible.Role, d *Deployment) *ansible.Role {

	var dependencies []map[string]string

	for _, dep := range d.Build.RuntimeCfg.Dependencies {

		dependencies = append(dependencies, map[string]string{
			"name": dep.Name,
		})
	}

	role.AddTask(ansible.Task{
		`name`:       `Install Dependencies`,
		`apt`:        `name={{ item.name }} state=latest update_cache=yes cache_valid_time=86400`,
		`with_items`: dependencies,
	})

	return role
}

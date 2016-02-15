package ansible

import (
	"crypto/sha1"
	"fmt"
	"github.com/gotoolz/env"
	"github.com/karhuteam/karhu/models"
	"github.com/karhuteam/karhu/ressources/file"
	"github.com/karhuteam/karhu/ressources/ssh"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"text/template"
)

const (
	HOSTS_FILENAME    = "hosts.ini"
	PLAYBOOK_FILENAME = "playbook.yml"
	VARS_FILENAME     = "vars.yml"
	CONFIG_FILENAME   = "ansible.cfg"
)

var (
	hostFileTemplate     *template.Template
	playbookFileTemplate *template.Template
	varsFileTemplate     *template.Template
	configFileTemplate   *template.Template
)

func init() {

	var err error

	if hostFileTemplate, err = template.New(HOSTS_FILENAME).Parse(`[all]
{{ range .Nodes }}{{ .Hostname }} ansible_ssh_host={{ .IP }} ansible_ssh_port={{ .SshPort }} ansible_ssh_user={{ .SshUser }}
{{ end }}
`); err != nil {
		panic(err)
	}

	if playbookFileTemplate, err = template.New(PLAYBOOK_FILENAME).Parse(`---

- hosts: all
  sudo: yes
  vars_files:
  - {{ .Vars }}
  roles:{{ range .Roles }}
  - {{ . }}{{ end }}
`); err != nil {
		panic(err)
	}

	if varsFileTemplate, err = template.New(VARS_FILENAME).Parse(`---

application_name: {{ .Application.Name }}
runtime_type: {{ .RuntimeConfig.Type }}
runtime_user: {{ .RuntimeConfig.User }}
runtime_bin: {{ .RuntimeConfig.Bin }}
runtime_workdir: {{ .RuntimeConfig.Workdir }}
{{ if len .RuntimeConfig.Dependencies }}runtime_dependencies:{{ range .RuntimeConfig.Dependencies }}
  - {{ . }}{{ end }}{{ end }}
runtime_files: {{ if eq .RuntimeConfig.Type "binary" }}
  - { src: '{{ .TmpPath }}/karhu/{{ .RuntimeConfig.Bin }}', dest: '{{ .RuntimeConfig.Workdir }}/bin/{{ .RuntimeConfig.Bin }}', mode: '0755', user: '{{ .RuntimeConfig.User }}', notify: { service: '{{ .Application.Name }}', state: 'restarted' } }{{ end }}
{{ range $index, $str := .RuntimeConfig.Static }}  - { src: '{{ $.TmpPath }}/karhu/{{ $.RuntimeConfig.Static.Src $index }}', dest: '{{ $.RuntimeConfig.Workdir }}/{{ $.RuntimeConfig.Static.Dest $index}}', mode: '{{ $.RuntimeConfig.Static.Mode $index }}', user: '{{ $.RuntimeConfig.User }}', notify: { service: '{{ if eq $.RuntimeConfig.Type "binary" }}{{ $.Application.Name }}{{ end }}', state: 'restarted' } }
{{ end }}{{ range .Configs }}  - { src: '{{ .Src }}', dest: '{{ .Dest }}', destdir: '{{ .DestDir }}', mode: '{{ .Mode }}', user: 'root', notify: { service: '{{ .Notify.Service }}', state: '{{ .Notify.State }}' } }
{{ end }}
runtime_services:{{ range $dep := .Services }}{{ range $dep.RuntimeCfg.Dependencies }}
  - {{ . }}{{ end }}{{ end }}
runtime_services_files:{{ range $cfg := .ServicesConfigs }}
  - { src: '{{ $cfg.Src }}', dest: '{{ $cfg.Dest }}', destdir: '{{ $cfg.DestDir }}', mode: '{{ $cfg.Mode }}', notify: { service: '{{ .Notify.Service }}', state: '{{ .Notify.State }}' } }{{ end }}
`); err != nil {
		panic(err)
	}

	if configFileTemplate, err = template.New(CONFIG_FILENAME).Parse(`[defaults]

jinja2_extensions = jinja2.ext.loopcontrols
inventory      = hosts.ini

# uncomment this to disable SSH key host checking
host_key_checking = False

# SSH timeout
timeout = 10

# default user to use for playbooks if user is not specified
# (/usr/bin/ansible will use current user as default)
#remote_user = root

# logging is off by default unless this path is defined
# if so defined, consider logrotate
#log_path = /var/log/ansible.log

# default module name for /usr/bin/ansible
module_name = setup

# use this shell for commands executed under sudo
# you may need to change this to bin/bash in rare instances
# if sudo is constrained
#executable = /bin/sh

# if inventory variables overlap, does the higher precedence one win
# or are hash values merged together?  The default is 'replace' but
# this can also be set to 'merge'.
#hash_behaviour = replace

# list any Jinja2 extensions to enable here:
#jinja2_extensions = jinja2.ext.do,jinja2.ext.i18n

# if set, always use this private key file for authentication, same as
# if passing --private-key to ansible or ansible-playbook
# private_key_file =


# the http user-agent string to use when fetching urls. Some web server
# operators block the default urllib user agent as it is frequently used
# by malicious attacks/scripts, so we set it to something unique to
# avoid issues.
#http_user_agent = ansible-agent

# if set to a persistent type (not 'memory', for example 'redis') fact values
# from previous runs in Ansible will be stored.  This may be useful when
# wanting to use, for example, IP information from one group of servers
# without having to talk to them in the same playbook run to get their
# current IP information.
fact_caching = memory


# retry files
retry_files_enabled = False
#retry_files_save_path = ~/.ansible-retry

[ssh_connection]
ssh_args = -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -i {{ .SshKey }}
`); err != nil {
		panic(err)
	}
}

func Run(depl *models.Deployment) error {

	// Get Nodes
	nodes, err := models.NodeMapper.FetchAllForApp(depl.Application)
	if err != nil {
		log.Println("ressources/ansible: Run: NodeMapper.FetchAllForApp:", err)
		return err
	}

	// Temp work dir
	tmpPath, err := ioutil.TempDir("", depl.Id.Hex())
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(path.Join(tmpPath, "karhu.log"), []byte(fmt.Sprintf("Deployment %s\n", depl.Id.Hex())), 0644); err != nil {
		return err
	}

	log.Println("ressources/ansible: Run: tmpPath:", tmpPath)
	depl.TmpPath = tmpPath
	if err := models.DeploymentMapper.Update(depl); err != nil {
		return err
	}
	defer os.RemoveAll(tmpPath)

	// build hosts
	if err := buildHosts(tmpPath, depl.Build.RuntimeCfg, nodes); err != nil {
		return err
	}

	// build playbook
	roles, err := buildPlaybook(tmpPath, depl.Build.RuntimeCfg, depl.Application)
	if err != nil {
		return err
	}

	if err := buildVars(tmpPath, depl.Build.RuntimeCfg, depl.Application); err != nil {
		return err
	}

	if err := buildConfig(tmpPath); err != nil {
		return err
	}

	// Copy required roles
	if err := copyRoles(tmpPath, roles); err != nil {
		return err
	}

	// Extract build zip
	if err := extractArchive(tmpPath, depl.Build); err != nil {
		return err
	}

	// run playbook
	errPlaybook := runPlaybook(tmpPath)

	data, err := ioutil.ReadFile(path.Join(tmpPath, "karhu.log"))
	if err != nil {
		return err
	}

	depl.Logs = string(data)

	return errPlaybook

}

func buildHosts(tmpPath string, runtimeCfg *models.RuntimeConfiguration, nodes models.Nodes) error {

	w, err := os.OpenFile(path.Join(tmpPath, HOSTS_FILENAME), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer w.Close()

	return hostFileTemplate.Execute(w, map[string]interface{}{
		"RuntimeConfig": runtimeCfg,
		"Nodes":         nodes,
	})
}

func buildPlaybook(tmpPath string, runtimeCfg *models.RuntimeConfiguration, app *models.Application) ([]string, error) {

	w, err := os.OpenFile(path.Join(tmpPath, PLAYBOOK_FILENAME), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer w.Close()

	var roles []string
	if len(app.Deps) > 0 {
		roles = append(roles, models.APPLICATION_TYPE_SERVICE)
	}

	// Because service SHOULD ALWAYS BE FIRST, because of logic
	roles = append(roles, runtimeCfg.Type)

	return roles, playbookFileTemplate.Execute(w, map[string]interface{}{
		"RuntimeConfig": runtimeCfg,
		"Vars":          VARS_FILENAME,
		"Roles":         roles,
	})
}

func buildVars(tmpPath string, runtimeCfg *models.RuntimeConfiguration, app *models.Application) error {

	w, err := os.OpenFile(path.Join(tmpPath, VARS_FILENAME), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer w.Close()

	configs, err := extractConfigs(tmpPath, runtimeCfg, app)
	if err != nil {
		return err
	}

	var services []*models.Build
	var servicesConfigs []ConfigFile
	for _, dep := range app.Deps {
		if dep.Type != models.APPLICATION_TYPE_SERVICE {
			continue
		}

		build, err := models.BuildMapper.FetchLast(dep)
		if err != nil {
			return err
		}

		if build == nil {
			log.Println("No last build for app:", app.Name, app.Id.Hex())
			continue
		}

		services = append(services, build)

		cfgs, err := extractConfigs(tmpPath, build.RuntimeCfg, dep)
		if err != nil {
			return err
		}

		servicesConfigs = append(servicesConfigs, cfgs...)
	}

	return varsFileTemplate.Execute(w, map[string]interface{}{
		"TmpPath":         tmpPath,
		"RuntimeConfig":   runtimeCfg,
		"Application":     app,
		"Configs":         configs,
		"Services":        services,
		"ServicesConfigs": servicesConfigs,
	})
}

func buildConfig(tmpPath string) error {

	w, err := os.OpenFile(path.Join(tmpPath, CONFIG_FILENAME), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer w.Close()

	return configFileTemplate.Execute(w, map[string]interface{}{
		"SshKey": ssh.PrivateKeyPath(),
	})
}

// Extract build zip into {{ tmpPath }} /karhu
func extractArchive(tmpPath string, build *models.Build) error {

	// Fetch zip
	data, err := file.Get(build.FilePath)
	if err != nil {
		return err
	}

	// Write zip file
	zipPath := path.Join(tmpPath, "karhu.zip")
	if err := ioutil.WriteFile(zipPath, data, 0644); err != nil {
		return err
	}

	destDir := path.Join(tmpPath, "karhu")
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	// Unzip file
	if err := file.Unzip(zipPath, destDir); err != nil {
		return err
	}

	return nil
}

type ConfigFile struct {
	Src     string
	Dest    string
	Mode    string
	DestDir string
	Notify  models.ConfigNotify // Linked service
}

// Copy application configs
func extractConfigs(tmpPath string, runtimeCfg *models.RuntimeConfiguration, app *models.Application) ([]ConfigFile, error) {

	// Get all configs
	configs, err := models.ConfigMapper.FetchAllEnabled(app)
	if err != nil {
		return nil, err
	}

	if len(configs) == 0 {
		return nil, nil
	}

	// create dest directory
	destDir := path.Join(tmpPath, "configs")
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return nil, err
	}

	var configFiles []ConfigFile

	for _, config := range configs {

		src := path.Join(destDir, fmt.Sprintf("%x", sha1.Sum([]byte(config.Path))))
		if err := ioutil.WriteFile(src, []byte(config.Content), 0644); err != nil {
			return nil, err
		}

		destPath := config.Path
		if destPath[0] != '/' {
			destPath = path.Join(runtimeCfg.Workdir, config.Path)
		}

		configFiles = append(configFiles, ConfigFile{
			Src:     path.Clean(src),      // Absolute path to src file
			Dest:    path.Clean(destPath), // absolute path
			DestDir: path.Dir(path.Clean(destPath)),
			Mode:    "0644",
			Notify:  config.Notify,
		})
	}

	return configFiles, nil
}

// Copy required roles
func copyRoles(tmpPath string, roles []string) error {

	if err := os.MkdirAll(path.Join(tmpPath, "roles"), 0755); err != nil {
		return err
	}

	for _, role := range roles {
		command := fmt.Sprintf("cp -rf %s %s", path.Join(env.GetDefault("ANSIBLE_ROLES_DIR", "ansible"), role), path.Join(tmpPath, "roles"))

		cmd := exec.Command("sh", "-c", command)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

func runPlaybook(tmpPath string) error {

	command := fmt.Sprintf("ansible-playbook -i %s %s > %s/karhu.log 2>&1", HOSTS_FILENAME, PLAYBOOK_FILENAME, tmpPath)
	log.Println("ressources/ansible: runPlaybook:", command)
	cmd := exec.Command("sh", "-c", fmt.Sprintf("cd %s && %s", tmpPath, command))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

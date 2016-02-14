package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gotoolz/env"
	"github.com/gotoolz/errors"
	"github.com/karhuteam/karhu/models"
	"github.com/karhuteam/karhu/ressources/logstash"
	"github.com/karhuteam/karhu/ressources/ssh"
	"log"
	"net/http"
)

type NodeController struct {
}

func NewNodeController(s *gin.RouterGroup) *NodeController {

	ctl := &NodeController{}

	s.GET("/nodes/register.sh", ctl.getRegisterSH)
	s.POST("/nodes", ctl.postNode)
	s.GET("/nodes", ctl.getAllNode)
	s.GET("/nodes/config/:filename", ctl.getNodeConfig)

	return ctl
}

func (pc *NodeController) getRegisterSH(c *gin.Context) {

	publicKey, err := ssh.GetPublicKey()
	if err != nil {
		panic(err)
	}

	karhuHost := env.GetDefault("PUBLIC_HOST", "http://127.0.0.1:8080")
	logstashIP := env.GetDefault("LOGSTASH_IP", "127.0.0.1")
	collectdUser, collectdPassword, err := logstash.ReadAuthfile()
	if err != nil {
		panic(err)
	}

	log.Println("URL:", c.Request.URL)

	clientIP := c.ClientIP()

	basicAuth := ""
	if auth := env.Get("BASIC_AUTH"); auth != "" {
		basicAuth = `"-u ` + auth + `"`
	}

	c.String(http.StatusOK, fmt.Sprintf(`
#!/bin/bash
echo "Registering host on Karhu..."; echo
KARHU_HOST=%s
LOGSTASH_IP=%s
PUBLIC_KEY='%s'
LOGSTASH_CRT_URL=$KARHU_HOST/api/nodes/config/logstash.crt
LOGSTASH_CRT_PATH=/etc/filebeat/certs/logstash.crt
FILEBEAT_CONFIG_URL=$KARHU_HOST/api/nodes/config/filebeat.yml
FILEBEAT_CONFIG_PATH=/etc/filebeat/filebeat.yml
AUTHORIZED_KEYS_DIR=%s
AUTHORIZED_KEYS_FILE=%s
CLIENT_IP=%s
SSH_PORT=%s
SSH_USER=$USER
BASIC_AUTH=%s
SETUP_MONITORING=%s
INFLUXDB_COLLECTD_HOST=%s
INFLUXDB_COLLECTD_PORT=%s
COLLECTD_USERNAME=%s
COLLECTD_PASSWORD=%s
COLLECTD_CONFIG_PATH=/etc/collectd/collectd.conf.d/karhu.conf
NO_REGISTER=%s

SUDO=
if [ "$SSH_USER" != "root" ]; then
	echo "Check sudo..."
	sudo -n true
	if [ "$?" != "0" ]; then
		echo "You need root access or sudo without password..."
		exit 1
	fi
	SUDO=sudo
fi

if [ "$NO_REGISTER" != "1" ]; then

	if [ ! -d "$AUTHORIZED_KEYS_DIR" ]; then
		mkdir -p $AUTHORIZED_KEYS_DIR || exit 1
	fi

	if [ ! -f "$AUTHORIZED_KEYS_FILE" ]; then
		touch $AUTHORIZED_KEYS_FILE || exit 1
	fi

	echo "Setting up ssh keys..."
	grep -q -F "$(echo $PUBLIC_KEY)" $AUTHORIZED_KEYS_FILE || echo $PUBLIC_KEY >> $AUTHORIZED_KEYS_FILE

	echo "Registering node..."
	curl --fail $BASIC_AUTH -X POST $KARHU_HOST/api/nodes -d hostname=$(hostname) -d ip=$CLIENT_IP -d ssh_port=$SSH_PORT -d ssh_user=$SSH_USER || exit 1
	echo
fi
if [ "$SETUP_MONITORING" = "1" ]; then

	echo "Setup logstash host"
	$SUDO sed '/ karhu$/{h;s/.*/'$LOGSTASH_IP' karhu/};${x;/^$/{s//'$LOGSTASH_IP' karhu/;H};x}' -i /etc/hosts

	echo "Setup monitoring..."
	if [ ! -f "$(which collectd)" ]; then
		$SUDO apt-get update || exit 1
		$SUDO apt-get install -y --no-install-recommends collectd || exit 1
	fi
	echo "LoadPlugin network
<Plugin "network">
    <Server \"karhu\" \"$INFLUXDB_COLLECTD_PORT\">
		SecurityLevel "Encrypt"
		Username "$COLLECTD_USERNAME"
		Password "$COLLECTD_PASSWORD"
	</Server>
</Plugin>" | $SUDO tee $COLLECTD_CONFIG_PATH || exit 1
	echo "Restard collectd"
	$SUDO service collectd restart || exit 1

	# Setup filebeat
	if [ ! -f "$(which filebeat)" ]; then
		echo "deb https://packages.elastic.co/beats/apt stable main" |  $SUDO tee /etc/apt/sources.list.d/beats.list || exit 1
		curl -L https://packages.elastic.co/GPG-KEY-elasticsearch | $SUDO apt-key add - || exit 1
		$SUDO apt-get update || exit 1
		$SUDO apt-get install -y filebeat || exit 1
	fi

	if [ ! -d "$(dirname $LOGSTASH_CRT_PATH)" ]; then
		$SUDO mkdir -p $(dirname $LOGSTASH_CRT_PATH) || exit 1
	fi

	# setup crt
	$SUDO curl $BASIC_AUTH -o $LOGSTASH_CRT_PATH $LOGSTASH_CRT_URL || exit 1
	# setup config
	$SUDO curl $BASIC_AUTH -o $FILEBEAT_CONFIG_PATH $FILEBEAT_CONFIG_URL || exit 1

	$SUDO service filebeat restart || exit 1
fi
echo "Done."`, karhuHost, logstashIP, publicKey, ssh.SSH_AUTHORIZED_KEYS_DIR, ssh.AuthorizedKeysPath(), clientIP, c.DefaultQuery("ssh_port", "22") /*, c.DefaultQuery("ssh_user", "root") */, basicAuth, c.DefaultQuery("monit", "1"), env.Get("INFLUXDB_COLLECTD_HOST"), env.Get("INFLUXDB_COLLECTD_PORT"), collectdUser, collectdPassword, c.DefaultQuery("noreg", "0")))
}

func (pc *NodeController) postNode(c *gin.Context) {

	var form models.NodeCreateForm
	if err := c.Bind(&form); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := form.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	node := models.NodeMapper.Create(&form)

	if err := models.NodeMapper.CheckSsh(node); err != nil {
		c.JSON(http.StatusBadRequest, errors.Error{
			Label: "bad_ssh_configuration",
			Field: "ip,ssh_port,ssh_user",
			Text:  err.Error(),
		})
		return
	}

	if err := models.NodeMapper.Save(node); err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, node)
}

func (pc *NodeController) getAllNode(c *gin.Context) {

	nodes, err := models.NodeMapper.FetchAll()
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, nodes)
}

func (pc *NodeController) getNodeConfig(c *gin.Context) {

	filename := c.Param("filename")

	var data []byte
	var err error

	switch filename {
	case "filebeat.yml":
		data, err = logstash.GetFilebeatConfig()
	case "logstash.crt":
		data, err = logstash.GetCert()
	}
	if err != nil {
		panic(err)
	}

	if data != nil {
		c.Data(http.StatusOK, "text/plain", data)
		return
	}

	c.AbortWithStatus(http.StatusNotFound)
}

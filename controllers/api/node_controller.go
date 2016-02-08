package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gotoolz/env"
	"github.com/gotoolz/errors"
	"github.com/karhuteam/karhu/models"
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

	return ctl
}

func (pc *NodeController) getRegisterSH(c *gin.Context) {

	publicKey, err := ssh.GetPublicKey()
	if err != nil {
		panic(err)
	}
	karhuHost := env.GetDefault("PUBLIC_HOST", "http://127.0.0.1:8080")

	log.Println("URL:", c.Request.URL)

	clientIP := c.ClientIP()

	basicAuth := ""
	if auth := env.Get("BASIC_AUTH"); auth != "" {
		basicAuth = "-u " + auth + " "
	}

	c.String(http.StatusOK, fmt.Sprintf(`
#!/bin/bash
echo "Registering host on Karhu..."; echo
PUBLIC_KEY='%s'
AUTHORIZED_KEYS_DIR=%s
AUTHORIZED_KEYS_FILE=%s
KARHU_HOST=%s
CLIENT_IP=%s
SSH_PORT=%s
SSH_USER=$USER
BASIC_AUTH="%s"
SETUP_MONITORING=%s
INFLUXDB_COLLECTD_HOST=%s
INFLUXDB_COLLECTD_PORT=%s
COLLECTD_CONFIG_PATH=/etc/collectd/collectd.conf.d/karhu.conf

SUDO=
if [ "$SSH_USER" != "root" ]; then
	echo "Check sudo..."
	sudo -n true || $(echo "You need root access or sudo without password..." && exit 1)
	SUDO=sudo
fi

if [ ! -d "$AUTHORIZED_KEYS_DIR" ]; then
	mkdir -p $AUTHORIZED_KEYS_DIR || exit 1
fi

if [ ! -f "$AUTHORIZED_KEYS_FILE" ]; then
	touch $AUTHORIZED_KEYS_FILE || exit 1
fi

echo "Setting up ssh keys..."
grep -q -F "$(echo $PUBLIC_KEY)" $AUTHORIZED_KEYS_FILE || echo $PUBLIC_KEY >> $AUTHORIZED_KEYS_FILE

echo "Registering node..."
curl --fail $BASIC_AUTH-X POST $KARHU_HOST/api/nodes -d hostname=$(hostname) -d ip=$CLIENT_IP -d ssh_port=$SSH_PORT -d ssh_user=$SSH_USER || exit 1
echo
if [ "$SETUP_MONITORING" = "1" ]; then
	echo "Setup monitoring..."
	$SUDO apt-get update || exit 1
	$SUDO apt-get install -y collectd || exit 1
	echo "LoadPlugin network
<Plugin "network">
    Server \"$INFLUXDB_COLLECTD_HOST\" \"$INFLUXDB_COLLECTD_PORT\"
</Plugin>" | $SUDO tee $COLLECTD_CONFIG_PATH || exit 1
	echo "Restard collectd"
	$SUDO service collectd restart || exit 1
fi
echo "Done."`, publicKey, ssh.SSH_AUTHORIZED_KEYS_DIR, ssh.AuthorizedKeysPath(), karhuHost, clientIP, c.DefaultQuery("ssh_port", "22") /*, c.DefaultQuery("ssh_user", "root") */, basicAuth, c.DefaultQuery("monit", "1"), env.Get("INFLUXDB_COLLECTD_HOST"), env.Get("INFLUXDB_COLLECTD_PORT")))
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

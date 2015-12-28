package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gotoolz/env"
	"github.com/gotoolz/errors"
	"github.com/karhuteam/karhu/models"
	"github.com/karhuteam/karhu/ressources/ssh"
	"log"
	"net"
	"net/http"
	"strings"
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

	clientIP, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		clientIP = c.Request.RemoteAddr
	}
	clientIP = strings.TrimSpace(clientIP)

	c.String(http.StatusOK, fmt.Sprintf(`
#!/bin/bash
echo "Registering host on Karhu..."; echo
PUBLIC_KEY='%s'
AUTHORIZED_KEYS_DIR=%s
AUTHORIZED_KEYS_FILE=%s
KARHU_HOST=%s
CLIENT_IP=%s

if [ ! -d "$AUTHORIZED_KEYS_DIR" ]; then
	mkdir -p $AUTHORIZED_KEYS_DIR || exit 1
fi

if [ ! -f "$AUTHORIZED_KEYS_FILE" ]; then
	touch $AUTHORIZED_KEYS_FILE || exit 1
fi

echo "Setting up ssh keys..."
grep -q -F "$(echo $PUBLIC_KEY)" $AUTHORIZED_KEYS_FILE || echo $PUBLIC_KEY >> $AUTHORIZED_KEYS_FILE

echo "Registering node..."
curl -X POST $KARHU_HOST/api/nodes -d hostname=$(hostname) -d ip=$CLIENT_IP -d ssh_port=%s -d ssh_user=%s; echo
echo "Done."`, publicKey, ssh.SSH_AUTHORIZED_KEYS_DIR, ssh.AuthorizedKeysPath(), karhuHost, clientIP, c.DefaultQuery("ssh_port", "22"), c.DefaultQuery("ssh_user", "root")))
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

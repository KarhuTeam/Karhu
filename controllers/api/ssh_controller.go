package api

import (
	"github.com/gin-gonic/gin"
	// "github.com/karhuteam/karhu/models"
	"github.com/karhuteam/karhu/ressources/ssh"
	"net/http"
)

type SshController struct {
}

func NewSshController(s *gin.RouterGroup) *SshController {

	ctl := &SshController{}

	s.GET("/ssh/public-key", ctl.getPublicKey)

	return ctl
}

func (sc *SshController) getPublicKey(c *gin.Context) {

	data, err := ssh.GetPublicKey()
	if err != nil {
		panic(err)
	}

	c.String(http.StatusOK, string(data))

}

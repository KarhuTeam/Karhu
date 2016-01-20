package front

import (
	"github.com/gin-gonic/gin"
	"github.com/karhuteam/karhu/web"
	"net/http"
)

type ConfigurationController struct {
}

func NewConfigurationController(s *web.Server) *ConfigurationController {

	ctl := &ConfigurationController{}

	// 1 - Add a configuration
	s.POST("/application/configuration/:app_id", ctl.postAddConfigurationAction)
	s.GET("/application/configuration/:app_id", ctl.getAddConfigurationAction)

	return ctl
}

/**
 * 1 - add a configuration
 */
func (ctl *ConfigurationController) postAddConfigurationAction(c *gin.Context) {

	c.HTML(http.StatusOK, "configuration_add.html", nil)
}

func (ctl *ConfigurationController) getAddConfigurationAction(c *gin.Context) {

	c.HTML(http.StatusOK, "configuration_add.html", nil)
}

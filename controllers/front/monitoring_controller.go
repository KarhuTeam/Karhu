package front

import (
	// "fmt"
	"github.com/gin-gonic/gin"
	// "github.com/gotoolz/errors"
	"github.com/karhuteam/karhu/models"
	"github.com/karhuteam/karhu/web"
	"net/http"
	// "strconv"
)

type MonitoringController struct {
}

func NewMonitoringController(s *web.Server) *MonitoringController {
	ctl := &MonitoringController{}

	s.GET("/monitoring", ctl.getMonitoringAction)

	return ctl
}

func (pc *MonitoringController) getMonitoringAction(c *gin.Context) {

	monitoring_data := models.GetDefaultMonitoring()

	c.HTML(http.StatusOK, "monitoring.html", map[string]interface{} {
		"Data":	monitoring_data.Name,
	})
}
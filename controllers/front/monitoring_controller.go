package front

import (
	"fmt"
	"github.com/gin-gonic/gin"
	// "github.com/gotoolz/errors"
	// "github.com/karhuteam/karhu/models"
	"github.com/karhuteam/karhu/web"
	"net/http"
	// "strconv"
)

type MonitoringController struct {
}

func NewMonitoringController(s *web.Server) *MonitoringController {
	ctl := &MonitoringController{}
	fmt.Printf("HELLO THERE")
	s.GET("/monitoring", ctl.getMonitoringAction)

	return ctl
}

func (pc *MonitoringController) getMonitoringAction(c *gin.Context) {
	fmt.Printf("HELLO THERE2")
	c.HTML(http.StatusOK, "monitoring.html", map[string]interface{} {

	})
}
package front

import (
	"github.com/gin-gonic/gin"
	// "github.com/gotoolz/env"
	"github.com/karhuteam/karhu/models"
	"github.com/karhuteam/karhu/web"
	"net/http"
)

type MonitoringController struct {
}

func NewMonitoringController(s *web.Server) *MonitoringController {
	ctl := &MonitoringController{}

	s.GET("/monitoring", ctl.getMonitoringAction)

	return ctl
}

func (pc *MonitoringController) getMonitoringAction(c *gin.Context) {

	host := c.Query("host")
	stat := c.DefaultQuery("stat", "all")
	t := c.DefaultQuery("time", "last1800")

	hosts, err := models.NodeMapper.FetchAll()
	if err != nil {
		panic(err)
	}

	var graphs models.Graphs

	var target []string
	if host != "" {
		target = []string{host}
	} else {
		for _, h := range hosts {
			target = append(target, h.Hostname)
		}
	}

	if graphs, err = models.GraphMapper.FetchAll(target, stat, t); err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "monitoring_show.html", map[string]interface{}{
		"host":   host,
		"stat":   stat,
		"time":   t,
		"hosts":  hosts,
		"stats":  models.GraphStats,
		"graphs": graphs,
	})
}

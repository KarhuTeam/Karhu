package front

import (
	"github.com/gin-gonic/gin"
	"github.com/gotoolz/env"
	"github.com/karhuteam/karhu/web"
	"net/http"
	"github.com/karhuteam/karhu/models"
)

type MonitoringController struct {
}

func NewMonitoringController(s *web.Server) *MonitoringController {
	ctl := &MonitoringController{}

	s.GET("/monitoring", ctl.getMonitoringAction)

	return ctl
}

func (pc *MonitoringController) getMonitoringAction(c *gin.Context) {

	grafana_url := env.GetDefault("GRAFANA_URL", "http://localhost:3000")
	grafana_url = grafana_url + "/dashboard/script/karhu.js"

	hosts := c.Request.URL.Query()["hosts"]
	page_title := "General"

	if len(hosts) > 0 {
		grafana_url = grafana_url + "?hosts=" + hosts[0]
		page_title = hosts[0]
	} else {
		nodes, err := models.NodeMapper.FetchAll()
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error_500.html", map[string]interface{}{
				"error": err,
			})
			return
		}
		h := ""
		for _, n := range nodes {
			if h == "" {
				h = n.Hostname	
			} else {
				h = h + "," + n.Hostname
			}
		}
		grafana_url = grafana_url + "?hosts=\"" + h + "\""
	}

	c.HTML(http.StatusOK, "monitoring_view.html", map[string]interface{}{
		"grafana_url": grafana_url,
		"page_title": page_title,
	})
}

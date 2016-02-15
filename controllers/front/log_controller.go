package front

import (
	"github.com/gin-gonic/gin"
	// "github.com/gotoolz/env"
	"github.com/karhuteam/karhu/models"
	"github.com/karhuteam/karhu/web"
	"net/http"
	"strconv"
)

type LogController struct {
}

func NewLogController(s *web.Server) *LogController {

	ctl := &LogController{}

	// 1 - Show log page
	s.GET("/logs", ctl.getLogsAction)

	return ctl
}

/**
 * 1 - Show log page
 */
func (ctl *LogController) getLogsAction(c *gin.Context) {

	tags := c.Request.URL.Query()["tags[]"]

	app := c.DefaultQuery("app", "")

	s := c.DefaultQuery("count", "30")
	count, _ := strconv.Atoi(s)

	query := c.DefaultQuery("query", "*")
	var logs models.Logs
	if query != "" {
		var err error
		logs, err = models.LogMapper.Search(query, tags, app, count)
		if err != nil {
			panic(err)
		}
	}

	applications, err := models.ApplicationMapper.FetchAll()
	if err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "log_show.html", map[string]interface{}{
		"query":        query,
		"result":       logs,
		"tags":         tags,
		"app":          app,
		"count":        s,
		"applications": applications,
	})
}

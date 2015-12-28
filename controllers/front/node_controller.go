package front

import (
	"github.com/gin-gonic/gin"
	"github.com/karhuteam/karhu/models"
	"github.com/karhuteam/karhu/web"
	"net/http"
)

type NodeController struct {
}

func NewNodeController(s *web.Server) *NodeController {

	ctl := &NodeController{}

	s.GET("/nodes", ctl.getNodesAction)
	s.GET("/node/edit/:id", ctl.getNodeAction)
	s.POST("/node/edit/:id", ctl.postNodeAction)

	return ctl
}

func (pc *NodeController) getNodesAction(c *gin.Context) {

	nodes, err := models.NodeMapper.FetchAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error_500.html", map[string]interface{}{
			"error": err,
		})
		return
	}
	c.HTML(http.StatusOK, "node_list.html", map[string]interface{}{
		"nodes": nodes,
	})
}

func (pc *NodeController) postNodeAction(c *gin.Context) {

	id := c.Param("id")

	var form models.NodeUpdateForm
	if err := c.Bind(&form); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	node, err := models.NodeMapper.FetchOneById(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error_500.html", map[string]interface{}{
			"error": err,
		})
		return
	}
	node.Tags = form.Tags
	node.Description = form.Description
	models.NodeMapper.Update(node)

	c.Redirect(http.StatusMovedPermanently, "/nodes")
}

func (pc *NodeController) getNodeAction(c *gin.Context) {

	id := c.Param("id")

	node, err := models.NodeMapper.FetchOneById(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error_500.html", map[string]interface{}{
			"error": err,
		})
		return
	}
	c.HTML(http.StatusOK, "node_edit.html", map[string]interface{}{
		"errors": nil,
		"node":   node,
	})
}
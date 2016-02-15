package front

import (
	"github.com/gin-gonic/gin"
	"github.com/gotoolz/env"
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
	s.GET("/node/add", ctl.getNodeAddAction)
	s.GET("/node/add/ec2", ctl.getNodeAddEc2Action)
	s.POST("/node/delete/:id", ctl.postDeleteNodeAction)

	return ctl
}

func (pc *NodeController) getNodeAddAction(c *gin.Context) {

	basicAuth := env.Get("BASIC_AUTH")
	if auth := env.Get("BASIC_AUTH"); auth != "" {
		basicAuth = "-u " + auth + " "
	}

	c.HTML(http.StatusOK, "node_add.html", map[string]interface{}{
		"PublicHost": c.DefaultQuery("karhu_url", env.Get("PUBLIC_HOST")),
		"SshUser":    c.DefaultQuery("ssh_user", "root"),
		"SshPort":    c.DefaultQuery("ssh_port", "22"),
		"Monit":      c.DefaultQuery("monit", "1"),
		"BasicAuth":  basicAuth,
	})
}

func (pc *NodeController) getNodeAddEc2Action(c *gin.Context) {

	basicAuth := env.Get("BASIC_AUTH")
	if auth := env.Get("BASIC_AUTH"); auth != "" {
		basicAuth = "-u " + auth + " "
	}

	c.HTML(http.StatusOK, "node_add_ec2.html", map[string]interface{}{
		"PublicHost": c.DefaultQuery("karhu_url", env.Get("PUBLIC_HOST")),
		"SshUser":    c.DefaultQuery("ssh_user", "root"),
		"SshPort":    c.DefaultQuery("ssh_port", "22"),
		"Monit":      c.DefaultQuery("monit", "1"),
		"BasicAuth":  basicAuth,
	})
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

	if node == nil {
		c.HTML(http.StatusNotFound, "error_404.html", map[string]interface{}{
			"text": "Node not found",
		})
		return
	}

	node.Tags = form.Tags
	node.Description = form.Description
	models.NodeMapper.Update(node)

	c.Redirect(http.StatusFound, "/nodes")
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
	if node == nil {
		c.HTML(http.StatusNotFound, "error_404.html", map[string]interface{}{
			"text": "Node not found",
		})
		return
	}

	c.HTML(http.StatusOK, "node_edit.html", map[string]interface{}{
		"errors": nil,
		"node":   node,
	})
}

func (pc *NodeController) postDeleteNodeAction(c *gin.Context) {

	id := c.Param("id")

	node, err := models.NodeMapper.FetchOneById(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error_500.html", map[string]interface{}{
			"error": err,
		})
		return
	}

	if node == nil {
		c.HTML(http.StatusNotFound, "error_404.html", map[string]interface{}{
			"text": "Node not found",
		})
		return
	}

	if err := models.NodeMapper.Delete(node); err != nil {
		c.HTML(http.StatusInternalServerError, "error_500.html", map[string]interface{}{
			"error": err,
		})
		return
	}

	c.Redirect(http.StatusFound, "/nodes")
}

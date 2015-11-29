package api

import (
	"github.com/gin-gonic/gin"
	"github.com/karhuteam/karhu/models"
	"net/http"
)

type ProjectController struct {
}

func NewProjectController(s *gin.RouterGroup) *ProjectController {

	ctl := &ProjectController{}

	s.POST("/projects", ctl.postProject)
	s.GET("/projects", ctl.getProjectList)
	s.GET("/projects/:id", ctl.getProject)
	s.PUT("/projects/:id", ctl.putProject)
	s.DELETE("/projects/:id", ctl.deleteProject)

	return ctl
}

func (pc *ProjectController) postProject(c *gin.Context) {

	var form models.ProjectCreateForm
	if err := c.Bind(&form); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := form.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	project := models.ProjectMapper.Create(&form)

	if err := models.ProjectMapper.Save(project); err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, project)
}

func (pc *ProjectController) getProjectList(c *gin.Context) {

	projects, err := models.ProjectMapper.FetchAll()
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, projects)
}

func (pc *ProjectController) getProject(c *gin.Context) {

	id := c.Param("id")

	project, err := models.ProjectMapper.FetchOne(id)
	if err != nil {
		panic(err)
	}

	if project == nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, project)
}

func (pc *ProjectController) putProject(c *gin.Context) {

	var form models.ProjectUpdateForm
	if err := c.Bind(&form); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := form.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	id := c.Param("id")

	project, err := models.ProjectMapper.FetchOne(id)
	if err != nil {
		panic(err)
	}

	if project == nil {
		c.JSON(http.StatusNotFound, project)
		return
	}

	project.Update(&form)

	if err := models.ProjectMapper.Update(project); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, project)
}

func (pc *ProjectController) deleteProject(c *gin.Context) {

	id := c.Param("id")

	project, err := models.ProjectMapper.FetchOne(id)
	if err != nil {
		panic(err)
	}

	if err := models.ProjectMapper.Delete(project); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, nil)
}

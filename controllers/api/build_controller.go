package api

import (
	"github.com/gin-gonic/gin"
	"github.com/karhuteam/karhu/models"
	"github.com/wayt/goerrors"
	"net/http"
)

type BuildController struct {
}

func NewBuildController(s *gin.RouterGroup) *BuildController {

	ctl := &BuildController{}

	s.POST("/projects/:id/builds", ctl.postBuild)
	s.GET("/projects/:id/builds", ctl.getBuildList)

	return ctl
}

func (pc *BuildController) postBuild(c *gin.Context) {

	var form models.BuildCreateForm
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
		c.JSON(http.StatusNotFound, goerrors.New(goerrors.Error{
			Label: "invalid_project",
			Field: "id",
			Text:  "Invalid project ID in URL",
		}))
		return
	}

	build := models.BuildMapper.Create(project, &form)

	if err := models.BuildMapper.Save(build); err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, build)
}

func (pc *BuildController) getBuildList(c *gin.Context) {

	id := c.Param("id")

	project, err := models.ProjectMapper.FetchOne(id)
	if err != nil {
		panic(err)
	}

	if project == nil {
		c.JSON(http.StatusNotFound, goerrors.New(goerrors.Error{
			Label: "invalid_project",
			Field: "id",
			Text:  "Invalid project ID in URL",
		}))
		return
	}

	builds, err := models.BuildMapper.FetchAll(project)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, builds)
}

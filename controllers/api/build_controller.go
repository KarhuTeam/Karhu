package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gotoolz/errors"
	"github.com/karhuteam/karhu/models"
	"net/http"
)

type BuildController struct {
}

func NewBuildController(s *gin.RouterGroup) *BuildController {

	ctl := &BuildController{}

	s.POST("/apps/:id/builds", ctl.postBuild)
	s.GET("/apps/:id/builds", ctl.getBuildList)

	return ctl
}

func (pc *BuildController) getApp(c *gin.Context) (*models.Application, error) {

	id := c.Param("id")

	app, err := models.ApplicationMapper.FetchOne(id)
	if err != nil {
		panic(err)
	}

	if app == nil {
		return nil, errors.New(errors.Error{
			Label: "invalid_application",
			Field: "id",
			Text:  "Invalid application ID in URL",
		})
	}

	return app, nil
}

func (pc *BuildController) postBuild(c *gin.Context) {

	app, err := pc.getApp(c)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	commitHash := c.DefaultPostForm("commit_hash", "")

	build := models.BuildMapper.Create(app, commitHash)

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.New(errors.Error{
			Label: "invalid_file",
			Field: "file",
			Text:  "Missing zip file",
		}))
		return
	}
	defer file.Close()

	if err := build.AttachFile(file); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if err := models.BuildMapper.Save(build); err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, build)
}

func (pc *BuildController) getBuildList(c *gin.Context) {

	app, err := pc.getApp(c)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	builds, err := models.BuildMapper.FetchAll(app)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, builds)
}

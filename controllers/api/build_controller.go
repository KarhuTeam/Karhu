package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gotoolz/errors"
	"github.com/karhuteam/karhu/models"
	"github.com/karhuteam/karhu/ressources/ansible"
	"log"
	"net/http"
	"runtime"
	"time"
)

type BuildController struct {
}

func NewBuildController(s *gin.RouterGroup) *BuildController {

	ctl := &BuildController{}

	s.POST("/apps/:id/builds", ctl.postBuild)
	s.GET("/apps/:id/builds", ctl.getBuildList)
	s.POST("/apps/:id/builds/:build_id/deploy", ctl.postBuildDeploy)

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
		c.JSON(http.StatusBadRequest, errors.New(errors.Error{
			Label: "invalid_file",
			Field: "file",
			Text:  err.Error(),
		}))
		return
	}

	if err := models.BuildMapper.Save(build); err != nil {
		panic(err)
	}

	if c.DefaultQuery("deploy", "0") == "1" {

		resp, err := http.Post(fmt.Sprintf("http://localhost:8080/api/apps/%s/builds/%s/deploy", app.Id.Hex(), build.Id.Hex()), "json/application", nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.New(errors.Error{
				Label: "deploy_issue",
				Field: "deploy",
				Text:  err.Error(),
			}))
			return
		}
		resp.Body.Close()
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

func (pc *BuildController) postBuildDeploy(c *gin.Context) {

	buildId := c.Param("build_id")

	app, err := pc.getApp(c)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	build, err := models.BuildMapper.FetchOne(app, buildId)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	if build == nil {
		c.JSON(http.StatusNotFound, errors.New(errors.Error{
			Label: "invalid_build",
			Field: "build_id",
			Text:  "Invalid build ID in URL",
		}))
		return
	}

	depl := models.DeploymentMapper.Create(app, build)

	if err := models.DeploymentMapper.Save(depl); err != nil {
		panic(err)
	}

	go func() {

		start := time.Now()

		// catch panic
		defer func() {
			if err := recover(); err != nil {
				log.Println("ansible:", err)

				trace := make([]byte, 2048)
				runtime.Stack(trace, true)
				log.Println(string(trace))

				depl.Status = models.STATUS_ERROR
				depl.Duration = time.Since(start)
				if err := models.DeploymentMapper.Update(depl); err != nil {
					log.Println(err)
				}
			}
		}()

		depl.Status = models.STATUS_RUNNING
		if err := models.DeploymentMapper.Update(depl); err != nil {
			return
		}

		if err := ansible.Run(depl); err != nil {
			panic(err)
		}

		depl.Status = models.STATUS_DONE
		depl.Duration = time.Since(start)
		if err := models.DeploymentMapper.Update(depl); err != nil {
			return
		}
	}()

	c.JSON(http.StatusOK, depl)
}

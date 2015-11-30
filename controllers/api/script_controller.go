package api

import (
	"github.com/gin-gonic/gin"
	"github.com/karhuteam/karhu/models"
	"github.com/wayt/goerrors"
	"net/http"
	"strings"
)

type ScriptController struct {
}

func NewScriptController(s *gin.RouterGroup) *ScriptController {

	ctl := &ScriptController{}

	s.POST("/projects/:id/scripts", ctl.postScript)
	s.PUT("/projects/:id/scripts/:script_id", ctl.putScript)
	s.GET("/projects/:id/scripts", ctl.getScriptList)
	s.GET("/projects/:id/scripts/:script_id", ctl.getScript)
	s.DELETE("/projects/:id/scripts/:script_id", ctl.deleteScript)

	return ctl
}

func (sc *ScriptController) postScript(c *gin.Context) {

	var form models.ScriptCreateForm
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

	script, err := models.ScriptMapper.FetchOneByName(project, form.Name)
	if err != nil {
		panic(err)
	}

	if script != nil {
		c.JSON(http.StatusConflict, goerrors.New(goerrors.Error{
			Label: "invalid_name",
			Field: "name",
			Text:  "Duplicate script name for this project",
		}))
		return
	}

	script = models.ScriptMapper.Create(project, &form)

	if err := models.ScriptMapper.Save(script); err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, script)
}

func (sc *ScriptController) putScript(c *gin.Context) {

	var form models.ScriptUpdateForm
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

	scriptId := c.Param("script_id")

	script, err := models.ScriptMapper.FetchOne(project, scriptId)
	if err != nil {
		panic(err)
	}

	if script == nil {
		c.JSON(http.StatusNotFound, goerrors.New(goerrors.Error{
			Label: "invalid_script",
			Field: "script_id",
			Text:  "Invalid script ID in URL",
		}))
		return
	}

	script.Update(&form)

	if err := models.ScriptMapper.Update(script); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, script)
}

func (sc *ScriptController) getScriptList(c *gin.Context) {

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

	scripts, err := models.ScriptMapper.FetchAll(project)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, scripts)
}

func (sc *ScriptController) getScript(c *gin.Context) {

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

	scriptId := c.Param("script_id")

	format := c.Query("f")

	if strings.HasSuffix(scriptId, ".sh") {
		format = "bash"
		scriptId = strings.TrimSuffix(scriptId, ".sh")
	}

	script, err := models.ScriptMapper.FetchOne(project, scriptId)
	if err != nil {
		panic(err)
	}

	if script == nil {
		c.JSON(http.StatusNotFound, goerrors.New(goerrors.Error{
			Label: "invalid_script",
			Field: "script_id",
			Text:  "Invalid script ID in URL",
		}))
		return
	}

	switch format {
	case "bash":
		c.String(http.StatusOK, script.Content)
	case "xml":
		c.XML(http.StatusOK, script)
	default:
		c.JSON(http.StatusOK, script)
	}
}

func (sc *ScriptController) deleteScript(c *gin.Context) {

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

	scriptId := c.Param("script_id")

	script, err := models.ScriptMapper.FetchOne(project, scriptId)
	if err != nil {
		panic(err)
	}

	if script == nil {
		c.JSON(http.StatusNotFound, goerrors.New(goerrors.Error{
			Label: "invalid_script",
			Field: "script_id",
			Text:  "Invalid script ID in URL",
		}))
		return
	}

	if err := models.ScriptMapper.Delete(script); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, nil)
}

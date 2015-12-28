package api

//
// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/gotoolz/errors"
// 	"github.com/karhuteam/karhu/models"
// 	"net/http"
// 	"strings"
// )
//
// type ScriptController struct {
// }
//
// func NewScriptController(s *gin.RouterGroup) *ScriptController {
//
// 	ctl := &ScriptController{}
//
// 	s.POST("/apps/:id/env/:env/scripts", ctl.postScript)
// 	s.PUT("/apps/:id/env/:env/scripts/:script_id", ctl.putScript)
// 	s.GET("/apps/:id/env/:env/scripts", ctl.getScriptList)
// 	s.GET("/apps/:id/env/:env/scripts/:script_id", ctl.getScript)
// 	s.DELETE("/apps/:id/env/:env/scripts/:script_id", ctl.deleteScript)
//
// 	return ctl
// }
//
// func (sc *ScriptController) getAppEnv(c *gin.Context) (*models.Application, *models.Environment, error) {
//
// 	id := c.Param("id")
//
// 	app, err := models.ApplicationMapper.FetchOne(id)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	if app == nil {
// 		return nil, nil, errors.New(errors.Error{
// 			Label: "invalid_application",
// 			Field: "id",
// 			Text:  "Invalid application ID in URL",
// 		})
// 	}
//
// 	envId := c.Param("env")
//
// 	env, err := models.EnvironmentMapper.FetchOne(app, envId)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	if env == nil {
// 		return nil, nil, errors.New(errors.Error{
// 			Label: "invalid_environment",
// 			Field: "id",
// 			Text:  "Invalid environment ID in URL",
// 		})
// 	}
//
// 	return app, env, nil
// }
//
// func (sc *ScriptController) postScript(c *gin.Context) {
//
// 	var form models.ScriptCreateForm
// 	if err := c.Bind(&form); err != nil {
// 		c.AbortWithStatus(http.StatusBadRequest)
// 		return
// 	}
//
// 	if err := form.Validate(); err != nil {
// 		c.JSON(http.StatusBadRequest, err)
// 		return
// 	}
//
// 	_, env, err := sc.getAppEnv(c)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, err)
// 		return
// 	}
//
// 	script, err := models.ScriptMapper.FetchOneByName(env, form.Name)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	if script != nil {
// 		c.JSON(http.StatusConflict, errors.New(errors.Error{
// 			Label: "invalid_name",
// 			Field: "name",
// 			Text:  "Duplicate script name for this environment",
// 		}))
// 		return
// 	}
//
// 	script = models.ScriptMapper.Create(env, &form)
//
// 	if err := models.ScriptMapper.Save(script); err != nil {
// 		panic(err)
// 	}
//
// 	c.JSON(http.StatusCreated, script)
// }
//
// func (sc *ScriptController) putScript(c *gin.Context) {
//
// 	var form models.ScriptUpdateForm
// 	if err := c.Bind(&form); err != nil {
// 		c.AbortWithStatus(http.StatusBadRequest)
// 		return
// 	}
//
// 	if err := form.Validate(); err != nil {
// 		c.JSON(http.StatusBadRequest, err)
// 		return
// 	}
//
// 	_, env, err := sc.getAppEnv(c)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, err)
// 		return
// 	}
//
// 	scriptId := c.Param("script_id")
//
// 	script, err := models.ScriptMapper.FetchOne(env, scriptId)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	if script == nil {
// 		c.JSON(http.StatusNotFound, errors.New(errors.Error{
// 			Label: "invalid_script",
// 			Field: "script_id",
// 			Text:  "Invalid script ID in URL",
// 		}))
// 		return
// 	}
//
// 	script.Update(&form)
//
// 	if err := models.ScriptMapper.Update(script); err != nil {
// 		panic(err)
// 	}
//
// 	c.JSON(http.StatusOK, script)
// }
//
// func (sc *ScriptController) getScriptList(c *gin.Context) {
//
// 	_, env, err := sc.getAppEnv(c)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, err)
// 		return
// 	}
//
// 	scripts, err := models.ScriptMapper.FetchAll(env)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	c.JSON(http.StatusOK, scripts)
// }
//
// func (sc *ScriptController) getScript(c *gin.Context) {
//
// 	_, env, err := sc.getAppEnv(c)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, err)
// 		return
// 	}
//
// 	scriptId := c.Param("script_id")
//
// 	format := c.Query("f")
//
// 	if strings.HasSuffix(scriptId, ".sh") {
// 		format = "bash"
// 		scriptId = strings.TrimSuffix(scriptId, ".sh")
// 	}
//
// 	script, err := models.ScriptMapper.FetchOne(env, scriptId)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	if script == nil {
// 		c.JSON(http.StatusNotFound, errors.New(errors.Error{
// 			Label: "invalid_script",
// 			Field: "script_id",
// 			Text:  "Invalid script ID in URL",
// 		}))
// 		return
// 	}
//
// 	switch format {
// 	case "bash":
// 		c.String(http.StatusOK, script.Content)
// 	case "xml":
// 		c.XML(http.StatusOK, script)
// 	default:
// 		c.JSON(http.StatusOK, script)
// 	}
// }
//
// func (sc *ScriptController) deleteScript(c *gin.Context) {
//
// 	_, env, err := sc.getAppEnv(c)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, err)
// 		return
// 	}
//
// 	scriptId := c.Param("script_id")
//
// 	script, err := models.ScriptMapper.FetchOne(env, scriptId)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	if script == nil {
// 		c.JSON(http.StatusNotFound, errors.New(errors.Error{
// 			Label: "invalid_script",
// 			Field: "script_id",
// 			Text:  "Invalid script ID in URL",
// 		}))
// 		return
// 	}
//
// 	if err := models.ScriptMapper.Delete(script); err != nil {
// 		panic(err)
// 	}
//
// 	c.JSON(http.StatusOK, nil)
// }

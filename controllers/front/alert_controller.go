package front

import (
	// "fmt"
	"github.com/gin-gonic/gin"
	"github.com/karhuteam/karhu/models"
	"github.com/karhuteam/karhu/web"
	"log"
	"net/http"
	"strings"
	// "strconv"
	// "strings"
)

type AlertController struct {
}

func NewAlertController(s *web.Server) *AlertController {

	ctl := &AlertController{}

	s.GET("/alerts", ctl.getAlertsListAction)
	s.POST("/alerts/acknowledge/:id", ctl.postAlertsAcknowledgeAction)
	s.POST("/alerts/close/:id", ctl.postAlertsCloseAction)
	s.POST("/alerts/message/:id", ctl.postAlertsMessageAction)

	s.GET("/alerts-policies", ctl.getAlertsPolicyListAction)
	s.GET("/alerts-policies/add", ctl.getAddAlertsPolicyAction)
	s.POST("/alerts-policies/add", ctl.postAddAlertsPolicyAction)
	s.GET("/alerts-policies/edit/:id", ctl.getEditAlertsPolicyAction)
	s.POST("/alerts-policies/edit/:id", ctl.postEditAlertsPolicyAction)

	s.GET("/alerts-groups", ctl.getAlertsGroupsListAction)
	s.GET("/alerts-groups/add", ctl.getAddAlertsGroupAction)
	s.POST("/alerts-groups/add", ctl.postAddAlertsGroupAction)
	s.GET("/alerts-groups/edit/:id", ctl.getEditAlertsGroupAction)
	s.POST("/alerts-groups/edit/:id", ctl.postEditAlertsGroupAction)
	s.POST("/alerts-groups/delete/:id", ctl.postDeleteAlertsGroupAction)

	return ctl
}

func (ctl *AlertController) getAlertsListAction(c *gin.Context) {

	status := c.DefaultQuery("status", "open")

	alerts, err := models.AlertMapper.FetchAllStatus(status)
	if err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "alert_list.html", map[string]interface{}{
		"alerts":        alerts,
		"status_filter": status,
	})
}

func (ctl *AlertController) postAlertsAcknowledgeAction(c *gin.Context) {

	alert, err := models.AlertMapper.FetchOne(c.Param("id"))
	if err != nil {
		panic(err)
	}

	alert.Status = models.STATUS_ACKNOWLEDGE

	if err := models.AlertMapper.Update(alert); err != nil {
		panic(err)
	}

	c.Redirect(http.StatusFound, "/alerts?status=acknowledge")
}

func (ctl *AlertController) postAlertsCloseAction(c *gin.Context) {

	alert, err := models.AlertMapper.FetchOne(c.Param("id"))
	if err != nil {
		panic(err)
	}

	alert.Status = models.STATUS_CLOSED

	if err := models.AlertMapper.Update(alert); err != nil {
		panic(err)
	}

	c.Redirect(http.StatusFound, "/alerts?status=closed")
}

func (ctl *AlertController) postAlertsMessageAction(c *gin.Context) {

	text := strings.TrimSpace(c.PostForm("text"))
	if text == "" {
		c.Redirect(http.StatusFound, "/alerts")
		return
	}

	alert, err := models.AlertMapper.FetchOne(c.Param("id"))
	if err != nil {
		panic(err)
	}

	alert.AddMessage(text)

	if err := models.AlertMapper.Update(alert); err != nil {
		panic(err)
	}

	c.Redirect(http.StatusFound, "/alerts?status="+alert.Status)
}

func (ctl *AlertController) getAlertsPolicyListAction(c *gin.Context) {

	alertPolicies, err := models.AlertPolicyMapper.FetchAll()
	if err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "alert_policy_list.html", map[string]interface{}{
		"policies": alertPolicies,
	})
}

func (ctl *AlertController) getAddAlertsPolicyAction(c *gin.Context) {

	c.HTML(http.StatusOK, "alert_policy_add.html", nil)
}

func (ctl *AlertController) postAddAlertsPolicyAction(c *gin.Context) {

	var form models.AlertPolicyForm
	if err := c.Bind(&form); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := form.Validate(); err != nil {
		c.HTML(http.StatusOK, "alert_policy_add.html", map[string]interface{}{
			"errors": err.Errors,
			"form":   form,
		})
		return
	}

	ap := models.AlertPolicyMapper.Create(&form)

	if err := models.AlertPolicyMapper.Save(ap); err != nil {
		panic(err)
	}

	c.Redirect(http.StatusFound, "/alerts-policies")
}

func (ctl *AlertController) getEditAlertsPolicyAction(c *gin.Context) {

	ap, err := models.AlertPolicyMapper.FetchOne(c.Param("id"))
	if err != nil {
		panic(err)
	}

	if ap == nil {
		c.Redirect(http.StatusFound, "/alerts-policies")
		return
	}

	var form models.AlertPolicyForm
	form.Hydrate(ap)

	log.Println("interval", form.Interval)

	c.HTML(http.StatusOK, "alert_policy_add.html", map[string]interface{}{
		"form":   form,
		"policy": ap,
	})
}

func (ctl *AlertController) postEditAlertsPolicyAction(c *gin.Context) {

	var form models.AlertPolicyForm
	if err := c.Bind(&form); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ap, err := models.AlertPolicyMapper.FetchOne(c.Param("id"))
	if err != nil {
		panic(err)
	}

	if ap == nil {
		c.Redirect(http.StatusFound, "/alerts-policies")
		return
	}

	if err := form.Validate(); err != nil {
		c.HTML(http.StatusOK, "alert_policy_add.html", map[string]interface{}{
			"errors": err.Errors,
			"form":   form,
			"policy": ap,
		})
		return
	}

	ap.Update(&form)

	if err := models.AlertPolicyMapper.Update(ap); err != nil {
		panic(err)
	}

	c.Redirect(http.StatusFound, "/alerts-policies")
}

func (ctl *AlertController) getAlertsGroupsListAction(c *gin.Context) {

	groups, err := models.AlertGroupMapper.FetchAll()
	if err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "alert_group_list.html", map[string]interface{}{
		"groups": groups,
	})
}

func (ctl *AlertController) getAddAlertsGroupAction(c *gin.Context) {

	c.HTML(http.StatusOK, "alert_group_add.html", nil)
}

func (ctl *AlertController) postAddAlertsGroupAction(c *gin.Context) {

	var form models.AlertGroupCreateForm
	if err := c.Bind(&form); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := form.Validate(); err != nil {
		c.HTML(http.StatusOK, "alert_group_add.html", map[string]interface{}{
			"errors": err.Errors,
			"form":   form,
		})
		return
	}

	ag := models.AlertGroupMapper.Create(&form)

	if err := models.AlertGroupMapper.Save(ag); err != nil {
		panic(err)
	}

	c.Redirect(http.StatusFound, "/alerts-groups")
}

func (ctl *AlertController) getEditAlertsGroupAction(c *gin.Context) {

	ag, err := models.AlertGroupMapper.FetchOne(c.Param("id"))
	if err != nil {
		panic(err)
	}

	if ag == nil {
		c.Redirect(http.StatusFound, "/alerts-groups")
		return
	}

	var form models.AlertGroupUpdateForm
	form.Hydrate(ag)

	c.HTML(http.StatusOK, "alert_group_edit.html", map[string]interface{}{
		"form": form,
		"ag":   ag,
	})
}

func (ctl *AlertController) postEditAlertsGroupAction(c *gin.Context) {

	var form models.AlertGroupUpdateForm
	if err := c.Bind(&form); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ag, err := models.AlertGroupMapper.FetchOne(c.Param("id"))
	if err != nil {
		panic(err)
	}

	if ag == nil {
		c.Redirect(http.StatusFound, "/alerts-groups")
		return
	}

	if err := form.Validate(ag); err != nil {
		c.HTML(http.StatusOK, "alert_group_edit.html", map[string]interface{}{
			"errors": err.Errors,
			"form":   form,
			"ag":     ag,
		})
		return
	}

	ag.Update(&form)

	if err := models.AlertGroupMapper.Update(ag); err != nil {
		panic(err)
	}

	c.Redirect(http.StatusFound, "/alerts-groups")
}

func (ctl *AlertController) postDeleteAlertsGroupAction(c *gin.Context) {

	ag, err := models.AlertGroupMapper.FetchOne(c.Param("id"))
	if err != nil {
		panic(err)
	}

	if ag == nil {
		c.Redirect(http.StatusFound, "/alerts-groups")
		return
	}

	if err := models.AlertGroupMapper.Delete(ag); err != nil {
		panic(err)
	}

	c.Redirect(http.StatusFound, "/alerts-groups")
}

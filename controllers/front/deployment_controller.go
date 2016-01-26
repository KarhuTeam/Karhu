package front

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/gotoolz/env"
	"github.com/karhuteam/karhu/models"
	"github.com/karhuteam/karhu/web"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

type DeploymentController struct {
	upgrader websocket.Upgrader
}

func NewDeploymentController(s *web.Server) *DeploymentController {

	ctl := &DeploymentController{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

	// 1 - Show deployment
	s.GET("/application/deployment/:app_id/:deploy_id", ctl.getDeploymentAction)
	// 2 - Create a deployment
	s.POST("/application/deployment/:app_id/:build_id", ctl.postDeploymentAction)
	// 3 - Stream the output (ws)
	s.GET("/ws/application/deployment/:app_id/:deploy_id", ctl.getDeploymentWSAction)

	return ctl
}

func (ctl *DeploymentController) getApplication(c *gin.Context, id string) *models.Application {

	application, err := models.ApplicationMapper.FetchOne(id)

	// Error 500
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error_500.html", map[string]interface{}{
			"error": err,
		})
		return nil
	}

	// Error 404
	if application == nil {
		c.HTML(http.StatusNotFound, "error_404.html", map[string]interface{}{
			"title": "Application not found",
			"text":  "Application not found... It's not my fault",
		})
		return nil
	}

	return application
}

func (ctl *DeploymentController) getDeployment(c *gin.Context, application *models.Application, id string) *models.Deployment {

	deployment, err := models.DeploymentMapper.FetchOne(application, id)

	// Error 500
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error_500.html", map[string]interface{}{
			"error": err,
		})
		return nil
	}

	// Error 404
	if deployment == nil {
		c.HTML(http.StatusNotFound, "error_404.html", map[string]interface{}{
			"title": "Deploy not found",
			"text":  "Deploy not found... It's not my fault",
		})
		return nil
	}

	return deployment
}

/**
 * 1 - Show a deployment
 */
func (ctl *DeploymentController) getDeploymentAction(c *gin.Context) {

	appId := c.Param("app_id")
	deployId := c.Param("deploy_id")

	application := ctl.getApplication(c, appId)
	if application == nil {
		return
	}

	deployment := ctl.getDeployment(c, application, deployId)
	if deployment == nil {
		return
	}

	karhuURL, _ := url.Parse(env.GetDefault("PUBLIC_HOST", "http://127.0.0.1:8080"))
	webSocketProto := "ws"
	if karhuURL.Scheme == "https" {
		webSocketProto = "wss"
	}

	c.HTML(http.StatusOK, "deployment_show.html", map[string]interface{}{
		"application":     application,
		"deployment":      deployment,
		"public_host":     karhuURL.Host,
		"websocket_proto": webSocketProto,
	})
}

/**
 * 2 - Create a deployment
 */
func (ctl *DeploymentController) postDeploymentAction(c *gin.Context) {

	appId := c.Param("app_id")
	buildId := c.Param("build_id")

	resp, err := http.Post(fmt.Sprintf("http://localhost:8080/api/apps/%s/builds/%s/deploy", appId, buildId), "json/application", nil)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error_500.html", map[string]interface{}{
			"error": err,
		})
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error_500.html", map[string]interface{}{
			"error": err,
		})
		return
	}
	var deployment models.Deployment
	if err := json.Unmarshal(body, &deployment); err != nil {
		c.HTML(http.StatusInternalServerError, "error_500.html", map[string]interface{}{
			"error": err,
		})
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/application/deployment/%s/%s", appId, deployment.Id.Hex()))
}

/**
 * 3 - Stream a deployment (ws)
 */
func (ctl *DeploymentController) getDeploymentWSAction(c *gin.Context) {

	appId := c.Param("app_id")
	deployId := c.Param("deploy_id")

	application := ctl.getApplication(c, appId)
	if application == nil {
		log.Println("ws: application not found")
		return
	}

	deployment := ctl.getDeployment(c, application, deployId)
	if deployment == nil {
		log.Println("ws: deployment not found")
		return
	}

	sock, err := ctl.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("ws: ", err)
	}
	defer sock.Close()

	file, err := os.Open(path.Join(deployment.TmpPath, "karhu.log"))
	if err != nil {
		log.Println("ws: ", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line, _, err := reader.ReadLine()

		if err == io.EOF {
			time.Sleep(1 * time.Second)
			continue
		} else if err != nil {
			log.Println("ws: ", err)
			break
		}
		line = append(line, '\n')
		sock.WriteMessage(websocket.TextMessage, line)
	}
}

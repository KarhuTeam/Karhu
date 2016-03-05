package front

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/ec2"
	"github.com/gotoolz/env"
	"github.com/karhuteam/karhu/models"
	// "github.com/karhuteam/karhu/ressources/ssh"
	"github.com/karhuteam/karhu/web"
	"log"
	"net/http"
	"time"
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
	s.POST("/node/add/ec2", ctl.postNodeAddEc2Action)
	s.GET("/node/add/do", ctl.getNodeAddDOAction)
	s.POST("/node/access/add", ctl.postAddAccessAction)
	s.POST("/node/access/delete/:type", ctl.postDeleteAccessAction)
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

	a, err := models.AccessMapper.FetchOne("ec2")
	if err != nil {
		panic(err)
	}

	if a == nil {

		c.HTML(http.StatusOK, "node_add_ec2.html", map[string]interface{}{})
		return
	}

	auth, err := aws.GetAuth(a.AccessKey, a.PrivateKey, "", time.Now().Add(time.Hour))
	if err != nil {
		panic(err)
	}

	var vpcs []ec2.VPC
	var securityGroups []ec2.SecurityGroupInfo
	region := c.Query("availability_zone")
	vpc := c.Query("vpc")
	securityGroup := c.Query("security_group")
	if region != "" {

		awsec2 := ec2.New(auth, aws.Regions[region])
		res, _ := awsec2.DescribeVpcs(nil, nil)

		if res != nil {
			vpcs = res.VPCs
		}

		if vpc != "" {
			if groups, _ := awsec2.SecurityGroups(nil, nil); groups != nil {
				for _, g := range groups.Groups {
					if g.VpcId == vpc {
						securityGroups = append(securityGroups, g)
					}
				}
			}
		}

	}

	log.Println("vpcs:", vpcs)

	c.HTML(http.StatusOK, "node_add_ec2.html", map[string]interface{}{
		"AccessKey":      a.AccessKey,
		"AWSRegions":     aws.Regions,
		"VPCs":           vpcs,
		"SecurityGroups": securityGroups,
		"query": map[string]interface{}{
			"availability_zone": region,
			"vpc":               vpc,
			"security_group":    securityGroup,
		},
	})
}

func (pc *NodeController) postNodeAddEc2Action(c *gin.Context) {

	var form models.EC2NodeCreateForm
	if err := c.Bind(&form); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := form.Validate(); err != nil {
		c.Redirect(http.StatusFound, c.Request.Referer())
		return
	}

	a, err := models.AccessMapper.FetchOne("ec2")
	if err != nil {
		panic(err)
	}

	if a == nil {
		c.Redirect(http.StatusFound, c.Request.Referer())
		return
	}

	basicAuth := env.Get("BASIC_AUTH")
	if auth := env.Get("BASIC_AUTH"); auth != "" {
		basicAuth = "-u " + auth + " "
	}

	auth, err := aws.GetAuth(a.AccessKey, a.PrivateKey, "", time.Now().Add(time.Hour))
	if err != nil {
		panic(err)
	}
	awsec2 := ec2.New(auth, aws.Regions[form.AvailabilityZone])
	// Create public key
	// Waiting for merge pull request https://github.com/goamz/goamz/pull/111
	// {
	// key, err := ssh.GetPublicKey()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	if _, err := awsec2.ImportKeyPair(&ImportKeyPairOptions{
	// 		KeyName:           "karhu",
	// 		PublicKeyMaterial: string(key),
	// 	}); err != nil {
	// 		panic(err)
	// 	}
	// }

	if _, err := awsec2.RunInstances(&ec2.RunInstancesOptions{
		ImageId:        "ami-e31a6594",
		MinCount:       1,
		MaxCount:       0,
		KeyName:        "karhu",
		InstanceType:   form.InstanceType,
		SecurityGroups: []ec2.SecurityGroup{{Id: form.SecurityGroup}},
		// KernelId               :  string
		// RamdiskId              :  string
		UserData: []byte(fmt.Sprintf(`#!/bin/bash
sudo apt-get update && \
sudo apt-get install -y curl && \
curl %s"%s/api/nodes/register.sh?monit=1&ssh_port=22" | bash`, basicAuth, env.Get("PUBLIC_HOST"))),
		AvailabilityZone: "eu-west-1c", // Waiting for https://github.com/goamz/goamz/pull/112
		// PlacementGroupName     :  string
		Tenancy:    "default",
		Monitoring: form.Monitoring == "on",
		SubnetId:   "subnet-425a4f27", // Waiting for https://github.com/goamz/goamz/pull/112
		// DisableAPITermination  :  bool
		// ShutdownBehavior       :  string
		// PrivateIPAddress       :  string
		// IamInstanceProfile      : IamInstanceProfile
		// BlockDevices            : []BlockDeviceMapping
		// EbsOptimized            : bool
		// AssociatePublicIpAddress :bool
	}); err != nil {
		panic(err)
	}

	c.Redirect(http.StatusFound, c.Request.Referer())

}

func (pc *NodeController) getNodeAddDOAction(c *gin.Context) {
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

func (nc *NodeController) postAddAccessAction(c *gin.Context) {

	var form models.AccessCreateForm
	if err := c.Bind(&form); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := form.Validate(); err != nil {
		log.Println(err)
		c.Redirect(http.StatusFound, c.Request.Referer())
		return
	}

	a := models.AccessMapper.Create(&form)

	if err := models.AccessMapper.Save(a); err != nil {
		panic(err)
	}

	c.Redirect(http.StatusFound, c.Request.Referer())
}
func (nc *NodeController) postDeleteAccessAction(c *gin.Context) {

	typ := c.Param("type")

	a, err := models.AccessMapper.FetchOne(typ)
	if err != nil {
		panic(err)
	}

	if err := models.AccessMapper.Delete(a); err != nil {
		panic(err)
	}

	c.Redirect(http.StatusFound, c.Request.Referer())
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

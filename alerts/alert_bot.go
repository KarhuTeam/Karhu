package alerts

import (
	"github.com/karhuteam/karhu/models"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strings"
	"time"
)

type TargetKarhuHandler interface {
	HandleTargetKarhu() error
}

type TargetNodeHandler interface {
	HandleTargetNode(*models.Node) error
}

func Run() {

	for {
		time.Sleep(time.Second * 10)

		policies, err := models.AlertPolicyMapper.FetchAllForExecution()
		if err != nil {
			log.Println("alerts.Run:", err)
			continue
		}

		log.Println("Policies:", len(policies))

		for _, policy := range policies {
			func() {
				defer func() {
					if err := recover(); err != nil {
						log.Println("alerts.Run:", err)
					}
				}()

				check(policy)
			}()
		}
	}
}

func check(policy *models.AlertPolicy) {

	log.Println("Checking", policy.Name)

	var err error
	var chk interface{}
	// switch policy.CondType {
	// case "cond-http":
	// 	chk, err = NewCheckHTTP(policy)
	// }
	chk, err = NewCheckNagios(policy)

	if err != nil {
		log.Println("alerts.check:", err)
		return
	}

	switch policy.TargetType {
	case "target-karhu":
		handler, ok := chk.(TargetKarhuHandler)
		if !ok {
			log.Println("Check status, TargetKarhuHandler isn't supported")
			return
		}
		err := handler.HandleTargetKarhu()
		updateAlert(policy, nil, err)

	case "target-tag":
		handler, ok := chk.(TargetNodeHandler)
		if !ok {
			log.Println("Check status, TargetNodeHandler isn't supported")
			return
		}

		var tags []string
		for _, t := range policy.Target.(bson.M)["tags"].([]interface{}) {
			tags = append(tags, t.(string))
		}

		nodes, err := models.NodeMapper.FetchAllTags(tags)
		if err != nil {
			log.Println("Check status, FetchAllTags:", err)
			return
		}

		for _, n := range nodes {

			err := handler.HandleTargetNode(n)
			updateAlert(policy, n, err)
		}
	case "target-node":
		handler, ok := chk.(TargetNodeHandler)
		if !ok {
			log.Println("Check status, TargetNodeHandler isn't supported")
			return
		}

		for _, n := range policy.Target.(bson.M)["nodes"].([]interface{}) {
			node, err := models.NodeMapper.FetchOne(n.(string))
			if err != nil {
				log.Println("Check status, NodeMapper.FetchOne:", err)
				continue
			}

			if node == nil {
				continue
			}

			err = handler.HandleTargetNode(node)
			updateAlert(policy, node, err)
		}
	case "target-all":
		handler, ok := chk.(TargetNodeHandler)
		if !ok {
			log.Println("Check status, TargetNodeHandler isn't supported")
			return
		}

		nodes, err := models.NodeMapper.FetchAll()
		if err != nil {
			log.Println("Check status, NodeMapper.FetchAll:", err)
			return
		}

		for _, node := range nodes {
			err = handler.HandleTargetNode(node)
			updateAlert(policy, node, err)
		}
	default:
		log.Println("Check status, unsupported TargetType:", policy.TargetType)
		return
	}

	policy.NextAt = time.Now().Add(policy.Interval)
	log.Println("Next", policy.Name, "at", policy.NextAt)
	if err := models.AlertPolicyMapper.Update(policy); err != nil {
		log.Println("Check status, AlertPolicyMapper.Update:", err)
		// no return
	}
}

func updateAlert(policy *models.AlertPolicy, node *models.Node, checkResult error) {

	if checkResult != nil && strings.Contains(checkResult.Error(), "ssh: connect to host") &&
		strings.Contains(checkResult.Error(), "Operation timed out") {
		// Ignore in case of ssh fail
		return
	}

	if checkResult != nil { // Check to open issue

		alert, err := models.AlertMapper.FetchOneByPolicy(policy.Id.Hex(), []string{models.STATUS_OPEN, models.STATUS_ACKNOWLEDGE}, node)
		if err != nil {
			log.Println("Check status, FetchOneByPolicy:", err)
			return
		}

		if alert != nil { // We already have an alert
			log.Println("Alert:", *alert.Node)
			return
		}

		alert = models.AlertMapper.Create(policy, checkResult)
		if node != nil {
			alert.Node = node
		}

		if err := models.AlertMapper.Save(alert); err != nil {
			log.Println("Check status, AlertMapper.Save:", err)
			return
		}
	} else {

		alert, err := models.AlertMapper.FetchOneByPolicy(policy.Id.Hex(), []string{models.STATUS_OPEN, models.STATUS_ACKNOWLEDGE}, node)
		if err != nil {
			log.Println("Check status, FetchOneByPolicy:", err)
			return
		}

		if alert == nil { // No alert
			return
		}

		alert.Status = models.STATUS_CLOSED
		if err := models.AlertMapper.Update(alert); err != nil {
			log.Println("Check status, AlertMapper.Save:", err)
			return
		}
	}

	Notify(policy, node, checkResult)
}

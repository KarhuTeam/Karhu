package alerts

import (
	"github.com/karhuteam/karhu/models"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type TargetUrlHandler interface {
	HandleTargetUrl(string) error
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
			check(policy)
		}
	}
}

func check(policy *models.AlertPolicy) {

	log.Println("Checking", policy.Name)

	var err error
	var chk interface{}
	switch policy.CondType {
	case "cond-http":
		chk, err = NewCheckHTTP(policy)
	}

	if err != nil {
		log.Println("alerts.check:", err)
		return
	}

	var checkResult error
	switch policy.TargetType {
	case "target-url":
		host := policy.Target.(bson.M)["host"].(string)
		handler, ok := chk.(TargetUrlHandler)
		if !ok {
			log.Println("Check status, TargetUrlHandler isn't supported by", policy.CondType)
			return
		}
		checkResult = handler.HandleTargetUrl(host)
	default:
		log.Println("Check status, unsupported TargetType:", policy.TargetType)
		return
	}

	log.Println("Check status:", checkResult)

	policy.NextAt = time.Now().Add(policy.Interval)
	log.Println("Next", policy.Name, "at", policy.NextAt)
	if err := models.AlertPolicyMapper.Update(policy); err != nil {
		log.Println("Check status, AlertPolicyMapper.Update:", err)
		// no return
	}

	if checkResult != nil { // Check to open issue

		alert, err := models.AlertMapper.FetchOneByPolicy(policy.Id.Hex(), []string{models.STATUS_OPEN, models.STATUS_ACKNOWLEDGE})
		if err != nil {
			log.Println("Check status, FetchOneByPolicy:", err)
			return
		}

		if alert != nil { // We already have an alert
			return
		}

		alert = models.AlertMapper.Create(policy, checkResult)

		if err := models.AlertMapper.Save(alert); err != nil {
			log.Println("Check status, AlertMapper.Save:", err)
			return
		}
	} else {

		alert, err := models.AlertMapper.FetchOneByPolicy(policy.Id.Hex(), []string{models.STATUS_OPEN, models.STATUS_ACKNOWLEDGE})
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
}

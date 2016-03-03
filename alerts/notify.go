package alerts

import (
	"fmt"
	"github.com/gotoolz/mail"
	"github.com/karhuteam/karhu/models"
	"log"
	"strings"
	"time"
)

func Notify(policy *models.AlertPolicy, node *models.Node, problem error) {

	for _, group := range policy.AlertGroups {

		ag, err := models.AlertGroupMapper.FetchOne(group)
		if err != nil {
			log.Println("AlertGroupMapper.FetchOne:", err)
			continue
		}

		if ag == nil {
			log.Println("AlertGroupMapper.FetchOne: unknown group:", group)
			continue
		}

		// Handle  notify
		for _, method := range ag.Methods {

			switch method.Type {
			case "email":
				notifyEmail(policy, node, problem, method.Value)
			}
		}
	}
}

func notifyEmail(policy *models.AlertPolicy, node *models.Node, problem error, email string) {

	level := "OK"
	if problem != nil {
		for _, l := range []string{"CRITICAL", "WARNING"} {
			if strings.Contains(problem.Error(), l) {
				level = l
				break
			}
		}
	}

	subject := fmt.Sprintf("%s - %s", level, policy.Name)
	if node != nil {
		subject += " - " + node.Hostname
	}

	if err := mail.Send(mail.Email{
		From:    "monitoring@karhu.my-sign.org",
		Subject: subject,
		Content: fmt.Sprintf(`%s
%s`, subject, time.Now().String()),
		To: []string{email},
	}); err != nil {
		log.Println("alerts.notifyEmail:", err)
	}
}

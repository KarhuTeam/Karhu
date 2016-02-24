package alerts

import (
	"github.com/karhuteam/karhu/models"
	"gopkg.in/mgo.v2/bson"
	// "log"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
)

type CheckHTTP struct {
	Policy  *models.AlertPolicy
	Request *http.Request
	Port    string
	Status  int
}

func NewCheckHTTP(policy *models.AlertPolicy) (*CheckHTTP, error) {

	cond := policy.Cond.(bson.M)

	hostname := cond["hostname"].(string)
	path := cond["path"].(string)
	port := strconv.Itoa(cond["port"].(int))
	proto := cond["proto"].(string)
	status := cond["status"].(int)

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	req.URL.Scheme = proto
	req.Header.Set("User-Agent", "Karhu HTTP Check Bot v0.1")
	req.Host = hostname
	req.Header.Set("Host", hostname)

	return &CheckHTTP{
		Policy:  policy,
		Request: req,
		Port:    port,
		Status:  status,
	}, nil
}

func (c *CheckHTTP) HandleTargetUrl(host string) error {
	log.Println("HandleTargetUrl:", host)

	c.Request.URL.Host = fmt.Sprintf("%s:%s", host, c.Port)

	data, _ := httputil.DumpRequestOut(c.Request, false)
	log.Println(string(data))

	resp, err := http.DefaultClient.Do(c.Request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != c.Status {
		return errors.New(fmt.Sprintf("CheckHTTP: bad StatusCode, expected %d, got %s", c.Status, resp.Status))
	}

	return nil
}

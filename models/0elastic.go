package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gotoolz/env"
	"io"
	"io/ioutil"
	// "log"
	"net/http"
	"strings"
)

var esHost = env.GetDefault("ES_HOST", "http://127.0.0.1:9200")

type Hits struct {
	Total    int64   `json:"total"`
	MaxScore float32 `json:"max_score"`
	Hits     []Hit   `json:"hits"`
}
type Hit struct {
	Index  string           `json:"_index"`
	Type   string           `json:"_type"`
	Id     string           `json:"_id"`
	Score  float32          `json:"_score"`
	Source *json.RawMessage `json:"_source"`
}

type SearchResult struct {
	Took     int   `json:"took"`
	TimedOut bool  `json:"timed_out"`
	Hits     *Hits `json:"hits"`
}

func Search(index, body string) (*SearchResult, error) {

	result := new(SearchResult)
	if err := Elastic("POST", index, "/_search?pretty", body, result); err != nil {
		return nil, err
	}

	return result, nil
}

func Elastic(method, index, queryPath, body string, obj interface{}) error {

	var reader io.Reader = nil
	if body != "" {
		reader = strings.NewReader(body)
	}

	url := fmt.Sprintf("%s/%s%s", esHost, index, queryPath)
	// log.Printf("url: %s\n", url)

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return err
	}

	if reader != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status + " - " + string(data))
	}
	// log.Println("data:", string(data))

	if obj != nil {
		if err := json.Unmarshal(data, obj); err != nil {
			return err
		}
	}

	return nil
}

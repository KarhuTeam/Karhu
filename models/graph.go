package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type Graph struct {
	Title     string                   `json:"title"`
	Type      string                   `json:"type"`
	Stacked   bool                     `json:"stacked"`
	DataType  string                   `json:"data_type"`
	DataOrder []string                 `json:"data_order"`
	Data      map[string][]interface{} `json:"data"`
}
type Graphs []*Graph

func (g *Graph) Marshal() string {

	data, err := json.Marshal(g)
	if err != nil {
		panic(err)
	}

	return string(data)
}

type graphMapper struct{}

var GraphMapper = &graphMapper{}

type GraphTemplate struct {
	CollectdType  string
	Type          string
	Stacked       bool
	DataType      string
	ValueField    string
	TypeInstances []string
}

var graphQueries = map[string]GraphTemplate{
	"memory": {
		CollectdType:  "memory",
		Type:          "line",
		Stacked:       true,
		DataType:      "bytes",
		ValueField:    "value",
		TypeInstances: []string{"used", "cached", "buffered", "free"},
	},
	"cpu": {
		CollectdType:  "cpu",
		Type:          "line",
		Stacked:       true,
		DataType:      "percent",
		ValueField:    "value",
		TypeInstances: []string{"system", "user", "nice", "idle", "iowait", "irq", "softirq", "steal", "guest"},
	},
	// "load": {
	// 	CollectdType:  "load",
	// 	ValueField:   []string{"shortterm", "midterm", "longterm"},
	// 	TypeInstances: []string{""},
	// },
}

var GraphStats = []string{"cpu", "memory"}

func (gm *graphMapper) FetchOne(stat, host string, tm time.Time) (*Graph, error) {

	template, ok := graphQueries[stat]
	if !ok {
		return nil, errors.New("invalid stat: " + stat)
	}

	data := make(map[string][]interface{})

	for _, instance := range template.TypeInstances {

		query := fmt.Sprintf(`host: "%s" AND collectd_type: "%s"`, host, template.CollectdType)
		if instance != "" {
			query += fmt.Sprintf(` AND type_instance: "%s"`, instance)
		}

		res, err := Search("collectd-*", fmt.Sprintf(`{
"size": 500,
"sort": [
{
  "@timestamp": {
    "order": "desc",
    "unmapped_type": "boolean"
  }
}
],
"highlight": {
"fields": {
  "*": {}
},
"require_field_match": true,
"fragment_size": 2147483647
},
"query": {
"filtered": {
  "query": {
    "bool": {
        "must": [
        {
            "query_string": {
            "query": %s,
            "analyze_wildcard": true
            }
        },
        {
            "range": {
                "@timestamp": {
                    "gte": "%s"
                }
            }
        }
        ]
    }
  }
}
}
}`, strconv.Quote(query), tm.Format(time.RFC3339)))
		if err != nil {
			panic(err)
		}

		if res.Hits != nil {
			for _, hit := range res.Hits.Hits {

				values := make(map[string]interface{})

				if err := json.Unmarshal(*hit.Source, &values); err != nil {
					return nil, err
				}

				data[instance] = append(data[instance], values[template.ValueField])
			}
		}
	}

	return &Graph{
		Title:     fmt.Sprintf("%s %s", host, stat),
		Type:      template.Type,
		Stacked:   template.Stacked,
		DataType:  template.DataType,
		DataOrder: template.TypeInstances,
		Data:      data,
	}, nil
}

func (gm *graphMapper) FetchAll(hosts []string, stat, t string) (Graphs, error) {

	var graphs Graphs

	var tm time.Time
	switch t {
	case "last900":
		tm = time.Now().Add(-time.Minute * 15)
	case "last1800":
		tm = time.Now().Add(-time.Minute * 30)
	case "last3600":
		tm = time.Now().Add(-time.Minute * 60)
	case "last86400":
		tm = time.Now().Add(-time.Hour * 24)
	default: // 30min
		tm = time.Now().Add(-time.Minute * 30)
	}

	for _, host := range hosts {
		for _, currstat := range GraphStats {

			if stat != "all" && stat != currstat {
				continue
			}

			g, err := gm.FetchOne(currstat, host, tm)
			if err != nil {
				return nil, err
			}

			graphs = append(graphs, g)
		}
	}

	return graphs, nil
}

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
	CollectdType    string
	Type            string
	Stacked         bool
	DataType        string
	ValueFields     []string
	TypeInstances   []string
	PluginInstances []string
}

var graphQueries = map[string]GraphTemplate{
	"memory": {
		CollectdType:  "memory",
		Type:          "area",
		Stacked:       true,
		DataType:      "bytes",
		ValueFields:   []string{"value"},
		TypeInstances: []string{"used", "cached", "buffered", "free"},
	},
	"cpu": {
		CollectdType:  "cpu",
		Type:          "area",
		Stacked:       true,
		DataType:      "percent",
		ValueFields:   []string{"value"},
		TypeInstances: []string{"system", "user", "nice", "idle", "iowait", "irq", "softirq", "steal", "guest"},
	},
	"load": {
		CollectdType:  "load",
		Type:          "line",
		Stacked:       false,
		DataType:      "",
		ValueFields:   []string{"shortterm", "midterm", "longterm"},
		TypeInstances: nil,
	},
	"disk_ops": {
		CollectdType:    "disk_ops",
		Type:            "line",
		Stacked:         false,
		DataType:        "",
		ValueFields:     []string{"read", "write"},
		TypeInstances:   nil,
		PluginInstances: []string{"sda", "sda1", "sda2", "sda3", "sda4", "sda5", "sda5", "sdb", "sdb1", "sdb2", "sdb3", "sdb4", "sdb5", "sdb5", "sdc", "sdc1", "sdc2", "sdc3", "sdc4", "sdc5", "sdc5"},
	},
	"if_packets": {
		CollectdType:    "if_packets",
		Type:            "line",
		Stacked:         false,
		DataType:        "",
		ValueFields:     []string{"rx", "tx"},
		TypeInstances:   nil,
		PluginInstances: []string{"eth0", "eth1", "eth2", "lo", "docker0"},
	},
}

var GraphStats = []string{"cpu", "memory", "load", "disk_ops", "if_packets"}

const graphQuery = `{
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
}`

func (gm *graphMapper) fetch(template GraphTemplate, host, collectdType, instance, pluginInstance string, tm time.Time) map[string][]interface{} {

	data := make(map[string][]interface{})

	query := fmt.Sprintf(`host: "%s" AND collectd_type: "%s"`, host, collectdType)
	if instance != "" {
		query += fmt.Sprintf(` AND type_instance: "%s"`, instance)
	}
	if pluginInstance != "" {
		query += fmt.Sprintf(` AND plugin_instance: "%s"`, pluginInstance)
	}

	res, err := Search("collectd-*", fmt.Sprintf(graphQuery, strconv.Quote(query), tm.Format(time.RFC3339)))
	if err != nil {
		panic(err)
	}

	if res.Hits != nil {
		for _, hit := range res.Hits.Hits {

			values := make(map[string]interface{})

			if err := json.Unmarshal(*hit.Source, &values); err != nil {
				panic(err)
			}

			for _, field := range template.ValueFields {

				name := instance
				if pluginInstance != "" {
					name = pluginInstance
				}
				if field != "value" {
					name = name + " " + field
				}

				data[name] = append(data[name], values[field])
			}
		}
	}

	return data
}
func (gm *graphMapper) FetchOne(stat, host string, tm time.Time) (*Graph, error) {

	template, ok := graphQueries[stat]
	if !ok {
		return nil, errors.New("invalid stat: " + stat)
	}

	data := make(map[string][]interface{})
	if template.TypeInstances != nil {

		for _, instance := range template.TypeInstances {
			result := gm.fetch(template, host, template.CollectdType, instance, "", tm)
			for k, v := range result {
				data[k] = v
			}
		}
	} else if template.PluginInstances != nil {

		for _, instance := range template.PluginInstances {
			result := gm.fetch(template, host, template.CollectdType, "", instance, tm)
			for k, v := range result {
				data[k] = v
			}
		}

	} else {
		data = gm.fetch(template, host, template.CollectdType, "", "", tm)
	}

	if len(data) == 0 {
		return nil, nil
	}

	order := template.TypeInstances
	if order == nil {
		order = template.ValueFields
	}

	return &Graph{
		Title:     fmt.Sprintf("%s %s", host, stat),
		Type:      template.Type,
		Stacked:   template.Stacked,
		DataType:  template.DataType,
		DataOrder: order,
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

			if g == nil {
				continue
			}

			graphs = append(graphs, g)
		}
	}

	return graphs, nil
}

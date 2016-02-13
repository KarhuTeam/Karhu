package models

import (
	"fmt"
	"github.com/gotoolz/env"
	"github.com/karhuteam/karhu/ressources/logstash"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func LogstashRefreshTagsFilters() error {

	// Get all node
	nodes, err := NodeMapper.FetchAll()
	if err != nil {
		return err
	}

	filters := logstash.NewTagFilters()

	for _, n := range nodes {

		if len(n.Tags) == 0 {
			continue
		}

		var tags []string
		for _, t := range n.Tags {
			tags = append(tags, strconv.Quote(t))

		}

		log.Println("tags:", tags)

		filters.AddFilter(
			logstash.NewFilter("[host] == " + strconv.Quote(n.Hostname)).
				Mutate(logstash.NewMutate("karhu_tags", fmt.Sprintf("[ %s ]", strings.Join(tags, ", ")))))
	}

	data, err := filters.Marshal()
	if err != nil {
		return err
	}

	log.Println(string(data))

	if err := ioutil.WriteFile(env.GetDefault("LOGSTASH_TAGS_FILTERS", "./logstash/conf.d/10-tags-filters.conf"), data, 0644); err != nil {
		return err
	}

	return nil
}

package models

import (
	"fmt"
	"github.com/karhuteam/karhu/ressources/logstash"
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

	return nil
}

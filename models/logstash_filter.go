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

		var tags = []string{`"karhu"`}
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

func LogstashRefreshApplicationsFilters() error {

	// Get all node
	apps, err := ApplicationMapper.FetchAll()
	if err != nil {
		return err
	}

	filters := logstash.NewTagFilters()

	for _, app := range apps {

		// Fetch logfiles
		logfiles, err := LogfileMapper.FetchAllEnabled(app)
		if err != nil {
			return err
		}

		if len(logfiles) == 0 {
			continue
		}

		var filePaths []string
		for _, lf := range logfiles {
			filePaths = append(filePaths, strconv.Quote(lf.Path))
		}

		// Forge application tags
		var conds []string
		if len(filePaths) > 1 {
			conds = append(conds, fmt.Sprintf("[source] in [ %s ]", strings.Join(filePaths, ", ")))
		} else {
			conds = append(conds, fmt.Sprintf("[source] == %s", filePaths[0]))
		}
		for _, t := range app.Tags {
			conds = append(conds, fmt.Sprintf("%s in [karhu_tags]", strconv.Quote(t)))
		}

		filters.AddFilter(
			logstash.NewFilter(strings.Join(conds, " and ")).
				Mutate(logstash.NewMutate("karhu_app", strconv.Quote(app.Name))))
	}

	data, err := filters.Marshal()
	if err != nil {
		return err
	}

	log.Println(string(data))

	if err := ioutil.WriteFile(env.GetDefault("LOGSTASH_APPS_FILTERS", "./logstash/conf.d/11-apps-filters.conf"), data, 0644); err != nil {
		return err
	}

	return nil
}

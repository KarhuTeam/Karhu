package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type logMapper struct{}

var LogMapper = &logMapper{}

const logIndex = "filebeat-*"

type Log struct {
	Message   string      `json:"message"`
	Timestamp time.Time   `json:"@timestamp"`
	Host      string      `json:"host"`
	Source    string      `json:"source"`
	InputType string      `json:"input_type"`
	Fields    interface{} `json:"fields"`
	Type      string      `json:"type"`
	Tags      []string    `json:"tags"`
}

func (l *Log) TagsString() string {
	return strings.Join(l.Tags, ", ")
}

type Logs []*Log

func (lm *logMapper) Search(q string, count int) (Logs, error) {

	q = strconv.Quote(q)

	query := fmt.Sprintf(`{
  "size": %d,
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
        "query_string": {
          "query": %s,
          "analyze_wildcard": true
        }
      }
    }
  }
}`, count, q)
	//   "filter": {
	//     "bool": {
	//       "must": [
	//         {
	//           "range": {
	//             "@timestamp": {
	//               "gte": 1452640709890,
	//               "lte": 1455232709890,
	//               "format": "epoch_millis"
	//             }
	//           }
	//         }
	//       ],
	//       "must_not": []
	//     }
	//   }
	//     }
	//   }
	// }`, count, q))

	// log.Println("query:", query)
	res, err := Search(logIndex, query)

	if err != nil {
		return nil, err
	}

	var logs Logs
	if res.Hits != nil {
		for _, hit := range res.Hits.Hits {

			l := new(Log)
			if err := json.Unmarshal(*hit.Source, l); err != nil {
				return nil, err
			}

			logs = append(logs, l)
		}
	}
	return logs, nil

}

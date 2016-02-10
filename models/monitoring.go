package models

// import (
// 	goerrors "errors"
// 	"github.com/gotoolz/errors"
// 	"github.com/gotoolz/validator"
// 	"github.com/karhuteam/karhu/ressources/ssh"
// 	"gopkg.in/mgo.v2"
// 	"gopkg.in/mgo.v2/bson"
// 	"log"
// 	"strconv"
// 	"time"
// )

import (
)

//// Datas list:
// cpu_value
// df_value
// disk_read
// disk_write
// entropy_value
// interface_rx
// interface_tx
// irq_value
// load_longterm
// load_midterm
// load_shortterm
// memory_value
// swap_value
// users_value



type NodeMonitor struct {
	Name	string		`json:"name" bson:"name"`
}

func GetDefaultMonitoring() *NodeMonitor {
	Influx()

	return &NodeMonitor {
		Name:	"COUCOULECPU",
	}
}
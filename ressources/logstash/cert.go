package logstash

import (
	"github.com/gotoolz/env"
	"io/ioutil"
)

func GetCert() ([]byte, error) {

	return ioutil.ReadFile(env.Get("LOGSTASH_TLS_CRT"))
}

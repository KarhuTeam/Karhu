package models

import (
	"github.com/wayt/goenv"
	"gopkg.in/mgo.v2"
	"log"
)

type Mongo struct {
	s  *mgo.Session
	db string
}

var mongo *Mongo

// Don't forget to close c.Database.Session
func (m *Mongo) C(name string) *mgo.Collection {
	return m.s.Copy().DB(m.db).C(name)
}

func C(name string) *mgo.Collection {

	return mongo.C(name)
}

func init() {

	hosts := goenv.GetDefault("MGO_HOSTS", "127.0.0.1")
	db := goenv.GetDefault("MGO_DB", "karhu")

	log.Println("Connecting to mgo:", hosts)
	session, err := mgo.Dial(hosts)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Eventual, true)

	mongo = &Mongo{
		s:  session,
		db: db,
	}
}

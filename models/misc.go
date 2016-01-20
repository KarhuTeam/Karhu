package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Useful for application query where both _id and name are unique
func findIdOrSlugQuery(idOrSlug string) bson.M {

	key := "name"
	var value interface{} = idOrSlug
	if bson.IsObjectIdHex(idOrSlug) {
		key = "_id"
		value = bson.ObjectIdHex(idOrSlug)
	}

	return bson.M{key: value}
}

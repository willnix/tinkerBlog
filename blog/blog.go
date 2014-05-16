package blog

import (
	"labix.org/v2/mgo"
	"time"

	"labix.org/v2/mgo/bson"
)

// struct for dbquery results
// mgo requires the string literal tags
type Entry struct {
	ObjId   bson.ObjectId `bson:"_id,omitempty" form:"-"`
	Id      string        `bson:"-" form:"id"`
	Title   string        `bson:"title" form:"title"`
	Author  string        `bson:"author,omitempty" form:"-"`
	Text    string        `bson:"text" form:"text"`
	Written time.Time     `bson:"written,omitempty" form:"-"`
}

// default collection name
const dbCollectionName = "blogEntries"

type MgoBlog struct {
	collectionName string

	db *mgo.Database
}

func NewMgoBlog(db *mgo.Database) Blogger {
	return MgoBlog{
		db: db,

		collectionName: dbCollectionName,
	}
}

func NewMgoBlogWithCollectionName(db *mgo.Database, collName string) Blogger {
	return MgoBlog{
		db: db,

		collectionName: collName,
	}
}

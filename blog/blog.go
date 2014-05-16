package blog

import (
	"net/http"
	"time"

	"github.com/martini-contrib/binding"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

// struct for dbquery results
// mgo requires the string literal tags
type Entry struct {
	ObjId   bson.ObjectId `bson:"_id,omitempty" form:"-"`
	Id      string        `bson:"-" form:"id"`
	Author  string        `bson:"author,omitempty" form:"-"`
	Written time.Time     `bson:"written,omitempty" form:"-"`
	Title   string        `bson:"title" form:"title" binding:"required"`
	Text    string        `bson:"text" form:"text" binding:"required"`
}

func (e Entry) Validate(errors binding.Errors, req *http.Request) binding.Errors {

	if len(e.Title) == 0 {
		errors = append(errors, binding.Error{
			FieldNames:     []string{"title"},
			Classification: "IncompleteError",
			Message:        "Title can't be empty",
		})
	}

	if len(e.Text) == 0 {
		errors = append(errors, binding.Error{
			FieldNames:     []string{"text"},
			Classification: "IncompleteError",
			Message:        "Text can't be empty",
		})
	}

	return errors
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

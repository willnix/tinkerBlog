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

// defaults
const (
	blogCollName = "blogEntries"
	blogDbName   = "blog"
)

type MgoBlog struct {
	dbName, collName string
	s                *mgo.Session
}

func (m MgoBlog) getCollection() (*mgo.Collection, *mgo.Session) {
	s := m.s.Clone()

	return s.DB(m.dbName).C(m.collName), s
}

type Options struct {
	DbName         string
	CollectionName string
}

func NewMgoBlog(session *mgo.Session, o *Options) Blogger {
	b := MgoBlog{s: session}

	if o == nil {
		b.collName = blogCollName
		b.dbName = blogDbName
	} else {
		if o.CollectionName != "" {
			b.collName = o.CollectionName
		} else {
			b.collName = blogCollName
		}

		if o.DbName != "" {
			b.dbName = o.DbName
		} else {
			b.dbName = blogDbName
		}
	}

	return b
}

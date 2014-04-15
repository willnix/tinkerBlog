package blog

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/russross/blackfriday"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"strconv"
	"time"
)

// struct for dbquery results
// mgo requires the string literal tags
type dbBlogEntry struct {
	ObjId   bson.ObjectId "_id,omitempty"
	Id      int           "id"
	Title   string        "title"
	Author  string        "author"
	Text    string        "text"
	Written time.Time     "written"
}

func BlogEntryList(ren render.Render, db *mgo.Database) {
	var results []dbBlogEntry

	// Load all Blogentries in the results slice
	// (sorted descending according to id)
	db.C("blogEntries").Find(nil).Sort("-id").All(&results)

	for i, _ := range results {
		results[i].Text = string(blackfriday.MarkdownBasic([]byte(results[i].Text)))
	}

	// render the template using the results from the db
	ren.HTML(200, "blogEntryList", results)
}

func BlogEntry(ren render.Render, db *mgo.Database, args martini.Params) {
	var result dbBlogEntry

	Id, _ := strconv.Atoi(args["Id"])

	// Find Blogentry by Id (should be only one)
	db.C("blogEntries").Find(bson.M{"id": Id}).One(&result)

	result.Text = string(blackfriday.MarkdownBasic([]byte(result.Text)))

	// render the template using the result from the db
	ren.HTML(200, "blogEntry", result)
}

func About(ren render.Render) {
	ren.HTML(200, "about", nil)
}

func Impressum(ren render.Render) {
	ren.HTML(200, "impressum", nil)
}

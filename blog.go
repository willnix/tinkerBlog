package main

import (
	"fmt"
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
	ObjId   bson.ObjectId `_id,omitempty`
	Id      int           `form:"id"`
	Title   string        `form:"title"`
	Author  string        `form:"author"`
	Text    string        `form:"text"`
	Written time.Time     `form:"written"`
}

func BlogEntryList(ren render.Render, db *mgo.Database) {
	var results []dbBlogEntry

	// Load all Blogentries in the results slice
	// (sorted descending according to id)
	db.C("blogEntries").Find(nil).Sort("-id").All(&results)

	for i, _ := range results {
		results[i].Text = string(blackfriday.MarkdownCommon([]byte(results[i].Text)))
	}

	// render the template using the results from the db
	ren.HTML(200, "blogEntryList", results)
}

func BlogEntry(ren render.Render, db *mgo.Database, args martini.Params) {
	var result dbBlogEntry

	Id, _ := strconv.Atoi(args["Id"])

	// Find Blogentry by Id (should be only one)
	db.C("blogEntries").Find(bson.M{"id": Id}).One(&result)

	fmt.Println(string(blackfriday.MarkdownCommon([]byte(result.Text))))

	result.Text = string(blackfriday.MarkdownCommon([]byte(result.Text)))

	// render the template using the result from the db
	ren.HTML(200, "blogEntry", result)
}

func addBlogEntrySubmit(blogEntry dbBlogEntry, ren render.Render, db *mgo.Database) {
	blogEntry.Written = time.Now()

	db.C("blogEntries").Insert(blogEntry)

	// render the template using the result from the db
	ren.HTML(200, "addBlogEntry", nil)
}

func addBlogEntry(ren render.Render) {
	ren.HTML(200, "addBlogEntry", nil)
}

func About(ren render.Render) {
	ren.HTML(200, "about", nil)
}

func Impressum(ren render.Render) {
	ren.HTML(200, "impressum", nil)
}

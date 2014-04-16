package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/russross/blackfriday"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

// struct for dbquery results
// mgo requires the string literal tags
type dbBlogEntry struct {
	ObjId   bson.ObjectId `_id,omitempty`
	Title   string        `form:"title"`
	Author  string        `form:"author"`
	Text    string        `form:"text"`
	Written time.Time     `form:"written"`
}

func BlogEntryList(ren render.Render, db *mgo.Database) {
	var results []dbBlogEntry

	// Load all Blogentries in the results slice
	// (sorted descending according to id)
	db.C("blogEntries").Find(nil).Sort("-written").All(&results)

	for i, _ := range results {
		results[i].Text = string(blackfriday.MarkdownCommon([]byte(results[i].Text)))
	}

	// render the template using the results from the db
	ren.HTML(200, "blogEntryList", results)
}

func BlogEntry(ren render.Render, db *mgo.Database, args martini.Params) {
	// validate the post id
	if bson.IsObjectIdHex(args["Id"]) {
		entryId := bson.ObjectIdHex(args["Id"])

		var result dbBlogEntry

		// Find Blogentry by Id (should be only one)
		db.C("blogEntries").Find(bson.M{"_id": entryId}).One(&result)

		result.Text = string(blackfriday.MarkdownCommon([]byte(result.Text)))

		// render the template using the result from the db
		ren.HTML(200, "blogEntry", result)

	} else {
		// invalid post id
		ren.Data(400, []byte("Your request was bad and you should feeld bad!"))
	}

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

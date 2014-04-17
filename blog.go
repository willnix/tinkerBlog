package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessionauth"
	"github.com/martini-contrib/sessions"
	"github.com/russross/blackfriday"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

// struct for dbquery results
// mgo requires the string literal tags
type dbBlogEntry struct {
	ObjId   bson.ObjectId `bson:"_id,omitempty" form:"-"`
	Id      string        `bson:"-" form:"id"`
	Title   string        `bson:"title" form:"title"`
	Author  string        `bson:"author" form:"author"`
	Text    string        `bson:"text" form:"text"`
	Written time.Time     `bson:"written" form:"written"`
}

// List all blog entries
func BlogEntryList(ren render.Render, db *mgo.Database) {
	var results []dbBlogEntry

	// Load all Blogentries in the results slice
	// (sorted descending according to id)
	err := db.C(dbCollectionEntries).Find(nil).Sort("-written").All(&results)
	if err != nil {
		ren.JSON(500, err)
		return
	}

	for i, _ := range results {
		results[i].Text = string(blackfriday.MarkdownCommon([]byte(results[i].Text)))
	}

	// render the template using the results from the db
	ren.HTML(200, "blogEntryList", results)
}

// Show single blog entry
func BlogEntry(ren render.Render, db *mgo.Database, args martini.Params) {
	// validate the post id
	if !bson.IsObjectIdHex(args["Id"]) {
		ren.Data(400, []byte("Your request was bad and you should feeld bad!"))
		return
	}
	entryId := bson.ObjectIdHex(args["Id"])
	var result dbBlogEntry

	// Find Blogentry by Id (should be only one)
	err := db.C("blogEntries").Find(bson.M{"_id": entryId}).One(&result)
	if err != nil {
		ren.JSON(500, err)
		return
	}
	result.Text = string(blackfriday.MarkdownCommon([]byte(result.Text)))

	// render the template using the result from the db
	ren.HTML(200, "blogEntry", result)
}

// Submit new or update existing blog entry
func BlogEntrySubmit(user sessionauth.User, blogEntry dbBlogEntry, ren render.Render, db *mgo.Database) {
	blogEntry.Written = time.Now()
	// validate the post id
	if !bson.IsObjectIdHex(blogEntry.Id) {
		ren.Data(400, []byte("Your request was bad and you should feeld bad!"))
		return
	}
	blogEntry.ObjId = bson.ObjectIdHex(blogEntry.Id)

	// Set author to session user
	var userData UserModel
	userData.GetById(user.UniqueId())
	blogEntry.Author = userData.Username

	_, err := db.C(dbCollectionEntries).Upsert(bson.M{"_id": blogEntry.ObjId}, blogEntry)
	if err != nil {
		ren.JSON(500, err)
		return
	}

	// render the template using the result from the db
	ren.HTML(200, "addBlogEntry", nil)
}

// Display empty form to write new blog entry
func AddBlogEntryForm(ren render.Render) {
	ren.HTML(200, "addBlogEntry", nil)
}

// Display prefilled form to edit existing blog entry
func EditBlogEntryForm(ren render.Render, db *mgo.Database, args martini.Params) {
	// validate the post id
	if !bson.IsObjectIdHex(args["Id"]) {
		ren.Data(400, []byte("Your request was bad and you should feeld bad!"))
	}
	entryId := bson.ObjectIdHex(args["Id"])
	var result dbBlogEntry

	// Find Blogentry by Id (should be only one)
	err := db.C("blogEntries").Find(bson.M{"_id": entryId}).One(&result)
	if err != nil {
		ren.JSON(500, err)
		return
	}

	// render the template using the result from the db
	ren.HTML(200, "editBlogEntry", result)
}

func LoginForm(ren render.Render) {
	ren.HTML(200, "login", nil)
}

func Logout(session sessions.Session, user sessionauth.User, ren render.Render) {
	sessionauth.Logout(session, user)
	ren.Redirect("/")
}

func About(ren render.Render) {
	ren.HTML(200, "about", nil)
}

func Impressum(ren render.Render) {
	ren.HTML(200, "impressum", nil)
}

package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessionauth"
	"github.com/martini-contrib/sessions"
	"github.com/willnix/blackfriday"
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
	Author  string        `bson:"author,omitempty" form:"-"`
	Text    string        `bson:"text" form:"text"`
	Written time.Time     `bson:"written,omitempty" form:"-"`
}

// List all blog entries
func BlogEntryList(ren render.Render, db *mgo.Database) {
	var results []dbBlogEntry

	// Load all Blogentries in the results slice
	// (sorted descending by date)
	err := db.C(dbCollectionEntries).Find(nil).Sort("-_written").All(&results)
	if err != nil {
		ren.JSON(500, err)
		return
	}

	for i, _ := range results {
		results[i].Text = string(blackfriday.MarkdownCommon([]byte(results[i].Text)))
	}

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

	// Find Blogentry by Id
	err := db.C("blogEntries").Find(bson.M{"_id": entryId}).One(&result)
	if err != nil {
		ren.JSON(500, err)
		return
	}

	result.Text = string(blackfriday.MarkdownCommon([]byte(result.Text)))

	ren.HTML(200, "blogEntry", result)
}

// Submit new or update existing blog entry
func BlogEntrySubmit(user sessionauth.User, blogEntry dbBlogEntry, ren render.Render, db *mgo.Database) {
	// if we have a valid ObjId we assume it's an update
	if bson.IsObjectIdHex(blogEntry.Id) {
		blogEntry.ObjId = bson.ObjectIdHex(blogEntry.Id)
	} else {
		// no valid ObjId, so we assume it's a new post
		// and generate a new one
		blogEntry.ObjId = bson.NewObjectId()
	}
	// Set author to session user
	var userData UserModel
	userData.GetById(user.UniqueId())
	blogEntry.Author = userData.Username
	// set creation datetime
	blogEntry.Written = time.Now()

	_, err := db.C(dbCollectionEntries).UpsertId(blogEntry.ObjId, blogEntry)
	if err != nil {
		ren.JSON(500, err)
		return
	}

	// show new or updated entry
	ren.Redirect("/post/" + blogEntry.ObjId.Hex())
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

	// Find Blogentry by Id
	err := db.C("blogEntries").Find(bson.M{"_id": entryId}).One(&result)
	if err != nil {
		ren.JSON(500, err)
		return
	}

	ren.HTML(200, "editBlogEntry", result)
}

// Delete entry
func DeleteBlogEntry(ren render.Render, db *mgo.Database, args martini.Params) {
	// validate the post id
	if !bson.IsObjectIdHex(args["Id"]) {
		ren.Data(400, []byte("Your request was bad and you should feeld bad!"))
	}
	entryId := bson.ObjectIdHex(args["Id"])

	// Delete entry
	err := db.C("blogEntries").Remove(bson.M{"_id": entryId})
	if err != nil {
		ren.JSON(500, err)
		return
	}

	ren.Redirect("/")
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

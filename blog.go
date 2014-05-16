package main

import (
	"time"

	"github.com/go-martini/martini"
	"github.com/gorilla/feeds"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessionauth"
	"github.com/russross/blackfriday"
	"github.com/willnix/tinkerBlog/blog"
)

// List all blog entries
func BlogEntryList(ren render.Render, b blog.Blogger) {

	results, err := b.LatestEntries()
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
func BlogEntry(ren render.Render, b blog.Blogger, args martini.Params) {
	result, err := b.FindById(args["Id"])
	switch {
	case err == nil:
		result.Text = string(blackfriday.MarkdownCommon([]byte(result.Text)))
		ren.HTML(200, "blogEntry", result)
	case err == blog.ErrBadObjectId:
		ren.Data(400, []byte("Your request was bad and you should feeld bad!"))
	default:
		ren.JSON(500, err)
	}

}

// Submit new or update existing blog entry
func BlogEntrySubmit(user sessionauth.User, blogEntry blog.Entry, ren render.Render, b blog.Blogger) {

	// Set author to session user
	var userData UserModel
	userData.GetById(user.UniqueId())
	blogEntry.Author = userData.Username

	err := b.Save(&blogEntry)
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
func EditBlogEntryForm(ren render.Render, b blog.Blogger, args martini.Params) {

	result, err := b.FindById(args["Id"])
	switch {
	case err == nil:
		ren.HTML(200, "editBlogEntry", result)
	case err == blog.ErrBadObjectId:
		ren.Data(400, []byte("Your request was bad and you should feeld bad!"))
	default:
		ren.JSON(500, err)
	}

}

// Delete entry
func DeleteBlogEntry(ren render.Render, b blog.Blogger, args martini.Params) {

	err := b.Delete(args["Id"])
	switch {
	case err == nil:
		ren.Redirect("/")
	case err == blog.ErrBadObjectId:
		ren.Data(400, []byte("Your request was bad and you should feeld bad!"))
	default:
		ren.JSON(500, err)
	}

}

func RSS(ren render.Render, b blog.Blogger) string {

	results, err := b.LatestEntries()
	if err != nil {
		ren.JSON(500, err)
		return ""
	}

	now := time.Now()
	feed := &feeds.Feed{
		Title:       "tinkerBlog",
		Link:        &feeds.Link{Href: "http://localhost:3000/"},
		Description: "longcat is long",
		Author:      &feeds.Author{"kantorkel", "mail@example.xkcd"},
		Created:     now,
	}

	feed.Items = []*feeds.Item{}
	for i, _ := range results {
		feed.Items = append(feed.Items,
			&feeds.Item{
				Title:       results[i].Title,
				Link:        &feeds.Link{Href: "http://localhost:3000/post/" + results[i].ObjId.Hex()},
				Description: results[i].Author,
				Created:     results[i].Written,
			})
	}

	atom, err := feed.ToAtom()
	if err != nil {
		ren.JSON(500, err)
		return ""
	}

	return atom
}

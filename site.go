package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/auth"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"html/template"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

// Middleware handler for mongodb
func Mongo() martini.Handler {
	session, err := mgo.Dial("localhost/blog")
	if err != nil {
		panic(err)
	}

	return func(c martini.Context) {
		reqSession := session.Clone()
		c.Map(reqSession.DB("blog"))
		defer reqSession.Close()

		c.Next()
	}
}

func main() {
	// BasicAuth credentials for admin functions
	username := "username"
	password := "password"

	m := martini.Classic()

	//needs import ("time")
	m.Use(render.Renderer(render.Options{
		Directory: "templates",
		Layout:    "layout",
		Funcs: []template.FuncMap{
			{
				"formatTime": func(args ...interface{}) string {
					t1 := time.Time(args[0].(time.Time))
					return t1.Format("Jan 2, 2006 at 3:04pm (MST)")
				},
				"formatId": func(args ...interface{}) string {
					id := args[0].(bson.ObjectId)
					return id.Hex()
				},
				"unescaped": func(args ...interface{}) template.HTML {
					return template.HTML(args[0].(string))
				},
			},
		},
	}))

	// Middleware for mongodb connection
	m.Use(Mongo())

	// Setup static file serving
	m.Use(martini.Static("assets"))

	// Setup routing
	m.Get("/", BlogEntryList)
	m.Post("/blog/add/submit", auth.Basic(username, password), binding.Form(dbBlogEntry{}), addBlogEntrySubmit)
	m.Get("/blog/add", auth.Basic(username, password), addBlogEntry)
	m.Get("/post/:Id", BlogEntry)
	m.Get("/about", About)
	m.Get("/impressum", Impressum)

	m.Run()
}

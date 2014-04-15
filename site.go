package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"html/template"
	"labix.org/v2/mgo"
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
	m.Get("/post/:Id", BlogEntry)
	m.Get("/about", About)
	m.Get("/impressum", Impressum)

	m.Run()
}

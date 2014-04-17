package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessionauth"
	"github.com/martini-contrib/sessions"

	"html/template"
	"labix.org/v2/mgo/bson"
	"time"
)

func main() {
	store := sessions.NewCookieStore([]byte("BestSecretEvvaaaarr!!!"))

	m := martini.Classic()

	// Setup Renderer with some format functions for time.Time,
	// bson.ObjId and unescaped HTML
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

	// Setup session auth
	store.Options(sessions.Options{
		MaxAge: 0,
	})
	m.Use(sessions.Sessions("blogSession", store))
	m.Use(sessionauth.SessionUser(GenerateAnonymousUser))
	sessionauth.RedirectUrl = "/new-login"
	sessionauth.RedirectParam = "new-next"

	// Middleware for mongodb connection
	m.Use(Mongo())

	// Setup static file serving
	m.Use(martini.Static("assets"))

	// Setup routing
	m.Get("/", BlogEntryList)
	m.Get("/post/:Id", BlogEntry)
	m.Get("/about", About)
	m.Get("/impressum", Impressum)

	// login stuff
	m.Get("/blog/add", sessionauth.LoginRequired, AddBlogEntryForm)
	m.Post("/blog/add", sessionauth.LoginRequired, binding.Form(dbBlogEntry{}), BlogEntrySubmit)

	m.Get("/blog/edit/:Id", sessionauth.LoginRequired, EditBlogEntryForm)
	m.Post("/blog/edit", sessionauth.LoginRequired, binding.Form(dbBlogEntry{}), BlogEntrySubmit)

	m.Get("/new-login", LoginForm)

	m.Post("/new-login", binding.Bind(UserModel{}), ValidateLogin)

	m.Get("/logout", sessionauth.LoginRequired, Logout)

	m.Run()
}

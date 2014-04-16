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
	m.Get("/blog/add", sessionauth.LoginRequired, addBlogEntry)
	m.Post("/blog/add/submit", sessionauth.LoginRequired, binding.Form(dbBlogEntry{}), addBlogEntrySubmit)

	m.Get("/new-login", func(r render.Render) {
		r.HTML(200, "login", nil)
	})

	m.Post("/new-login", binding.Bind(UserModel{}), ValidateLogin)

	m.Get("/logout", sessionauth.LoginRequired, func(session sessions.Session, user sessionauth.User, r render.Render) {
		sessionauth.Logout(session, user)
		r.Redirect("/")
	})

	m.Run()
}

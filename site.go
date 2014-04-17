package main

import (
	"html/template"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessionauth"
	"github.com/martini-contrib/sessions"
)

func main() {
	store := sessions.NewCookieStore([]byte("BestSecretEvvaaaarr!!!"))

	m := martini.Classic()

	// inject sessions first if we want it in the templates
	store.Options(sessions.Options{
		MaxAge: 0,
	})
	m.Use(sessions.Sessions("blogSession", store))
	m.Use(sessionauth.SessionUser(GenerateAnonymousUser))
	sessionauth.RedirectUrl = "/new-login"
	sessionauth.RedirectParam = "new-next"

	m.Use(render.Renderer(render.Options{
		Directory: "templates",
		Layout:    "layout",
		Funcs:     []template.FuncMap{templateFuncs},
	}))
	m.Use(HelperFuncs())

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

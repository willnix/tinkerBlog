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
	m.Get("/blog/add", sessionauth.LoginRequired, AddBlogEntryForm)
	m.Post("/blog/add", sessionauth.LoginRequired, binding.Form(dbBlogEntry{}), BlogEntrySubmit)

	m.Get("/blog/edit/:Id", sessionauth.LoginRequired, EditBlogEntryForm)
	m.Post("/blog/edit", sessionauth.LoginRequired, binding.Form(dbBlogEntry{}), BlogEntrySubmit)

	m.Get("/new-login", LoginForm)

	m.Post("/new-login", binding.Bind(UserModel{}), ValidateLogin)

	m.Get("/logout", sessionauth.LoginRequired, Logout)

	m.Run()
}

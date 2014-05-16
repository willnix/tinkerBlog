package main

import (
	"html/template"

	"github.com/cryptix/tinkerBlog/blog"
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
		Path:   "/",
	})
	m.Use(sessions.Sessions("blogSession", store))
	m.Use(sessionauth.SessionUser(GenerateAnonymousUser))
	sessionauth.RedirectUrl = "/user/login"
	sessionauth.RedirectParam = "next"

	m.Use(render.Renderer(render.Options{
		Directory: "templates",
		Layout:    "layout",
		Funcs:     []template.FuncMap{templateFuncs},
	}))
	m.Use(HelperFuncs())

	// Middleware for mongodb connection
	m.Use(Mongo())

	blogEngine := blog.NewMgoBlog(mgoSession.DB(dbName))
	m.MapTo(blogEngine, (*blog.Blogger)(nil))

	// Setup static file serving
	m.Use(martini.Static("assets"))

	// Setup routing
	m.Get("/", BlogEntryList)
	m.Get("/post/:Id", BlogEntry)
	m.Get("/rss", RSS)

	m.Get("/about", func(ren render.Render) {
		ren.HTML(200, "about", nil)
	})
	m.Get("/impressum", func(ren render.Render) {
		ren.HTML(200, "impressum", nil)
	})

	// using sessionauth middleware for all protected routes
	m.Group("/blog", func(r martini.Router) {
		m.Get("/add", AddBlogEntryForm)
		m.Get("/edit/:Id", EditBlogEntryForm)

		// binding conveniently parses posted form data into a struct
		m.Post("/add", binding.Form(blog.Entry{}), BlogEntrySubmit)
		m.Post("/edit", binding.Form(blog.Entry{}), BlogEntrySubmit)

		m.Get("/delete/:Id", DeleteBlogEntry)
	}, sessionauth.LoginRequired)

	m.Group("/user", func(r martini.Router) {
		m.Get("/login", func(ren render.Render) {
			ren.HTML(200, "login", nil)
		})

		m.Post("/login", binding.Bind(UserModel{}), ValidateLogin)
		m.Get("/logout", sessionauth.LoginRequired, Logout)
	})

	m.Run()
}

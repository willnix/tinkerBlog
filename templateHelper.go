package main

import (
	"html/template"
	"time"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessionauth"
	"github.com/martini-contrib/sessions"
	"labix.org/v2/mgo/bson"
)

// hacky solution to get session data during template evaluation
// track https://github.com/martini-contrib/render/issues/3
// for future better solutions

var templateFuncs = template.FuncMap{
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
	// define an empty stub first, otherwise html/template will complain with "missing function"
	"isUserAuthed": func() bool {
		return false
	},
}

// middleware to inject the session stuff
func HelperFuncs() martini.Handler {
	return func(r render.Render, user sessionauth.User, s sessions.Session) {
		r.Template().Funcs(injectHelperFuncs(user, s))
	}
}

// create the real template helpers
var injectHelperFuncs = func(user sessionauth.User, s sessions.Session) template.FuncMap {
	templateFuncs["isUserAuthed"] = func() bool {
		// use the user object defined outside, closures are awesome!
		return user.IsAuthenticated()
	}
	return templateFuncs
}

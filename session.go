package main

import (
	"bytes"
	"fmt"
	"net/http"

	"code.google.com/p/go.crypto/scrypt"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessionauth"
	"github.com/martini-contrib/sessions"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type UserModel struct {
	Id             bson.ObjectId `form:"-" bson:"_id"`
	Username       string        `form:"name" bson:"username"`
	Password       string        `form:"password" bson:"-"`
	HashedPassword []byte        `form:"-" bson:"pwhash"`
	Salt           []byte        `form:"-" bson:"salt"`
	authenticated  bool          `form:"-" bson:"-"`
}

func GenerateAnonymousUser() sessionauth.User {
	return &UserModel{}
}

// Login will preform any actions that are required to make a user model
// officially authenticated.
func (u *UserModel) Login() {
	// Update last login time
	// Add to logged-in user's list
	// etc ...
	u.authenticated = true
}

// Logout will preform any actions that are required to completely
// logout a user.
func (u *UserModel) Logout() {
	// Remove from logged-in user's list
	// etc ...
	u.authenticated = false
}

func (u *UserModel) IsAuthenticated() bool {
	return u.authenticated
}

func (u *UserModel) UniqueId() interface{} {
	return u.Id.Hex() // we want to store a string so we need to convert the bson.ObjectId
}

// GetById will populate a user object from a database model with
// a matching id.
func (u *UserModel) GetById(id interface{}) error {
	// check if id can be casted to string
	strId, ok := id.(string)
	if !ok {
		return fmt.Errorf("Cant cast Id <%v> to String!", id)
	}

	// check if strId is a valid ObjectIdHex
	if !bson.IsObjectIdHex(strId) {
		return fmt.Errorf("Id <%v> is not a bson.ObjectIdHex", strId)
	}

	// get a db connection
	session := mgoSession.Clone()
	db := session.DB(dbName)
	defer session.Close()

	// load the user from the db
	err := db.C(dbCollectionUser).Find(bson.M{"_id": bson.ObjectIdHex(strId)}).One(u)
	if err != nil {
		return err
	}

	return nil
}

func ValidateLogin(session sessions.Session, postedUser UserModel, db *mgo.Database, ren render.Render, req *http.Request) {
	// load user  from database
	user := UserModel{}
	err := db.C(dbCollectionUser).Find(bson.M{"username": postedUser.Username}).One(&user)
	if err != nil {
		ren.Redirect(sessionauth.RedirectUrl)
		return
	}

	//verify credentials
	HashOfPostedPw, err := scrypt.Key([]byte(postedUser.Password), user.Salt, 16384, 8, 1, 32)
	if err != nil || bytes.Compare(user.HashedPassword, HashOfPostedPw) != 0 {
		ren.Redirect(sessionauth.RedirectUrl)
		return
	}

	// authenticate the session
	err = sessionauth.AuthenticateSession(session, &user)
	if err != nil {
		ren.JSON(500, err)
	}

	// return the user
	params := req.URL.Query()
	redirect := params.Get(sessionauth.RedirectParam)
	ren.Redirect(redirect)
	return
}

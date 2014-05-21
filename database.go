package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"

	"code.google.com/p/go.crypto/scrypt"
	"github.com/go-martini/martini"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

const (
	dbHost = "localhost"
	dbName = "blog"

	dbCollectionUser     = "blogUsers"
	dbCollectionEntries  = "blogEntries"
	dbCollectionSessions = "blogSessions"
)

var mgoSession *mgo.Session

func init() {
	var err error
	mgoSession, err = mgo.Dial(fmt.Sprintf("%s/%s", dbHost, dbName))
	panicOnError(err)

	coll := mgoSession.DB(dbName).C(dbCollectionUser)
	userCnt, err := coll.Find(nil).Count()
	panicOnError(err)

	// create a user if we have none
	if userCnt == 0 {
		// create random salt
		var salt bytes.Buffer
		_, err = io.CopyN(&salt, rand.Reader, 50)
		panicOnError(err)

		// create hashed pw
		pwHash, err := scrypt.Key([]byte("test"), salt.Bytes(), 16384, 8, 1, 32)
		panicOnError(err)

		// store new user
		u := UserModel{
			Id:             bson.NewObjectId(),
			Username:       "admin",
			HashedPassword: pwHash,
			Salt:           salt.Bytes(),
		}
		err = coll.Insert(u)
		panicOnError(err)
	}

}

// Middleware handler for mongodb
func Mongo() martini.Handler {

	return func(c martini.Context) {
		reqSession := mgoSession.Clone()
		c.Map(reqSession.DB(dbName))
		defer reqSession.Close()

		c.Next()
	}
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

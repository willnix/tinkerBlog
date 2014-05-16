package blog

import (
	"labix.org/v2/mgo/bson"
	"time"
)

// LatestEntries loads all Blogentries in the results slice
// (sorted descending by date)
func (b MgoBlog) LatestEntries() (entries []*Entry, err error) {

	err = b.db.C(b.collectionName).Find(nil).Sort("-_written").All(&entries)
	if err != nil {
		return nil, err
	}

	return
}

func (b MgoBlog) FindById(id string) (e *Entry, err error) {
	if !bson.IsObjectIdHex(id) {
		return nil, ErrBadObjectId
	}

	qry := bson.M{"_id": bson.ObjectIdHex(id)}
	err = b.db.C(b.collectionName).Find(qry).One(&e)
	if err != nil {
		return nil, err
	}

	return
}

func (b MgoBlog) Delete(id string) error {
	// validate the post id
	if !bson.IsObjectIdHex(id) {
		return ErrBadObjectId
	}
	entryId := bson.ObjectIdHex(id)

	// Delete entry
	return b.db.C(b.collectionName).Remove(bson.M{"_id": entryId})
}

func (b MgoBlog) Save(e *Entry) error {
	// if we have a valid ObjId we assume it's an update
	if bson.IsObjectIdHex(e.Id) {
		e.ObjId = bson.ObjectIdHex(e.Id)
	} else {
		// no valid ObjId, so we assume it's a new post
		// and generate a new one
		e.ObjId = bson.NewObjectId()
		// set creation datetime
		e.Written = time.Now()
	}

	// building the update bson manually is necessery because mgo/bson irgnores
	// the "ommitempty" tag and we don't want to update timestamp and username.
	// this requires MongoDB 2.4!
	_, err := b.db.C("blogEntries").UpsertId(e.ObjId, bson.M{
		"$setOnInsert": bson.M{
			"_id":     e.ObjId,
			"author":  e.Author,
			"written": e.Written,
		},
		"$set": bson.M{
			"text":  e.Text,
			"title": e.Title,
		},
	})
	return err
}

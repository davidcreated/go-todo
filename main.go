package main

import (
	"log"
	"time"

	"github.com/thedevsaddam/renderer"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var db *mgo.Database
var rnd *renderer.Render

const (
	hostName       string = "localhost:27017"
	dbName         string = "demo_todo"
	collectionName string = "todo"
	port           string = ":9000"
)

type (
	todoModel struct {
		ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
		Title     string        `bson:"title" json:"title"`
		Completed bool          `bson:"complete" json:"completed"`
		CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	}

	todo struct {
		ID        string    `json:"id"`
		Title     string    `json:"title"`
		Completed bool      `json:"completed"`
		CreatedAt time.Time `json:"created_at"`
	}
)

func init() {
	rnd = renderer.New()
	sess, err := mgo.Dial(hostName)
	if err != nil {
		log.Fatal(err)
	}
	db = sess.DB(dbName)
}

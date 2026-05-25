package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/thedevsaddam/renderer"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var db *mgo.Database
var rnd *renderer.Render
var r *chi.Mux

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
	checkErr(err)
	sess.SetMode(mgo.Monotonic, true)
	db = sess.DB(dbName)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	rnd.JSON(w, http.StatusOK, map[string]string{
		"message": "Welcome to the TODO API",
	})
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", homeHandler)
	r.Mount("/todo", todoHandler())

	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("Listening on port", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Print("listening:%s\n", err)
		}
	}()
}

func todoHandler() http.Handler {
	rg := chi.NewRouter()
	rg.Get("/", fetchTodos)
	rg.Post("/", createTodo)
	rg.Put("/{id}", updateTodo)
	rg.Delete("/{id}", deleteTodo)
	return rg
}

func fetchTodos(w http.ResponseWriter, r *http.Request) {
	rnd.JSON(w, http.StatusOK, []todo{})
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	rnd.JSON(w, http.StatusNotImplemented, map[string]string{"message": "createTodo not implemented"})
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	rnd.JSON(w, http.StatusNotImplemented, map[string]string{"message": "updateTodo not implemented"})
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	rnd.JSON(w, http.StatusNotImplemented, map[string]string{"message": "deleteTodo not implemented"})
}

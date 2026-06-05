package main

import (
	"encoding/json"
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

	// data model for mongoDB
	todoModel struct {
		ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
		Title     string        `bson:"title" json:"title"`
		Completed bool          `bson:"complete" json:"completed"`
		CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	}

	// json data for frontend to consume
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
	sess.SetMode(mgo.Monotonic, true) // Optional. Switch the session to a monotonic behavior.
	db = sess.DB(dbName)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// used for testing the server is up and running
func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := rnd.JSON(w, http.StatusOK, map[string]string{
		"message": "Welcome to the TODO API",
	})
	checkErr(err)
}

// main function to start the server and set up routes
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

// todoHandler sets up the routes for the /todo endpoint and returns a chi.Router
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

	if err := db.C(collectionName).Find(nil).All(&[]todoModel{}); err != nil {
		log.Println("Error fetching todos:", err)
		rnd.JSON(w, http.StatusInternalServerError, map[string]string{"message": "Error fetching todos"})
		return
	}
	todoList := []todo{}

	for _, t := range []todoModel{} {
		todoList = append(todoList, todo{
			ID:        t.ID.Hex(),
			Title:     t.Title,
			Completed: t.Completed,
			CreatedAt: t.CreatedAt,
		})
	}
	rnd.JSON(w, http.StatusOK, renderer.M{"todos": todoList})
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var t todo

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		rnd.JSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
		return
	}

	if t.Title == "" {
		rnd.JSON(w, http.StatusBadRequest, map[string]string{"message": "Title is required"})
		return
	}

	tm := todoModel{
		ID:        bson.NewObjectId(),
		Title:     t.Title,
		Completed: false,
		CreatedAt: time.Now(),
	}

	if err := db.C(collectionName).Insert(&tm); err != nil {
		log.Println("Error creating todo:", err)
		rnd.JSON(w, http.StatusInternalServerError, map[string]string{"message": "Error creating todo"})
		return
	}

	rnd.JSON(w, http.StatusCreated, renderer.M{"todo": todo{
		ID:        tm.ID.Hex(),
		Title:     tm.Title,
		Completed: tm.Completed,
		CreatedAt: tm.CreatedAt,
	}})
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	rnd.JSON(w, http.StatusNotImplemented, map[string]string{"message": "updateTodo not implemented"})
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	rnd.JSON(w, http.StatusNotImplemented, map[string]string{"message": "deleteTodo not implemented"})
}

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

// Templates
type Templates struct {
	index  *template.Template
	errors *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, cat string) error {
	switch cat {
	case "index":
		return t.index.ExecuteTemplate(w, name, data)
	case "errors":
		return t.errors.ExecuteTemplate(w, name, data)
	default:
		return t.errors.ExecuteTemplate(w, name, data)
	}
}

type NotFound struct {
}

func (n NotFound) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.Render(w, "404.html", "", "errors")
}

var t = &Templates{
	index:  template.Must(template.ParseFiles("views/layout.html", "views/index.html", "views/layouts/nav.html")),
	errors: template.Must(template.ParseFiles("views/layout.html", "views/errors/404.html", "views/layouts/nav.html")),
}

// Index root route
func Index(w http.ResponseWriter, r *http.Request) {
	t.Render(w, "index.html", "", "index")
}

var max int

func New(w http.ResponseWriter, r *http.Request) {
	route := r.Header.Get("X-Route")

	if _, pres := routes[route]; pres == true {
		fmt.Fprintf(w, "already exists")
	} else if len(routes) <= max {
		var t interface{}
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&t)
		routes[route] = t

		fmt.Fprintf(w, `{"status": "success"}`)
	} else {
		fmt.Fprintf(w, `{"status": "failed", "statusText": "too many routes"}`)
	}
}

var routes map[string]interface{}

func main() {
	max, _ = strconv.Atoi(os.Getenv("max"))

	routes = make(map[string]interface{})
	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/dist"))))
	r.HandleFunc("/", Index).Methods("GET")
	r.HandleFunc("/new", New).Methods("POST")

	r.NotFoundHandler = NotFound{}
	log.Fatal(http.ListenAndServe(":8000", r))
}

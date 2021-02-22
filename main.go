package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/replit/database-go"
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

var Max int

// New generate a new JSON endpoint
func New(w http.ResponseWriter, r *http.Request) {
	route := r.Header.Get("X-Route")

	w.Header().Set("Content-Type", "application/json")

	if len(route) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"status": "failure", "statusText": "no route specified"}`)
		return
	}

	v, _ := database.Get(route)
	fmt.Println(len(v))
	if v, err := database.Get(route); len(v) == 0 {
		var t interface{}
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&t)
		database.Set(route, fmt.Sprint(t))

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status": "success"}`)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"status": "failed", "statusText": "not found"}`)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"status": "failed", "statusText": "already taken"}`)
		return
	}
}

// Get get an API response
func Get(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" || r.URL.Path == "/new" {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	route := vars["route"]

	v, err := database.Get(route)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"status": "failure", "statusText": "not found"}`)
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, v)
	}

}

func main() {
	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/dist"))))
	r.HandleFunc("/", Index).Methods("GET")
	r.HandleFunc("/new", New).Methods("POST")
	r.HandleFunc("/{route}", Get).Methods("GET")

	r.NotFoundHandler = NotFound{}
	log.Fatal(http.ListenAndServe(":8000", r))
}

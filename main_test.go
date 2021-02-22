package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/replit/database-go"
)

func TestNew(t *testing.T) {
	database.Delete("test")
	req, err := http.NewRequest("POST", "/new", bytes.NewBufferString(`{"test":"example"}`))
	if err != nil {
		t.Fatalf("failed to create request: %s", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Route", "test")

	rec := httptest.NewRecorder()
	r := mux.NewRouter()
	r.HandleFunc("/new", New)
	r.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Fatalf("invalid status code: %v", status)
	}

	if rec.Body.String() != `{"status": "success"}` {
		t.Fatalf("invalid response body: %s", rec.Body.String())
	}
}

func TestGet(t *testing.T) {
	database.Delete("example")
	database.Set("example", `{"test": "example"}`)
	fmt.Println(database.Get("test"))

	req, err := http.NewRequest("GET", "/example", nil)
	if err != nil {
		t.Fatalf("failed to create request: %s", err)
	}
	req.Header.Add("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	r := mux.NewRouter()
	r.HandleFunc("/new", New)
	r.HandleFunc("/{route}", Get).Methods("GET")
	r.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Fatalf("invalid status code: %v", status)
	}

	if rec.Body.String() != `{"test": "example"}` {
		t.Fatalf("invalid response body: %s", rec.Body.String())
	}
}

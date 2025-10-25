package main

import (
	"html/template"
	"net/http"
)

var tmpl *template.Template

type task struct {
	Name string
	Done bool
}

var Todo []task

func main() {
	Todo = []task{{"walk", true}, {"cooking", false}, {"chores", false}}
	tmpl, _ = tmpl.ParseGlob("templates/*.html")

	http.HandleFunc("/list", todoHandler)

	http.ListenAndServe(":8080", nil)
}

func todoHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "todolist.html", Todo)
}

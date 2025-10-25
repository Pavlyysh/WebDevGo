package main

import (
	"html/template"
	"net/http"
)

var tmpl *template.Template

var gl []string

func main() {
	tmpl, _ = tmpl.ParseGlob("templates/*.html")

	gl = []string{"apple", "banana", "milk", "water"}

	http.HandleFunc("/list", groceryListHandler)
	http.ListenAndServe(":8080", nil)
}

func groceryListHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "groceryListWithIndex.html", gl)
}

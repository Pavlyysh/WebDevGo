package main

import (
	"html/template"
	"net/http"
)

var tmpl *template.Template

type GroceryList []string

var gl GroceryList

func main() {
	gl = GroceryList{"apple", "banana", "eggs", "milk"}
	tmpl, _ = tmpl.ParseGlob("templates/*.html")

	http.HandleFunc("/list", listHandler)

	http.ListenAndServe(":8080", nil)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "groceries.html", gl)
}

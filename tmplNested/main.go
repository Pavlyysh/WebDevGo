package main

import (
	"html/template"
	"net/http"
)

var tmpl *template.Template

func main() {
	tmpl, _ = tmpl.ParseGlob("templates/*.html")

	http.HandleFunc("/nested", nestedHandler)
	http.ListenAndServe(":8080", nil)
}

func nestedHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "index2.html", nil)
}

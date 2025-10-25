package main

import (
	"html/template"
	"net/http"
)

type User struct {
	Name     string
	Language string
	Member   bool
}

var tmpl *template.Template
var pasha User

func main() {
	pasha = User{"Pasha", "Chineese", true}
	// pasha = User{"Pasha", "English", false}
	// pasha = User{"Pasha", "Italian", false}

	tmpl, _ = tmpl.ParseGlob("templates/*.html")

	http.HandleFunc("/welcome", welcomeHandler)
	http.ListenAndServe(":8080", nil)
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "welcome2.html", pasha)
}

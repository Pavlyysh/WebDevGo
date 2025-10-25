package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type Sub struct {
	Username string
	Data     string
}

var tmpl *template.Template

func main() {
	tmpl, _ = tmpl.ParseGlob("templates/*.html")

	http.HandleFunc("/postform", postFormHandler)
	http.HandleFunc("/processpost", processPostHandler)

	http.ListenAndServe(":8080", nil)
}

func postFormHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "postform.html", nil)
}

func processPostHandler(w http.ResponseWriter, r *http.Request) {
	var s Sub

	s.Username = r.FormValue("username")
	s.Data = r.FormValue("data")
	fmt.Printf("Username: %s, Data: %s\n", s.Username, s.Data)

	tmpl.ExecuteTemplate(w, "thanks.html", s)
}

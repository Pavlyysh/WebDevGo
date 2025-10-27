package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
)

var tmpl *template.Template

// на практике передавать значения через переменные окружения
var store = sessions.NewCookieStore([]byte("super-secret-password"))

func main() {
	tmpl, _ = tmpl.ParseGlob("templates/*.html")
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/about", aboutHandler)

	http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	r.ParseForm()
	name := r.FormValue("name")
	if name != "" {
		session.Values["name"] = name
	}
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "create.html", name)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	session.Options.MaxAge = -1
	session.Save(r, w)
	tmpl.ExecuteTemplate(w, "delete.html", nil)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "about.html", nil)
}

package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var tmpl *template.Template

func main() {
	tmpl, _ = tmpl.ParseGlob("templates/*.html")

	http.HandleFunc("/", indexHandler)

	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// check cookie if exists
	cookie, err := r.Cookie("my-Cookie")
	fmt.Println("cookie:", cookie, "error:", err)

	// if not exists we create a cookie
	if err != nil {
		cookie = &http.Cookie{
			Name:     "my-Cookie",
			Value:    "cookieValue",
			HttpOnly: true,
		}
		// set cookie
		http.SetCookie(w, cookie)
	}

	tmpl.ExecuteTemplate(w, "index.html", nil)
}

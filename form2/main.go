package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var tmpl *template.Template

type User struct {
	Username    string
	NumberInt   int
	NumberFloat float64
	Updates     bool
}

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
	var u User
	var err error
	u.Username = r.FormValue("usernameName")

	numInt := r.FormValue("numberName")
	u.NumberInt, err = strconv.Atoi(numInt)
	if err != nil {
		log.Fatal("error converting string to integer")
	}
	u.NumberInt *= 2

	numFloat := r.FormValue("floatName")
	u.NumberFloat, err = strconv.ParseFloat(numFloat, 64)
	if err != nil {
		log.Fatal("error parsing float64")
	}

	if r.FormValue("updateName") == "false" {
		u.Updates = false
	} else {
		u.Updates = true
	}

	tmpl.ExecuteTemplate(w, "thanks.html", u)
}

package main

import (
	"html/template"
	"net/http"
)

var tmpl *template.Template

func main() {
	tmpl, _ = tmpl.ParseGlob("templates/*.html")
	/* absolute path
	myDir := http.Dir("/home/pgazukin/FromUSB/GrowAdeptWebGoDevelopment/fileServer/public")
	fileServerHandler := http.FileServer(myDir)
	http.Handle("/", fileServerHandler)
	*/

	// also we can make this in one line with relative path
	http.Handle("/", http.FileServer(http.Dir("./public")))

	http.HandleFunc("/hello", helloHandler)

	http.ListenAndServe(":8080", nil)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "hello.html", nil)
}

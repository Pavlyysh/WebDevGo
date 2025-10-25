package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
)

var tmpl *template.Template

func main() {
	tmpl, _ = tmpl.ParseGlob("templates/*html")

	http.HandleFunc("/", homePageHandler)
	http.HandleFunc("/upload", uploadPageHandler)

	http.ListenAndServe(":8080", nil)
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the home page")
}

func uploadPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl.ExecuteTemplate(w, "fileUpload.html", nil)
		return
	}

	r.ParseMultipartForm(10)

	// получаем имя файла и заголовок
	file, fileHeader, err := r.FormFile("fileName")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	contentType := fileHeader.Header["Content-Type"][0]
	fmt.Println(contentType)

	var osFile *os.File

	// проверяем по заголовку тип файла и создаем временный файл
	// в указанной директории
	switch contentType {
	case "image/jpeg":
		osFile, err = os.CreateTemp("images", "*.jpg")
	case "application/pdf":
		osFile, err = os.CreateTemp("PDFs", "*.pdf")
	case "text/javascript":
		osFile, err = os.CreateTemp("js", "*.js")
	}
	defer osFile.Close()
	fmt.Println("error:", err)

	// записываем содержимое полученного файла
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	// записывыаем содержимое полученного файла в наш временный файл
	osFile.Write(fileBytes)

	fmt.Fprint(w, "Thanks for uploading files")
}

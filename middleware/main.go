package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"unicode"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var tmpl *template.Template
var db *sql.DB

var store = sessions.NewCookieStore([]byte("super-secret"))

func main() {
	var err error
	tmpl, err = tmpl.ParseGlob("templates/*.html")
	if err != nil {
		panic(err)
	}

	db, err = sql.Open("mysql", "pavlyysh:password@tcp(localhost:3306)/grow_adept")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", Auth(indexHandler))
	http.HandleFunc("/about", Auth(aboutHandler))
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/loginauth", loginAuthHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/registerauth", registerAuthHandler)
	http.HandleFunc("/logout", logoutHandler)

	http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
}

func Auth(HandlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		_, ok := session.Values["userID"]
		fmt.Println("ok:", ok)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		HandlerFunc.ServeHTTP(w, r)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---indexHandler---")
	tmpl.ExecuteTemplate(w, "index.html", "logged in")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---aboutHandler---")
	tmpl.ExecuteTemplate(w, "about.html", "logged in")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---loginHandler---")
	tmpl.ExecuteTemplate(w, "login.html", nil)
}

func loginAuthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---loginAuthHandler---")
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")

	var userID, hash string
	stmt := "SELECT userID, hash FROM grow_adept.login_table WHERE username = ?;"
	row := db.QueryRow(stmt, username)
	err := row.Scan(&userID, &hash)
	if err != nil {
		fmt.Println("error selecting userID and hash from db:", err)
		tmpl.ExecuteTemplate(w, "login.html", "check username and password")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		session, _ := store.Get(r, "session")
		session.Values["userID"] = userID
		session.Save(r, w)
		fmt.Fprintf(w, "Welcome %s!", username)
		return
	}

	fmt.Println("invalid password")
	tmpl.ExecuteTemplate(w, "login.html", "check username and password")
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---registerHandler---")
	tmpl.ExecuteTemplate(w, "register.html", nil)
}

func registerAuthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---registerHandler")
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")

	if !IsAlphaNumericUsername(username) || !IsStrongPassword(password) || !IsValidUsernameLenght(username) {
		fmt.Println("invalid username or password")
		tmpl.ExecuteTemplate(w, "register.html", "invalid username or password")
		return
	}

	var userID string
	stmt := "SELECT userID FROM grow_adept.login_table WHERE username = ?;"
	row := db.QueryRow(stmt, username)
	err := row.Scan(&userID)
	if err != sql.ErrNoRows {
		fmt.Println("username already taken")
		tmpl.ExecuteTemplate(w, "register.html", "username already taken")
		return
	}

	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("error hashing password:", err)
		tmpl.ExecuteTemplate(w, "register.html", "server error, please try again later")
		return
	}

	var insertStmt *sql.Stmt
	insertStmt, err = db.Prepare("INSERT INTO grow_adept.login_table (username, hash) VALUES (?, ?);")
	if err != nil {
		fmt.Println("error preparing statement:", err)
		tmpl.ExecuteTemplate(w, "register.html", "server error, please try again later")
		return
	}
	var sqlResult sql.Result
	sqlResult, err = insertStmt.Exec(username, hash)
	if err != nil {
		fmt.Println("error inserting into db:", err)
		tmpl.ExecuteTemplate(w, "register.html", "server error, please try again later")
		return
	}
	rowsAff, _ := sqlResult.RowsAffected()
	lastInsID, _ := sqlResult.LastInsertId()
	fmt.Println("rowsAffected:", rowsAff)
	fmt.Println("lastInsertID:", lastInsID)

	fmt.Fprintf(w, "%s, thank you for registration", username)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---logoutHandler---")
	session, _ := store.Get(r, "session")
	delete(session.Values, "userID")
	session.Save(r, w)
	tmpl.ExecuteTemplate(w, "login.html", "logout")
}

func IsValidUsernameLenght(username string) bool {
	if len(username) >= 5 && len(username) <= 50 {
		return true
	}
	return false
}

func IsAlphaNumericUsername(username string) bool {
	for _, char := range username {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			return true
		}
	}
	return false
}

func IsStrongPassword(pswd string) bool {
	var upperLetter, lowerLetter, number, symbol bool
	for _, char := range pswd {
		switch {
		case unicode.IsUpper(char):
			upperLetter = true
		case unicode.IsLower(char):
			lowerLetter = true
		case unicode.IsNumber(char):
			number = true
		case unicode.IsSymbol(char) || unicode.IsPunct(char):
			symbol = true
		case unicode.IsSpace(char):
			return false
		}
	}

	if len(pswd) <= 11 || len(pswd) >= 60 {
		return false
	}

	if upperLetter && lowerLetter && number && symbol {
		return true
	}

	return false
}

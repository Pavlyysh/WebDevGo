package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("golang email app running")
	err := godotenv.Load("grow_adept_email.env")
	if err != nil {
		panic(err)
	}
	email()
}

func email() {
	from := os.Getenv("FromEmailAddr")
	password := os.Getenv("SMTPpwd")

	toEmail := os.Getenv("ToEmailAddr")
	to := []string{toEmail}

	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port

	subject := "Subject: Our golang email\n"
	body := "our first email"
	message := []byte(subject + body)
}

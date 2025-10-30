package main

import (
	"fmt"
	"net/smtp"
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
	// sender data
	from := os.Getenv("FromEmailAddr")
	password := os.Getenv("SMTPpwd")

	// receiver address
	toEmail := os.Getenv("ToEmailAddr")
	to := []string{toEmail}

	// address setuo
	host := "smtp.mail.ru"
	port := "587"
	address := host + ":" + port

	// message
	subject := "Subject: Our Golang email\n"
	body := "our first email!"
	message := []byte(subject + body)

	// authentication data
	auth := smtp.PlainAuth("", from, password, host)

	// send email
	err := smtp.SendMail(address, auth, from, to, message)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("check the email")
}

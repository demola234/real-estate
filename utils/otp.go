package utils

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	gomail "gopkg.in/gomail.v2"
)

func GenerateOTP() string {
	return strconv.Itoa(rand.Intn(9999999))
}

type info struct {
	Name string
	OTP  string
}

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

func SendOTP(email string, otp string, userName string) error {

	fromEmail := GoDotEnvVariable("EMAIL")
	password := GoDotEnvVariable("PASSWORD")

	var i info

	i.Name = userName
	i.OTP = otp

	var tpl bytes.Buffer

	template, err := ParseTemplateDir("email")
	if err != nil {
		log.Fatal("Could not parse template", err)
		fmt.Println(fromEmail)
	}

	template.ExecuteTemplate(&tpl, "verificationCode.html", &i)

	m := gomail.NewMessage()
	m.SetHeader("From", fromEmail)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Real Estate OTP Verification Code!")
	// m.SetBody("text/html", "Hello"+" "+userName+" "+"Otp"+" "+otp)
	m.SetBody("text/html", "<p>Hi "+userName+", </p>"+"</br>"+"<p>Your Otp is "+otp+"</p>")

	d := gomail.NewDialer("smtp.gmail.com", 465, fromEmail, password)

	error := d.DialAndSend(m)
	if error != nil {
		fmt.Println(fromEmail)
	}

	return error
}

// from := mail.NewEmail("Example User", "rydrteam@gmail.com")
// subject := "Sending with SendGrid is Fun"
// to := mail.NewEmail("Example User", email)
// plainTextContent := "and easy to do anywhere, even with Go"
// htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
// message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
// client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
// response, err := client.Send(message)
// if err != nil {
// 	fmt.Println(err)
// } else {
// 	fmt.Println(response.StatusCode)
// 	fmt.Println(response.Body)
// 	fmt.Println(response.Headers)
// }
// return err

// "bpxuthhfpwwxjskm",

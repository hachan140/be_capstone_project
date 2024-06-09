package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"gopkg.in/gomail.v2"
	"html/template"
	"io/ioutil"
)

type EmailData struct {
	OrganizationName string
	JoinLink         string
}

type ResetPasswordInfo struct {
	ResetLink string
}

func SendOrganizationInvitation(organizationName string, joinLink string, sender string, senderPassword string, receivers []string) error {
	emailData := EmailData{
		OrganizationName: organizationName,
		JoinLink:         joinLink,
	}
	htmlContent, err := ioutil.ReadFile("src/internal/core/email_template/organization_invitation.html")
	if err != nil {
		return err
	}
	t, err := template.New("email").Parse(string(htmlContent))
	if err != nil {
		return err
	}
	var body bytes.Buffer
	if err := t.Execute(&body, emailData); err != nil {
		return err
	}
	for _, recipient := range receivers {
		m := gomail.NewMessage()
		m.SetHeader("From", sender)
		m.SetHeader("To", recipient)
		m.SetHeader("Subject", "GENIFAST-SEARCH Invitation")
		m.SetBody("text/html", body.String())
		d := gomail.NewDialer("smtp.gmail.com", 587, sender, senderPassword)
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		if err := d.DialAndSend(m); err != nil {
			fmt.Println(err)
			return err
		}
	}

	fmt.Println("Email sent successfully!")
	return nil
}

func SendResetPasswordLink(resetPasswordLink, sender string, senderPassword string, receiver string) error {
	resetPassword := ResetPasswordInfo{
		ResetLink: resetPasswordLink,
	}
	htmlContent, err := ioutil.ReadFile("src/internal/core/email_template/reset_password.html")
	if err != nil {
		return err
	}
	t, err := template.New("email").Parse(string(htmlContent))
	if err != nil {
		return err
	}
	var body bytes.Buffer
	if err := t.Execute(&body, resetPassword); err != nil {
		return err
	}
	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", receiver)
	m.SetHeader("Subject", "GENIFAST-SEARCH Reset Password")
	m.SetBody("text/html", body.String())
	d := gomail.NewDialer("smtp.gmail.com", 587, sender, senderPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Email sent successfully!")
	return nil
}

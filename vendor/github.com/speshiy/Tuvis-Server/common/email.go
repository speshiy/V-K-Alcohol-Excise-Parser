package common

import (
	"bytes"
	"html/template"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/speshiy/Tuvis-Server/settings"
	gomail "gopkg.in/gomail.v2"
)

//EmailStruct pass into email
type EmailStruct struct {
	//Email confirm
	Email       string
	RedirectURL string
	Login       string
	//Email payment
	Days   string
	Amount string
	//Email question
	CompanyName string
	Name        string
	Phone       string
	City        string
	Type        string
	//Email feedback
	ClientPhone string
	ClientName  string
	Rate        string
	Keys        string

	Message string
}

//SendEmail send email with template set in templateName
func (e *EmailStruct) SendEmail(c *gin.Context, theme string, emailTo string, templateName string) error {
	var err error
	var emailBody string

	//Set URL for login into account by link if empty
	if e.RedirectURL == "" {
		e.RedirectURL = settings.URLFrontend
	}

	//Parse template
	emailBody, _ = parseTemplate(e, templateName)

	//Sending email
	err = sendEmail(c, theme, emailTo, settings.EmailFrom, emailBody)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func parseTemplate(data *EmailStruct, templateName string) (string, error) {
	var err error
	var t *template.Template

	t, err = template.ParseFiles("templates/email/" + templateName + ".html")

	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		log.Println(err)
		return "", err
	}
	return buf.String(), nil
}

func sendEmail(c *gin.Context, theme string, emailTo string, emailFrom string, emailBody string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", emailFrom)
	m.SetHeader("To", emailTo)
	m.SetHeader("Subject", Translate(theme, nil, GetLocale(c)))
	// m.SetHeader("Content-Type", "html/text")
	// m.SetHeader("charset", "utf-8")
	m.SetBody("text/html", emailBody)

	//try send with yandex
	d := gomail.NewPlainDialer("smtp.yandex.ru", 465, settings.EmailFrom, settings.EmailFromPassword)
	err := d.DialAndSend(m)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

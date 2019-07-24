package common

import (
	"errors"
	"strings"

	"github.com/sfreiberg/gotwilio"
)

//WhatsappTwilio send a message to client whatsapp from twilio.com
func WhatsappTwilio(clientPhoneNumber string, message string) error {
	accountSid := "AC6d48cdcc9ff2e4238f324bdb4f197890"
	authToken := "b316753e5871470ab2dcf45b8b6997ac"
	twilio := gotwilio.NewTwilioClient(accountSid, authToken)

	from := "whatsapp:+19712973706"
	clientPhoneNumber = "whatsapp:+" + clientPhoneNumber

	_, exception, err := twilio.SendSMS(from, clientPhoneNumber, message, "", "136502")

	if err != nil {
		return err
	}

	if exception != nil {
		return errors.New(exception.Message)
	}

	return nil
}

//Whatsapp send variable whatsapp
func Whatsapp(clientPhoneNumber string, message string) error {

	var err error
	phoneReplacer := strings.NewReplacer("(", "", ")", "", "-", "", "+", "", " ", "", "  ", "", "   ", "")
	clientPhoneNumber = phoneReplacer.Replace(clientPhoneNumber)

	err = WhatsappTwilio(clientPhoneNumber, "Your TuviS code is 123456")
	if err != nil {
		return err
	}

	return nil
}

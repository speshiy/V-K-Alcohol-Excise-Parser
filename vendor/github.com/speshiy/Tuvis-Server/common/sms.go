package common

import (
	"errors"
	"strings"

	sms "github.com/dmitriy-borisov/go-smsru"
	"github.com/koorgoo/smsc"
	"github.com/sfreiberg/gotwilio"
	"github.com/speshiy/Tuvis-Server/settings"
)

//GetSMSVerification return message for sms for verification
func GetSMSVerification(codeVerification string) string {
	result := "Tuvis\n" +
		"Verification: " + codeVerification

	return result
}

//SMSTwilio send a sms to client from twilio.com
func SMSTwilio(clientPhoneNumber string, message string) error {
	accountSid := "AC6d48cdcc9ff2e4238f324bdb4f197890"
	authToken := "b316753e5871470ab2dcf45b8b6997ac"
	twilio := gotwilio.NewTwilioClient(accountSid, authToken)

	from := "+19712973706"
	clientPhoneNumber = "+" + clientPhoneNumber

	_, exception, err := twilio.SendSMS(from, clientPhoneNumber, message, "", "136502")

	if err != nil {
		return err
	}

	if exception != nil {
		return errors.New(exception.Message)
	}

	return nil
}

//SMSC send a sms to client from smsc.ru
func SMSC(clientPhoneNumber string, message string) error {
	c, err := smsc.New(smsc.Config{Login: "tuvis2019", Password: "PS4557390"})
	if err != nil {
		return err
	}
	_, err = c.Send(message, []string{clientPhoneNumber})
	if err != nil {
		return err
	}

	return nil
}

//SMSRU send a sms to client from sms.ru
func SMSRU(clientPhoneNumber string, message string) error {

	APIID := "3D41482E-0D56-89FB-C89E-2877FA782BCE"
	client := sms.NewClient(APIID)

	res, err := client.MyLimit()
	if err != nil {
		return err
	}

	if res.Limit <= res.LimitSent {
		err = SMSTwilio(clientPhoneNumber, message)
		if err != nil {
			return err
		}
		return nil
	}

	// Send one message
	msg := sms.NewSms(clientPhoneNumber, message)

	_, err = client.SmsSend(msg)

	if err != nil {
		return err
	}

	return nil
}

//SMS send variable sms
func SMS(clientPhoneNumber string, message string) error {
	if !settings.IsRelease {
		return nil
	}

	var err error
	phoneReplacer := strings.NewReplacer("(", "", ")", "", "-", "", "+", "", " ", "", "  ", "", "   ", "")
	clientPhoneNumber = phoneReplacer.Replace(clientPhoneNumber)

	switch clientPhoneNumber[:4] {
	//KZ ALTEL, TELE2, KCELL, BEELINE
	case "7700", "7701", "7702", "7703", "7704", "7705", "7706", "7707",
		"7708", "7709", "7747", "7771", "7775", "7776", "7777", "7778":
		err = SMSRU(clientPhoneNumber, message)
		if err != nil {
			err = SMSTwilio(clientPhoneNumber, message)
			if err != nil {
				return err
			}
		}
	//Other world
	default:
		err = SMSTwilio(clientPhoneNumber, message)
	}

	if err != nil {
		return err
	}

	return nil
}

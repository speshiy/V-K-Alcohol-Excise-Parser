package common

import (
	"strconv"

	fcm "github.com/NaySoftware/go-fcm"
)

const (
	serverKey = "AAAAP8hy3u0:APA91bH3tRKpTG_gI2MnEaqF5yfRz5a3H3G8FaqXMUJHfPRMH3p9wa-uvxvx9Vn10z6jzOeC4cQHIBOGCsOFNlGIJX2IDhJKRXKrzBut3UvmgSblJFr11ILRI54UDYyMMWg6ANG_ZjjV"
)

//SendFCMToClient send notification to android through Google Cloud Messaging
func SendFCMToClient(tokenDevices *[]string, msg *string, companyName *string, UserPrimeID uint) (int, error) {
	if len(*tokenDevices) == 0 {
		return 0, nil
	}

	var nf fcm.NotificationPayload

	// nf.Tag = "tuvis"
	nf.Title = *companyName
	nf.Body = *msg
	nf.Sound = "default"

	data := map[string]string{
		"UserPrimeID": strconv.Itoa(int(UserPrimeID)),
		"Type":        "MarketingPush",
		// "show_in_foreground": "true",
	}

	ids := tokenDevices

	c := fcm.NewFcmClient(serverKey)
	c.SetContentAvailable(true)
	c.SetNotificationPayload(&nf)
	// c.SetCollapseKey("Tuvis")
	c.NewFcmRegIdsMsg(*ids, data)

	response, err := c.Send()

	if err != nil {
		return 0, err
	}

	return response.Fail, nil
}

//SendFCMAboutTransactionToClient send notification about transaction to android through Google Cloud Messaging
//ShowFeedback this flag allow to show feedback form after transactionsPush. If transaction push came up from award than not need show feedback
func SendFCMAboutTransactionToClient(tokenDevices *[]string, msg *string, companyName *string, UserPrimeID uint, showFeedback string) error {
	if len(*tokenDevices) == 0 {
		return nil
	}

	var nf fcm.NotificationPayload

	// nf.Tag = "tuvis"
	nf.Title = *companyName
	nf.Body = *msg
	nf.Sound = "default"

	data := map[string]string{
		"UserPrimeID":  strconv.Itoa(int(UserPrimeID)),
		"Type":         "TransactionPush",
		"ShowFeedback": showFeedback,
		// "show_in_foreground": "true",
	}

	ids := tokenDevices

	c := fcm.NewFcmClient(serverKey)
	c.SetContentAvailable(true)
	c.SetNotificationPayload(&nf)
	// c.SetCollapseKey("Tuvis")
	c.NewFcmRegIdsMsg(*ids, data)

	_, err := c.Send()

	if err != nil {
		return err
	}

	return nil
}

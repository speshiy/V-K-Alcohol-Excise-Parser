package common

import (
	"strconv"
	"strings"

	"github.com/ttacon/libphonenumber"
)

//StripPhone strip phone
func StripPhone(phone *string) {
	phoneReplacer := strings.NewReplacer("(", "", ")", "", "-", "", "+", "", " ", "", "  ", "", "   ", "")
	*phone = phoneReplacer.Replace(*phone)
}

//IsValidPhone check valid phone with country code and return formatted phone
func IsValidPhone(phone *string) bool {
	phoneReplacer := strings.NewReplacer("(", "", ")", "", "-", "", "+", "", " ", "", "  ", "", "   ", "")
	*phone = phoneReplacer.Replace(*phone)
	if len(*phone) == 0 {
		return false
	}

	*phone = "+" + *phone

	defaultCountryCode := "+7"

	num, _ := libphonenumber.Parse(*phone, defaultCountryCode)
	*phone = libphonenumber.Format(num, libphonenumber.INTERNATIONAL)
	return libphonenumber.IsValidNumber(num)
}

//GetContryPhoneCodeByPhone return country phone code
func GetContryPhoneCodeByPhone(phone string) string {
	phoneReplacer := strings.NewReplacer("(", "", ")", "", "-", "", "+", "", " ", "", "  ", "", "   ", "")
	phone = phoneReplacer.Replace(phone)
	if len(phone) == 0 {
		return ""
	}

	phone = "+" + phone

	defaultCountryCode := "+7"

	num, _ := libphonenumber.Parse(phone, defaultCountryCode)
	return "+" + strconv.Itoa(int(*num.CountryCode))
}

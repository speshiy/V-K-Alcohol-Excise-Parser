package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle = i18n.NewBundle(language.English)

// CreateLocalizerBundle reads language files and registers them in i18n bundle
func CreateLocalizerBundle(dir string) (*i18n.Bundle, error) {
	langFiles, _ := ioutil.ReadDir(dir)

	// Enable bundle to understand yaml
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	var translations []byte
	var err error
	for _, file := range langFiles {

		if file.IsDir() || !strings.Contains(file.Name(), ".json") {
			continue
		}

		// Read our language json file
		translations, err = ioutil.ReadFile(dir + "/" + file.Name())
		if err != nil {
			fmt.Printf("Unable to read translation file %s", err.Error())
			return nil, err
		}

		// It parses the bytes in buffer to add translations to the bundle
		bundle.MustParseMessageFileBytes(translations, fmt.Sprintf("i18n/%s", file.Name()))
	}

	return bundle, nil
}

//Translate string to selected locale
func Translate(key string, data map[string]interface{}, locale string) string {
	localizer := i18n.NewLocalizer(bundle, locale)
	msg, err := localizer.Localize(
		&i18n.LocalizeConfig{
			MessageID:    key,
			TemplateData: data,
		},
	)
	if err != nil {
		return err.Error()
	}
	return msg
}

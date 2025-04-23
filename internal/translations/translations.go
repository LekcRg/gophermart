package translations

import (
	"fmt"
	"strings"
)

// possible declension in the future
var trFields = map[string]map[string]string{
	"password": {
		"ru": "пароль",
		"en": "password",
	},
	"login": {
		"ru": "логин",
		"en": "login",
	},
}

func getValidLang(lang string) string {
	lang = strings.ToLower(lang)
	if lang != "en" && lang != "ru" {
		return "en"
	}

	return lang
}

func GetTranslatedField(field, lang string) string {
	lang = getValidLang(lang)
	tr, ok := trFields[field][lang]
	if !ok {
		tr = field
	}
	return tr
}

const (
	ErrAlreadyExists = "already_exists"
)

var trErrs = map[string]map[string]string{
	"already_exists": {
		"ru": "занят",
		"en": "already exists",
	},
}

func GetErr(err, field, lang string) (string, error) {
	lang = getValidLang(lang)
	trField := GetTranslatedField(field, lang)

	trErr, ok := trErrs[err]
	if !ok {
		return "", fmt.Errorf("error %s not found", err)
	}
	trErrMsg, ok := trErr[lang]
	if !ok {
		return "", fmt.Errorf("tr with lang %s for err %s not found",
			lang, err)
	}

	trFullMsgErr := fmt.Sprintf("%s %s", trField, trErrMsg)

	return trFullMsgErr, nil
}

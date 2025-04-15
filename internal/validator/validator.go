package validator

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/LekcRg/gophermart/internal/logger"
	"github.com/LekcRg/gophermart/internal/translations"
	en "github.com/go-playground/locales/en"
	ru "github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslation "github.com/go-playground/validator/v10/translations/en"
	ruTranslation "github.com/go-playground/validator/v10/translations/ru"
	"go.uber.org/zap"
)

type ValidationErrors = validator.ValidationErrors

type Translator map[string]ut.Translator

type Validator struct {
	validate *validator.Validate
	trans    Translator
	ut       *ut.UniversalTranslator
}

func New() *Validator {
	validate := validator.New(validator.WithRequiredStructEnabled())

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	en := en.New()
	ru := ru.New()
	uni := ut.New(en, en, ru)

	enTrans, found := uni.GetTranslator("en")
	if !found {
		logger.Log.Error("not found en translations")
	}

	ruTrans, found := uni.GetTranslator("ru")
	if !found {
		logger.Log.Error("not found ru translations")
	}

	enTranslation.RegisterDefaultTranslations(validate, enTrans)
	ruTranslation.RegisterDefaultTranslations(validate, ruTrans)

	valid := &Validator{
		validate: validate,
		trans: Translator{
			"en": enTrans,
			"ru": ruTrans,
		},
		ut: uni,
	}

	err := validate.RegisterValidation("password", validPassword)
	if err != nil {
		panic(err)
	}
	valid.addTranslates(
		"password",
		"{0} must include uppercase, lowercase letter",
		"{0} должен включать заглавную, строчную букву",
	)

	err = validate.RegisterValidation("luhn-order", validLuhn)
	if err != nil {
		panic(err)
	}
	valid.addTranslates(
		"luhn-order",
		"incorrect order number",
		"неправильный номер заказа",
	)

	return valid
}

func (v *Validator) addTranslates(key, en, ru string) {
	v.regTranslation(v.trans["en"], key, en)
	v.regTranslation(v.trans["ru"], key, ru)
}

func (v *Validator) regTranslation(trans ut.Translator, key, transStr string) {
	v.validate.RegisterTranslation(key, trans, func(ut ut.Translator) error {
		return ut.Add(key, transStr, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(key, fe.Field())
		return t
	})
}

func validPassword(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(value)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(value)
	// hasNumber := regexp.MustCompile(`[0-9]`).MatchString(value)

	return hasUpper && hasLower
	// return hasUpper && hasLower && hasNumber
}

func validLuhn(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	sum := 0
	valueSlice := []rune(value)
	length := len(value)
	parity := length % 2

	for i := range length {
		num, err := strconv.Atoi(string(valueSlice[i]))
		if err != nil {
			return false
		}

		if i%2 == parity {
			num *= 2
		}

		if num > 9 {
			num -= 9
		}

		sum += num
	}

	return sum%10 == 0
}

func (v *Validator) GetValidTranslateErrs(
	err error, lang string,
) map[string]string {
	var validErrs validator.ValidationErrors

	if !errors.As(err, &validErrs) {
		return map[string]string{}
	}
	trans, ok := v.trans[lang]
	if !ok {
		logger.Log.Info("not found translation",
			zap.String("lang", lang))
		lang = "en"
		trans = v.trans[lang]
	}

	tr := validErrs.Translate(trans)
	prettyTr := make(map[string]string, len(validErrs))
	for key, val := range tr {
		// remove type from key
		prettyKey := key[strings.Index(key, ".")+1:]
		prettyVal := strings.Replace(
			val,
			prettyKey,
			translations.GetTranslatedField(prettyKey, lang),
			1,
		)
		prettyTr[prettyKey] = prettyVal
	}

	return prettyTr
}

func (v *Validator) ValidateStruct(s any) error {
	return v.validate.Struct(s)
}

func (v *Validator) ValidateVar(s any, keys string) error {
	return v.validate.Var(s, keys)
}

func (v *Validator) GetTrans(lang string) ut.Translator {
	language, ok := v.trans[lang]
	if !ok {
		return v.trans["en"]
	}

	return language
}

package validator

import (
	"errors"
	"reflect"
	"regexp"
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

type Transalors map[string]ut.Translator

type Validator struct {
	validate *validator.Validate
	trans    Transalors
}

func New() *Validator {
	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.RegisterValidation("password", validPassword)
	if err != nil {
		panic(err)
	}

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
	regTrans(validate, enTrans, "password",
		"{0} must include uppercase, lowercase letter and digit")

	ruTrans, found := uni.GetTranslator("ru")
	if !found {
		logger.Log.Error("not found ru translations")
	}
	regTrans(validate, ruTrans, "password",
		"{0} должен включать заглавную, строчную букву и цифру")

	enTranslation.RegisterDefaultTranslations(validate, enTrans)
	ruTranslation.RegisterDefaultTranslations(validate, ruTrans)

	return &Validator{
		validate: validate,
		trans: Transalors{
			"en": enTrans,
			"ru": ruTrans,
		},
	}
}

func regTrans(valid *validator.Validate, trans ut.Translator, key, transStr string) {
	valid.RegisterTranslation(key, trans, func(ut ut.Translator) error {
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
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(value)

	return hasUpper && hasLower && hasNumber
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
		logger.Log.Warn("not found translation",
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

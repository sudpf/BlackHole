package validation

import (
	"BlackHole/pkg/constant"
	"BlackHole/pkg/env"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"golang.org/x/text/language"
)

type Translator struct {
	uni     *ut.UniversalTranslator
	matcher language.Matcher
}

func NewTranslator() (*Translator, error) {
	translator := &Translator{
		uni:     ut.New(en.New(), zh.New(), en.New()),
		matcher: language.NewMatcher([]language.Tag{language.English, language.SimplifiedChinese}),
	}

	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return translator, nil
	}

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name, _, _ := strings.Cut(fld.Tag.Get("json"), ",")
		if name == "-" {
			return ""
		}
		if name != "" {
			return name
		}

		name, _, _ = strings.Cut(fld.Tag.Get("form"), ",")
		if name == "-" {
			return ""
		}
		return name
	})

	transEn, ok := translator.uni.GetTranslator(constant.LangEnglish)
	if !ok {
		return nil, fmt.Errorf("uni.GetTranslator(%s) failed", constant.LangEnglish)
	}
	if err := enTranslations.RegisterDefaultTranslations(v, transEn); err != nil {
		return nil, err
	}

	transZh, ok := translator.uni.GetTranslator(constant.LangChinese)
	if !ok {
		return nil, fmt.Errorf("uni.GetTranslator(%s) failed", constant.LangChinese)
	}
	if err := zhTranslations.RegisterDefaultTranslations(v, transZh); err != nil {
		return nil, err
	}

	return translator, nil
}

func (t *Translator) TranslateErrors(requestEnv *env.Env, err error) map[string]string {
	if t == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return nil
	}

	trans, ok := t.uni.GetTranslator(t.Language(requestEnv))
	if !ok {
		return nil
	}
	fields := validationErrors.Translate(trans)
	result := map[string]string{}
	for field, err := range fields {
		parts := strings.Split(field, ".")
		result[parts[len(parts)-1]] = err
	}
	return result
}

func (t *Translator) Language(requestEnv *env.Env) string {
	value := constant.LangEnglish
	if requestEnv != nil && requestEnv.Lang != "" {
		value = requestEnv.Lang
	}

	tags, _, err := language.ParseAcceptLanguage(value)
	if err != nil || len(tags) == 0 {
		return constant.LangEnglish
	}

	_, index, _ := t.matcher.Match(tags...)
	if index == 1 {
		return constant.LangChinese
	}
	return constant.LangEnglish
}

package apperror

import (
	"BlackHole/pkg/constant"
	"BlackHole/pkg/env"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Definition struct {
	Code       Code
	HTTPStatus int
	English    string
	Chinese    string
}

type Catalog struct {
	definitions map[Code]Definition
	localizerZh *i18n.Localizer
	localizerEn *i18n.Localizer
	uni         *ut.UniversalTranslator
	matcher     language.Matcher
}

func NewCatalog(definitions ...Definition) (*Catalog, error) {
	catalog := &Catalog{
		definitions: make(map[Code]Definition, len(definitions)),
		uni:         ut.New(en.New(), zh.New(), en.New()),
		matcher:     language.NewMatcher([]language.Tag{language.English, language.SimplifiedChinese}),
	}

	for _, definition := range definitions {
		if definition.Code < 0 {
			return nil, fmt.Errorf("invalid error code %d", definition.Code)
		}
		if definition.HTTPStatus < 100 || definition.HTTPStatus > 599 {
			return nil, fmt.Errorf("invalid HTTP status %d for error code %d", definition.HTTPStatus, definition.Code)
		}
		if definition.Code == Success && (definition.HTTPStatus < 200 || definition.HTTPStatus > 299) {
			return nil, fmt.Errorf("success code %d requires a 2xx HTTP status", Success)
		}
		if _, exists := catalog.definitions[definition.Code]; exists {
			return nil, fmt.Errorf("duplicate error code %d", definition.Code)
		}
		catalog.definitions[definition.Code] = definition
	}

	if _, exists := catalog.definitions[Success]; !exists {
		return nil, fmt.Errorf("success code %d is not defined", Success)
	}
	enMessages := make(map[string]string, len(definitions))
	zhMessages := make(map[string]string, len(definitions))
	for _, definition := range definitions {
		if definition.English == "" || definition.Chinese == "" {
			return nil, fmt.Errorf("messages are not defined for error code %d", definition.Code)
		}
		messageID := MessageID(definition.Code)
		enMessages[messageID] = definition.English
		zhMessages[messageID] = definition.Chinese
	}

	bundle := i18n.NewBundle(language.English)
	for id, translation := range enMessages {
		if err := bundle.AddMessages(language.English, &i18n.Message{ID: id, Other: translation}); err != nil {
			return nil, fmt.Errorf("add english message %q: %w", id, err)
		}
	}
	for id, translation := range zhMessages {
		if err := bundle.AddMessages(language.Chinese, &i18n.Message{ID: id, Other: translation}); err != nil {
			return nil, fmt.Errorf("add chinese message %q: %w", id, err)
		}
	}
	catalog.localizerZh = i18n.NewLocalizer(bundle, constant.LangChinese)
	catalog.localizerEn = i18n.NewLocalizer(bundle, constant.LangEnglish)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name, _, _ := strings.Cut(fld.Tag.Get("json"), ",")
			if name == "-" {
				return ""
			}
			return name
		})

		transEn, ok := catalog.uni.GetTranslator(constant.LangEnglish)
		if !ok {
			return nil, fmt.Errorf("uni.GetTranslator(%s) failed", constant.LangEnglish)
		}
		if err := enTranslations.RegisterDefaultTranslations(v, transEn); err != nil {
			return nil, err
		}

		transZh, ok := catalog.uni.GetTranslator(constant.LangChinese)
		if !ok {
			return nil, fmt.Errorf("uni.GetTranslator(%s) failed", constant.LangChinese)
		}
		if err := zhTranslations.RegisterDefaultTranslations(v, transZh); err != nil {
			return nil, err
		}
	}

	return catalog, nil
}

func (c *Catalog) Lookup(code Code) (Definition, bool) {
	if c == nil {
		return Definition{}, false
	}
	definition, exists := c.definitions[code]
	return definition, exists
}

func (c *Catalog) MessageIDs() []string {
	if c == nil {
		return nil
	}

	codes := make([]int, 0, len(c.definitions))
	for code := range c.definitions {
		codes = append(codes, int(code))
	}
	sort.Ints(codes)

	messageIDs := make([]string, 0, len(codes))
	for _, code := range codes {
		messageIDs = append(messageIDs, MessageID(Code(code)))
	}
	return messageIDs
}

func MessageID(code Code) string {
	return fmt.Sprintf("error_%d", code)
}

func (c *Catalog) Localize(requestEnv *env.Env, messageID string, templateData map[string]any) (string, error) {
	localizer := c.localizerEn
	if c.Language(requestEnv) == constant.LangChinese {
		localizer = c.localizerZh
	}
	return localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})
}

func (c *Catalog) TranslateErrors(requestEnv *env.Env, err error) map[string]string {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return nil
	}

	trans, ok := c.uni.GetTranslator(c.Language(requestEnv))
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

func (c *Catalog) Language(requestEnv *env.Env) string {
	value := constant.LangEnglish
	if requestEnv != nil && requestEnv.Lang != "" {
		value = requestEnv.Lang
	}

	tags, _, err := language.ParseAcceptLanguage(value)
	if err != nil || len(tags) == 0 {
		return constant.LangEnglish
	}

	_, index, _ := c.matcher.Match(tags...)
	if index == 1 {
		return constant.LangChinese
	}
	return constant.LangEnglish
}

package env

import (
	"BlackHole/pkg/constant"
	"BlackHole/pkg/requestctx"
	"context"
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
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Provider struct {
	localizerZh *i18n.Localizer
	localizerEn *i18n.Localizer
	uni         *ut.UniversalTranslator
	matcher     language.Matcher
	messageIDs  map[string]struct{}
}

func NewProvider(enMessages, zhMessages map[string]string) (*Provider, error) {
	if err := validateMessages(enMessages, zhMessages); err != nil {
		return nil, err
	}

	localizerEn, localizerZh, err := newLocalizers(enMessages, zhMessages)
	if err != nil {
		return nil, fmt.Errorf("initialize localizers: %w", err)
	}

	return &Provider{
		localizerZh: localizerZh,
		localizerEn: localizerEn,
		uni:         newUniversalTranslator(),
		matcher:     language.NewMatcher([]language.Tag{language.English, language.SimplifiedChinese}),
		messageIDs:  messageIDs(enMessages),
	}, nil
}

func validateMessages(enMessages, zhMessages map[string]string) error {
	for id := range enMessages {
		if _, exists := zhMessages[id]; !exists {
			return fmt.Errorf("missing chinese message %q", id)
		}
	}
	for id := range zhMessages {
		if _, exists := enMessages[id]; !exists {
			return fmt.Errorf("missing english message %q", id)
		}
	}
	return nil
}

func messageIDs(messages map[string]string) map[string]struct{} {
	result := make(map[string]struct{}, len(messages))
	for id := range messages {
		result[id] = struct{}{}
	}
	return result
}

func newUniversalTranslator() *ut.UniversalTranslator {
	zhT := zh.New()
	enT := en.New()
	return ut.New(enT, zhT, enT)
}

func InitValidatorTranslations(provider *Provider) error {
	if provider == nil {
		return fmt.Errorf("env provider is required")
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name, _, _ := strings.Cut(fld.Tag.Get("json"), ",")
			if name == "-" {
				return ""
			}
			return name
		})

		transEn, ok := provider.uni.GetTranslator(constant.LangEnglish)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", constant.LangEnglish)
		}
		if err := enTranslations.RegisterDefaultTranslations(v, transEn); err != nil {
			return err
		}

		transZh, ok := provider.uni.GetTranslator(constant.LangChinese)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", constant.LangChinese)
		}
		if err := zhTranslations.RegisterDefaultTranslations(v, transZh); err != nil {
			return err
		}
	}

	return nil
}

func newLocalizers(enMessages, zhMessages map[string]string) (*i18n.Localizer, *i18n.Localizer, error) {
	bundle := i18n.NewBundle(language.English)

	for id, translation := range enMessages {
		if err := bundle.AddMessages(language.English, &i18n.Message{
			ID:    id,
			Other: translation,
		}); err != nil {
			return nil, nil, fmt.Errorf("add english message %q: %w", id, err)
		}
	}

	for id, translation := range zhMessages {
		if err := bundle.AddMessages(language.Chinese, &i18n.Message{
			ID:    id,
			Other: translation,
		}); err != nil {
			return nil, nil, fmt.Errorf("add chinese message %q: %w", id, err)
		}
	}

	return i18n.NewLocalizer(bundle, "en"), i18n.NewLocalizer(bundle, "zh"), nil
}

type Env struct {
	Lang      string
	ClientIp  string
	RequestId string
	Trans     ut.Translator

	localizer *i18n.Localizer
}

func (p *Provider) NewEnv(lang string, clientIp string) *Env {
	lang = p.matchLanguage(lang)
	trans, _ := p.uni.GetTranslator(lang)

	localizer := p.localizerEn
	if lang == constant.LangChinese {
		localizer = p.localizerZh
	}

	return &Env{Lang: lang, ClientIp: clientIp, Trans: trans, localizer: localizer}
}

func (p *Provider) NewEnvFromContext(ctx context.Context) *Env {
	lang := requestctx.Language(ctx)
	if lang == "" {
		lang = constant.LangEnglish
	}

	env := p.NewEnv(lang, requestctx.ClientIP(ctx))
	env.RequestId = requestctx.TraceID(ctx)
	return env
}

func (p *Provider) ValidateMessages(messageIDs []string) error {
	for _, id := range messageIDs {
		if _, exists := p.messageIDs[id]; !exists {
			return fmt.Errorf("message %q is not defined", id)
		}
	}
	return nil
}

func (p *Provider) matchLanguage(value string) string {
	tags, _, err := language.ParseAcceptLanguage(value)
	if err != nil || len(tags) == 0 {
		return constant.LangEnglish
	}

	_, index, _ := p.matcher.Match(tags...)
	if index == 1 {
		return constant.LangChinese
	}
	return constant.LangEnglish
}

func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		parts := strings.Split(field, ".")
		res[parts[len(parts)-1]] = err
	}
	return res
}

func (ev *Env) TranslateErrors(err error) map[string]string {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return nil
	}
	return removeTopStruct(validationErrors.Translate(ev.Trans))
}

func (ev *Env) Localize(messageID string, templateData map[string]any) (string, error) {
	return ev.localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})
}

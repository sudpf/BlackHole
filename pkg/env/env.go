package env

import (
	"BlackHole/pkg/constant"
	"BlackHole/pkg/requestctx"
	"context"
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
}

func NewProvider(enMessages, zhMessages map[string]string) (*Provider, error) {
	uni, err := setupTranslations()
	if err != nil {
		return nil, fmt.Errorf("setup translations: %w", err)
	}

	localizerEn, localizerZh, err := newLocalizers(enMessages, zhMessages)
	if err != nil {
		return nil, fmt.Errorf("initialize localizers: %w", err)
	}

	return &Provider{
		localizerZh: localizerZh,
		localizerEn: localizerEn,
		uni:         uni,
	}, nil
}

func setupTranslations() (*ut.UniversalTranslator, error) {
	zhT := zh.New()
	enT := en.New()
	uni := ut.New(enT, zhT, enT)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name, _, _ := strings.Cut(fld.Tag.Get("json"), ",")
			if name == "-" {
				return ""
			}
			return name
		})

		transEn, ok := uni.GetTranslator(constant.LangEnglish)
		if !ok {
			return nil, fmt.Errorf("uni.GetTranslator(%s) failed", constant.LangEnglish)
		}
		if err := enTranslations.RegisterDefaultTranslations(v, transEn); err != nil {
			return nil, err
		}

		transZh, ok := uni.GetTranslator(constant.LangChinese)
		if !ok {
			return nil, fmt.Errorf("uni.GetTranslator(%s) failed", constant.LangChinese)
		}
		if err := zhTranslations.RegisterDefaultTranslations(v, transZh); err != nil {
			return nil, err
		}
	}

	return uni, nil
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

func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		parts := strings.Split(field, ".")
		res[parts[len(parts)-1]] = err
	}
	return res
}

func (ev *Env) TranslatErrors(err error) map[string]string {
	errs, _ := err.(validator.ValidationErrors)

	return removeTopStruct(errs.Translate(ev.Trans))
}

func (ev *Env) MustLocalize(message string) string {
	return ev.localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: message,
	})
}

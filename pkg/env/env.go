package env

import (
	"BlackHole/pkg/constant"
	"BlackHole/pkg/locales"
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

var (
	localizerZh *i18n.Localizer
	localizerEn *i18n.Localizer
	uni         *ut.UniversalTranslator
)

func SetupTranslations() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New()
		enT := en.New()
		uni = ut.New(enT, zhT, enT)

		transEn, ok := uni.GetTranslator(constant.LangEnglish)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", constant.LangEnglish)
		}
		if err := enTranslations.RegisterDefaultTranslations(v, transEn); err != nil {
			return err
		}

		transZh, ok := uni.GetTranslator(constant.LangChinese)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", constant.LangChinese)
		}
		if err := zhTranslations.RegisterDefaultTranslations(v, transZh); err != nil {
			return err
		}
	}

	// 注册自定义验证器和翻译器

	return nil
}

func InitLocalizer() {
	// 创建一个新的 i18n bundle
	bundle := i18n.NewBundle(language.English)

	// 加载翻译到 bundle 中
	for id, translation := range locales.EnTranslations {
		bundle.AddMessages(language.English, &i18n.Message{
			ID:    id,
			Other: translation,
		})
	}

	for id, translation := range locales.ZhTranslations {
		bundle.AddMessages(language.Chinese, &i18n.Message{
			ID:    id,
			Other: translation,
		})
	}

	localizerZh = i18n.NewLocalizer(bundle, "zh")
	localizerEn = i18n.NewLocalizer(bundle, "en")
}

type Env struct {
	Lang      string
	ClientIp  string
	RequestId string
	Trans     ut.Translator
}

func NewEnv(lang string, clientIp string) *Env {
	trans, _ := uni.GetTranslator(lang)

	return &Env{Lang: lang, ClientIp: clientIp, Trans: trans}
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
	if ev.Lang == constant.LangChinese {
		return localizerZh.MustLocalize(&i18n.LocalizeConfig{
			MessageID: message,
		})
	}

	return localizerEn.MustLocalize(&i18n.LocalizeConfig{
		MessageID: message,
	})
}

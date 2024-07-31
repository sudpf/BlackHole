package env

import (
	"BlackHole/pkg/constant"
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
)

var (
	uni *ut.UniversalTranslator
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

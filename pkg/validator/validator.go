// Package validator 扩展 gin 的参数校验：注册自定义规则并输出中文错误信息。
package validator

import (
	"errors"
	"reflect"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	validatorlib "github.com/go-playground/validator/v10"
	zhtrans "github.com/go-playground/validator/v10/translations/zh"

	"github.com/seq-tech/seq-common-service/pkg/hash"
)

var (
	trans     ut.Translator
	phoneRule = regexp.MustCompile(`^1[3-9]\d{9}$`)
	// 用户名：字母开头，允许字母数字下划线，长度 4-32。
	usernameRule     = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{3,31}$`)
	alphaNumDashRule = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

// Init 注册自定义校验规则与中文翻译，需在路由注册前调用一次。
func Init() error {
	v, ok := binding.Validator.Engine().(*validatorlib.Validate)
	if !ok {
		return errors.New("gin validator engine is not *validator.Validate")
	}

	// 使用 json tag 作为错误信息中的字段名。
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" || name == "" {
			return fld.Name
		}
		return name
	})

	locale := zh.New()
	uni := ut.New(locale, locale)
	t, found := uni.GetTranslator("zh")
	if !found {
		return errors.New("zh translator not found")
	}
	if err := zhtrans.RegisterDefaultTranslations(v, t); err != nil {
		return err
	}
	trans = t

	rules := []struct {
		tag     string
		fn      validatorlib.Func
		message string
	}{
		{"phone", validatePhone, "{0}必须是有效的手机号"},
		{"password", validatePassword, "{0}至少 8 位且需同时包含字母和数字"},
		{"username", validateUsername, "{0}需以字母开头，仅含字母、数字、下划线，长度 4-32"},
		{"alphanumdash", validateAlphaNumDash, "{0}只能包含字母、数字、下划线和中划线"},
	}
	for _, r := range rules {
		if err := v.RegisterValidation(r.tag, r.fn); err != nil {
			return err
		}
		if err := registerTranslation(v, t, r.tag, r.message); err != nil {
			return err
		}
	}
	return nil
}

func registerTranslation(v *validatorlib.Validate, t ut.Translator, tag, message string) error {
	return v.RegisterTranslation(tag, t,
		func(ut ut.Translator) error { return ut.Add(tag, message, true) },
		func(ut ut.Translator, fe validatorlib.FieldError) string {
			msg, _ := ut.T(tag, fe.Field())
			return msg
		},
	)
}

func validatePhone(fl validatorlib.FieldLevel) bool {
	return phoneRule.MatchString(fl.Field().String())
}

func validatePassword(fl validatorlib.FieldLevel) bool {
	return hash.CheckStrength(fl.Field().String()) == nil
}

func validateUsername(fl validatorlib.FieldLevel) bool {
	return usernameRule.MatchString(fl.Field().String())
}

func validateAlphaNumDash(fl validatorlib.FieldLevel) bool {
	return alphaNumDashRule.MatchString(fl.Field().String())
}

// Translate 把校验错误转换为可直接返回给客户端的中文提示。
func Translate(err error) string {
	if err == nil {
		return ""
	}
	var ve validatorlib.ValidationErrors
	if !errors.As(err, &ve) || trans == nil {
		return "请求参数错误"
	}
	msgs := make([]string, 0, len(ve))
	for _, fe := range ve {
		msgs = append(msgs, fe.Translate(trans))
	}
	return strings.Join(msgs, "；")
}

package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/i18n/gi18n"
)

var (
	translate *gi18n.Manager
)

func GinI18nLocalize() gin.HandlerFunc {
	return func(c *gin.Context) {
		locale := c.GetHeader("lang")
		if locale != "" {
			c.Request.Header.Set("Accept-Language", locale)
		}
		headreLang := c.GetHeader("Accept-Language")
		lang := "en"
		if strings.HasPrefix(headreLang, "zh") || strings.HasPrefix(headreLang, "ZH") {
			lang = "zh"
		}
		var (
			t = GetTranslate()
		)
		t.SetLanguage(lang)
		c.Next()
	}
}

func GetEXPERTPath() string {
	basePath := os.Getenv("EXPERT_SERVICE")
	if basePath == "" {
		basePath, _ = os.Getwd()
		//panic("please set EXPERT_SERVICE env")
	}
	return basePath
}

func GetTranslate() *gi18n.Manager {
	if translate != nil {
		return translate
	}
	translate = gi18n.New(gi18n.Options{Path: GetEXPERTPath() + "/i18n"})
	return translate
}

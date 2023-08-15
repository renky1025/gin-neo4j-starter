package middleware

import (
	"github.com/gin-gonic/gin"
)

var (
	NotLoginUri = []string{"/api/v1.0/products"}
	//NotAuthUri  = []string{"/rsagen", "/localverify", "/licenselog/sync"}
)

// TokenAuth Token认证中间件
func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 免登录接口
		// tokenCheck := c.Request.Header.Get("SkipVerifyToken")
		// if utils.Contains(NotLoginUri, c.Request.RequestURI) || tokenCheck == "true" {
		// 	log.Println("skip token chek, go on ==>")
		// 	c.Next()
		// 	return
		// }
		// // Token是否为空
		// auth := c.Request.Header.Get("Authorization")
		// token := strings.TrimPrefix(auth, "Bearer ")
		// if token == auth || token == "" {
		// 	response.Fail(c, response.TokenEmpty)
		// 	c.Abort()
		// 	return
		// }
		c.Next()
		return
	}
}

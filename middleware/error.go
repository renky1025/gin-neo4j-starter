package middleware

import (
	"runtime/debug"

	"go-gin-restful-service/log"
	"go-gin-restful-service/response"

	"github.com/gin-gonic/gin"
)

// ErrorRecover 异常恢复中间件
func ErrorRecover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				switch v := r.(type) {
				// 自定义类型
				case response.RespType:
					log.Logger.Warnf(
						"Request Fail by recover: url=[%s], resp=[%+v]", c.Request.URL.Path, v)
					var data interface{}
					if v.Data() == nil {
						data = []string{}
					}
					response.Result(c, v, data)
				// 其他类型
				default:
					log.Logger.Errorf("stacktrace from panic: %+v\n%s", r, string(debug.Stack()))
					response.Fail(c, response.SystemError)
				}
				c.Abort()
				return
			}
		}()
		c.Next()
	}
}

package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func TenantFilterM() gin.HandlerFunc {
	return func(c *gin.Context) {
		headers := c.Request.Header
		for key, value := range headers {
			//log.Printf("Headers from request ==> headerName=%s, headerValue = %s", key, value)
			if strings.ToLower(key) == "x-tenant" && len(value) > 0 {
				c.Request.Header.Set("tenantId", value[0])
			}
		}
		// tenantId := c.GetHeader("x-tenant")
		// client := c.GetHeader("x-client")
		// version := c.GetHeader("version")
		// log.Printf("Headers from request ==> tenantId=%s, x-client = %s, version = %s \n", tenantId, client, version)
		// c.Request.Header.Set("tenantId", tenantId)
		c.Next()
	}
}

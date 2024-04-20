package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/cookiejar"
	"wiki-user/server/model"
)

func SessionManager() gin.HandlerFunc {
	return func(c *gin.Context) {
		userName := c.Request.Header.Get("username")
		session, exists := c.Get(userName)
		if !exists {
			jar, _ := cookiejar.New(nil)
			client := &http.Client{Jar: jar}
			session = model.UserSession{Client: client}
			c.Set(userName, session)
		}
		c.Next()
	}
}

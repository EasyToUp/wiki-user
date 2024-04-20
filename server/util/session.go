package utils

import (
	"github.com/gin-gonic/gin"
	"wiki-user/server/model"
)

func GetSession(c *gin.Context) *model.UserSession {
	userName := c.Request.Header.Get("username")
	session := c.MustGet(userName).(model.UserSession)
	return &session
}

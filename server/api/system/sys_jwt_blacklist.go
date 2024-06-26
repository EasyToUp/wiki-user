package system

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"wiki-user/server/global"
	"wiki-user/server/model/response"
	system2 "wiki-user/server/model/system"
)

type JwtApi struct{}

// JsonInBlacklist
// @Tags      Jwt
// @Summary   jwt加入黑名单
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{msg=string}  "jwt加入黑名单"
// @Router    /jwt/jsonInBlacklist [post]
func (j *JwtApi) JsonInBlacklist(c *gin.Context) {
	token := c.Request.Header.Get("x-token")
	jwt := system2.JwtBlacklist{Jwt: token}
	err := jwtService.JsonInBlacklist(jwt)
	if err != nil {
		global.WK_LOG.Error("jwt作废失败!", zap.Error(err))
		response.FailWithMessage("jwt作废失败", c)
		return
	}
	response.OkWithMessage("jwt作废成功", c)
}

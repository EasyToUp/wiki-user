package system

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"wiki-user/server/global"
	"wiki-user/server/model/system"
	systemRes "wiki-user/server/model/system/response"

	"wiki-user/server/model"
	"wiki-user/server/model/response"
	utils "wiki-user/server/util"
)

type AccountApi struct{}

func (s *AccountApi) Login(c *gin.Context) {
	{

		resp, err := accountService.LoginUser(c)
		if err != nil {
			response.Result(500, "Login failed", err.Error(), c)
			return
		}
		response.Result(200, resp, "Login success", c)
	}

	//var l request.Login
	//key := c.ClientIP()
	//
	//// 判断验证码是否开启
	//openCaptcha := global.WK_CONFIG.Captcha.OpenCaptcha               // 是否开启防爆次数
	//openCaptchaTimeOut := global.WK_CONFIG.Captcha.OpenCaptchaTimeOut // 缓存超时时间
	//v, ok := global.BlackCache.Get(key)
	//if !ok {
	//	global.BlackCache.Set(key, 1, time.Second*time.Duration(openCaptchaTimeOut))
	//}
	//
	//var oc bool = openCaptcha == 0 || openCaptcha < interfaceToInt(v)
	//
	//if !oc || (l.CaptchaId != "" && l.Captcha != "" && store.Verify(l.CaptchaId, l.Captcha, true)) {
	//	//u := &system.SysUser{Username: l.Username, Password: l.Password}
	//
	//	// todo 登陆逻辑 补充
	//
	////=========================================//
	//user, err := accountService.Login(u)
	//if err != nil {
	//	global.WK_LOG.Error("登陆失败! 用户名不存在或者密码错误!", zap.Error(err))
	//	// 验证码次数+1
	//	global.BlackCache.Increment(key, 1)
	//	response.FailWithMessage("用户名不存在或者密码错误", c)
	//	return
	//}
	//if user.Enable != 1 {
	//	global.WK_LOG.Error("登陆失败! 用户被禁止登录!")
	//	// 验证码次数+1
	//	global.BlackCache.Increment(key, 1)
	//	response.FailWithMessage("用户被禁止登录", c)
	//	return
	//}
	//s.TokenNext(c, *user)
	//return
	//}
	//// 验证码次数+1
	//global.BlackCache.Increment(key, 1)
	//response.FailWithMessage("验证码错误", c)
}

func (s *AccountApi) Register(c *gin.Context) {
	resp, err := accountService.RegisterUser(c)
	if err != nil {
		response.Result(500, "Registration failed", err.Error(), c)
		return
	}
	response.Result(200, resp, "Registration success", c)
}

// TokenNext 登录以后签发jwt
func (s *AccountApi) TokenNext(c *gin.Context, user system.SysUser) {
	j := &utils.JWT{SigningKey: []byte(global.WK_CONFIG.JWT.SigningKey)} // 唯一签名
	claims := j.CreateClaims(model.BaseClaims{
		UUID:        user.UUID,
		ID:          user.ID,
		NickName:    user.NickName,
		Username:    user.Username,
		AuthorityId: user.AuthorityId,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		global.WK_LOG.Error("获取token失败!", zap.Error(err))
		response.FailWithMessage("获取token失败", c)
		return
	}
	if !global.WK_CONFIG.System.UseMultipoint {
		response.OkWithDetailed(systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, "登录成功", c)
		return
	}

	if jwtStr, err := jwtService.GetRedisJWT(user.Username); err == redis.Nil {
		if err := jwtService.SetRedisJWT(token, user.Username); err != nil {
			global.WK_LOG.Error("设置登录状态失败!", zap.Error(err))
			response.FailWithMessage("设置登录状态失败", c)
			return
		}
		response.OkWithDetailed(systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, "登录成功", c)
	} else if err != nil {
		global.WK_LOG.Error("设置登录状态失败!", zap.Error(err))
		response.FailWithMessage("设置登录状态失败", c)
	} else {
		var blackJWT system.JwtBlacklist
		blackJWT.Jwt = jwtStr
		if err := jwtService.JsonInBlacklist(blackJWT); err != nil {
			response.FailWithMessage("jwt作废失败", c)
			return
		}
		if err := jwtService.SetRedisJWT(token, user.Username); err != nil {
			response.FailWithMessage("设置登录状态失败", c)
			return
		}
		response.OkWithDetailed(systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, "登录成功", c)
	}
}

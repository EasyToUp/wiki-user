package system

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
	"wiki-user/server/global"
	"wiki-user/server/model/system"
	systemRes "wiki-user/server/model/system/response"

	"wiki-user/server/model"
	"wiki-user/server/model/response"
	"wiki-user/server/model/system/request"
	utils "wiki-user/server/util"
)

type AccountApi struct{}

func (s *AccountApi) Login(c *gin.Context) {

	var l request.Login
	key := c.ClientIP()

	// 判断验证码是否开启
	openCaptcha := global.WK_CONFIG.Captcha.OpenCaptcha               // 是否开启防爆次数
	openCaptchaTimeOut := global.WK_CONFIG.Captcha.OpenCaptchaTimeOut // 缓存超时时间
	v, ok := global.BlackCache.Get(key)
	if !ok {
		global.BlackCache.Set(key, 1, time.Second*time.Duration(openCaptchaTimeOut))
	}

	var oc bool = openCaptcha == 0 || openCaptcha < interfaceToInt(v)

	if !oc || (l.CaptchaId != "" && l.Captcha != "" && store.Verify(l.CaptchaId, l.Captcha, true)) {
		u := &system.SysUser{Username: l.Username, Password: l.Password}

		// todo 登陆逻辑 补充
		{
			var loginInfo model.LoginInfo
			if err := c.ShouldBind(&loginInfo); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			token, err := getToken(c, global.WK_CONFIG.System.ApiUrl, "login")
			if err != nil {
				c.JSON(500, gin.H{"error": "get login token failed"})
				return
			}
			resp, err := accountService.LoginUser(c, &loginInfo, global.WK_CONFIG.System.ApiUrl, token)
			if err != nil {
				c.JSON(500, gin.H{"error": "Login failed"})
				return
			}
			c.JSON(200, gin.H{"resp": resp})
		}
		//=========================================//
		user, err := accountService.Login(u)
		if err != nil {
			global.WK_LOG.Error("登陆失败! 用户名不存在或者密码错误!", zap.Error(err))
			// 验证码次数+1
			global.BlackCache.Increment(key, 1)
			response.FailWithMessage("用户名不存在或者密码错误", c)
			return
		}
		if user.Enable != 1 {
			global.WK_LOG.Error("登陆失败! 用户被禁止登录!")
			// 验证码次数+1
			global.BlackCache.Increment(key, 1)
			response.FailWithMessage("用户被禁止登录", c)
			return
		}
		s.TokenNext(c, *user)
		return
	}
	// 验证码次数+1
	global.BlackCache.Increment(key, 1)
	response.FailWithMessage("验证码错误", c)
}

func (s *AccountApi) Register(c *gin.Context) {
	var registerInfo model.RegisterInfo
	if err := c.ShouldBind(&registerInfo); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	token, err := getToken(c, global.WK_CONFIG.System.ApiUrl, "createaccount")
	if err != nil {
		c.JSON(500, gin.H{"error": "get tken failed"})
		return
	}
	resp, err := accountService.RegisterUser(c, &registerInfo, global.WK_CONFIG.System.ApiUrl, token)
	if err != nil {
		c.JSON(500, gin.H{"error": "Registration failed"})
		return
	}
	c.JSON(200, gin.H{"resp": resp})
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

// 改进的getToken函数，加入了错误处理和明确的日志记录
func getToken(c *gin.Context, apiUrl string, tokenType string) (string, error) {

	client := utils.GetSession(c).Client

	data := url.Values{}
	data.Set("action", "query")
	data.Set("meta", "tokens")
	data.Set("type", tokenType)
	data.Set("format", "json")

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("creating token request failed: %v", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("performing token request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading token response body failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("unmarshaling token response failed: %v", err)
	}

	tokenPath := fmt.Sprintf("%stoken", tokenType)
	token, ok := result["query"].(map[string]interface{})["tokens"].(map[string]interface{})[tokenPath].(string)
	if !ok {
		return "", fmt.Errorf("token not found in the response")
	}
	fmt.Println("======>>>>tokenType", tokenType)
	fmt.Println("======>>>>token", token)

	return token, nil
}

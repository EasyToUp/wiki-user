package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"wiki-user/server/global"
	"wiki-user/server/model"
	"wiki-user/server/model/system"
	utils "wiki-user/server/util"
)

type AccountService struct{}

// LoginUser logs in a user with the given login information
func (accountService *AccountService) LoginUser(c *gin.Context) (interface{}, error) {
	//client := utils.GetSession(c).Client
	client, _ := utils.CreateHTTPClient()

	var loginInfo model.LoginInfo
	if err := c.ShouldBind(&loginInfo); err != nil {
		return "参数错误", err
	}
	token, err := getToken(client, global.WK_CONFIG.System.ApiUrl, "login")
	if err != nil {
		return "get clientlogin token failed", err
	}

	data := url.Values{}
	data.Set("action", "clientlogin")
	data.Set("username", loginInfo.Username)
	data.Set("password", loginInfo.Password)
	data.Set("format", "json")
	data.Set("logintoken", token)
	data.Set("loginreturnurl", global.WK_CONFIG.System.ApiUrl)

	resp, err := client.PostForm(global.WK_CONFIG.System.ApiUrl, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// 获取并设置cookies
	url, _ := url.Parse(global.WK_CONFIG.System.ApiUrl)
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	status := result["clientlogin"].(map[string]interface{})["status"].(string)
	if status == "PASS" {
		global.WK_LOG.Info("Login successful! ", zap.Any("response", result["clientlogin"]))
	} else {
		global.WK_LOG.Info("Login failed:", zap.Any("response", result["clientlogin"]))
		jsonData, _ := json.Marshal(result["clientlogin"])
		err = errors.New(string(jsonData))
	}

	cookies := client.Jar.Cookies(url)
	for _, cookie := range cookies {
		c.SetCookie(cookie.Name, cookie.Value, int(cookie.MaxAge), cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
	}
	// Extract token from body or handle errors based on actual API response
	return result["clientlogin"], err // Placeholder for actual token extraction logic
}

// RegisterUser registers a user with the given registration information
func (accountService *AccountService) RegisterUser(c *gin.Context) (interface{}, error) {

	client, _ := utils.CreateHTTPClient()

	var registerInfo model.RegisterInfo
	if err := c.ShouldBind(&registerInfo); err != nil {
		return "参数错误", err
	}
	token, err := getToken(client, global.WK_CONFIG.System.ApiUrl, "createaccount")
	if err != nil {
		return "get createaccount token failed", err
	}

	data := url.Values{}
	data.Set("action", "createaccount")
	data.Set("username", registerInfo.Username)
	data.Set("password", registerInfo.Password)
	data.Set("retype", registerInfo.Password)
	data.Set("email", registerInfo.Email)
	data.Set("format", "json")
	data.Set("createtoken", token)
	data.Set("createreturnurl", global.WK_CONFIG.System.ApiUrl)

	resp, err := client.PostForm(global.WK_CONFIG.System.ApiUrl, data)
	if err != nil {
		return "发送登陆请求失败", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "读取wiki返回失败", err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if result["createaccount"] != nil && result["createaccount"].(map[string]interface{})["status"].(string) == "PASS" {
		global.WK_LOG.Info("Account creation successful!", zap.Any("response", result["createaccount"]))

	} else {
		global.WK_LOG.Error("Account creation failed:", zap.Any("response", result["createaccount"]))
		jsonData, _ := json.Marshal(result["createaccount"])
		err = errors.New(string(jsonData))
	}
	// 构建请求 URL
	loginURL, _ := url.Parse(global.WK_CONFIG.System.ApiUrl)

	// 日志记录 cookies
	logCookies(loginURL, client.Jar.(*cookiejar.Jar))
	// Handle API response to determine if registration was successful
	return result["createaccount"], err // Placeholder for error handling based on API response
}

func (accountService *AccountService) Login(u *system.SysUser) (userInter *system.SysUser, err error) {
	if nil == global.WK_DB {
		return nil, fmt.Errorf("db not init")
	}

	var user system.SysUser
	err = global.WK_DB.Where("username = ?", u.Username).Preload("Authorities").Preload("Authority").First(&user).Error
	if err == nil {
		if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
			return nil, errors.New("密码错误")
		}
	}
	return &user, err
}

func logCookies(u *url.URL, jar *cookiejar.Jar) {
	cookies := jar.Cookies(u)
	for _, cookie := range cookies {
		fmt.Println("Cookie: ", cookie.Name, cookie.Value)
	}
}

// 改进的getToken函数，加入了错误处理和明确的日志记录
func getToken(client *http.Client, apiUrl string, tokenType string) (string, error) {

	//client := utils.GetSession(c).Client

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

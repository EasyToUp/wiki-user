package user

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http/cookiejar"
	"net/url"
	"wiki-user/server/global"
	"wiki-user/server/model"
	"wiki-user/server/model/system"
	utils "wiki-user/server/util"
)

type AccountService struct{}

// LoginUser logs in a user with the given login information
func (accountService *AccountService) LoginUser(c *gin.Context, loginInfo *model.LoginInfo, apiURL string, token string) (string, error) {
	client := utils.GetSession(c).Client

	data := url.Values{}
	data.Set("action", "login")
	data.Set("lgname", loginInfo.Username)
	data.Set("lgpassword", loginInfo.Password)
	data.Set("format", "json")
	data.Set("lgtoken", token)

	resp, err := client.PostForm(apiURL, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Extract token from body or handle errors based on actual API response
	return string(body), nil // Placeholder for actual token extraction logic
}

// RegisterUser registers a user with the given registration information
func (accountService *AccountService) RegisterUser(c *gin.Context, registerInfo *model.RegisterInfo, apiURL string, token string) (string, error) {

	client := utils.GetSession(c).Client

	data := url.Values{}
	data.Set("action", "createaccount")
	data.Set("username", registerInfo.Username)
	data.Set("password", registerInfo.Password)
	data.Set("retype", registerInfo.Password)
	data.Set("email", registerInfo.Email)
	data.Set("format", "json")
	data.Set("createtoken", token)
	data.Set("createreturnurl", "http://localhost:9090/index.php")

	resp, err := client.PostForm(apiURL, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// 构建请求 URL
	loginURL, _ := url.Parse("http://localhost:9090/api.php")

	// 日志记录 cookies
	logCookies(loginURL, client.Jar.(*cookiejar.Jar))
	// Handle API response to determine if registration was successful
	return string(body), nil // Placeholder for error handling based on API response
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

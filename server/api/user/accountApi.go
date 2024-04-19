package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"wiki-user/server/global"
	"wiki-user/server/model"
)

type AccountApi struct{}

func (s *AccountApi) Login(c *gin.Context) {
	var loginInfo model.LoginInfo
	if err := c.ShouldBind(&loginInfo); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	token, err := getToken(global.CONFIG.System.ApiUrl, "login")
	if err != nil {
		c.JSON(500, gin.H{"error": "get login token failed"})
	}
	resp, err := accountService.LoginUser(&loginInfo, global.CONFIG.System.ApiUrl, token)
	if err != nil {
		c.JSON(500, gin.H{"error": "Login failed"})
		return
	}
	c.JSON(200, gin.H{"resp": resp})
}

func (s *AccountApi) Register(c *gin.Context) {
	var registerInfo model.RegisterInfo
	if err := c.ShouldBind(&registerInfo); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	token, err := getToken(global.CONFIG.System.ApiUrl, "createaccount")
	if err != nil {
		c.JSON(500, gin.H{"error": "get tken failed"})
	}
	resp, err := accountService.RegisterUser(&registerInfo, global.CONFIG.System.ApiUrl, token)
	if err != nil {
		c.JSON(500, gin.H{"error": "Registration failed"})
		return
	}
	c.JSON(200, gin.H{"resp": resp})
}

// 改进的getToken函数，加入了错误处理和明确的日志记录
func getToken(apiUrl string, tokenType string) (string, error) {
	client := &http.Client{}
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

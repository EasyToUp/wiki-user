package user

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"wiki-user/server/model"
)

type AccountService struct{}

// LoginUser logs in a user with the given login information
func (accountService *AccountService) LoginUser(loginInfo *model.LoginInfo, apiURL string, token string) (string, error) {
	data := url.Values{}
	data.Set("action", "login")
	data.Set("lgname", loginInfo.Username)
	data.Set("lgpassword", loginInfo.Password)
	data.Set("format", "json")
	data.Set("lgtoken", token)

	resp, err := http.PostForm(apiURL, data)
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
func (accountService *AccountService) RegisterUser(registerInfo *model.RegisterInfo, apiURL string, token string) (string, error) {

	data := url.Values{}
	data.Set("action", "createaccount")
	data.Set("username", registerInfo.Username)
	data.Set("password", registerInfo.Password)
	data.Set("email", registerInfo.Email)
	data.Set("format", "json")
	data.Set("createtoken", token)
	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Handle API response to determine if registration was successful
	return string(body), nil // Placeholder for error handling based on API response
}

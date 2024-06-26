package utils

import (
	"net/http"
	"net/http/cookiejar"
)

func CreateHTTPClient() (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	return &http.Client{Jar: jar}, nil
}

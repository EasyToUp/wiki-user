package model

import "net/http"

type UserSession struct {
	Client   *http.Client
	Username string `form:"username" json:"username" binding:"required"`
}

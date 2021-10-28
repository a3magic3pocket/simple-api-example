package controllers

import "encoding/json"

type SuccessResponse struct {
	Data json.RawMessage `json:"data"`
}

type FailResponse struct {
	Error json.RawMessage `json:"error"`
}

type AuthTokenResponse struct {
	Code    int    `json:"code"`
	Expire  string `json:"expire"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

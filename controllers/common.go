package controllers

import (
	"encoding/json"
)

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

type SwagFail struct {
	Error string `json:"error"`
}

type SwagSucc struct {
	Data string `json:"data"`
}

type SwagSuccRetrieveLockers struct {
	Data []LockerOutput `json:"data"`
}

type SwagSuccAuth struct {
	Data AuthTokenResponse `json:"data"`
}

type SwagSuccRetrieveUser struct {
	Data UserOutput `json:"data"`
}

package test

import (
	"encoding/json"
	"log"
	"simple-api-example/controllers"
	"simple-api-example/utils"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestLogout(t *testing.T) {
	TI.Reset()
	url := TI.Server.URL + "/logout"
	method := "POST"

	// 정상 요청
	data := map[string]string{}
	headers := map[string]string{
		"Authorization": "Bearer " + TI.Bearers["normal1"],
		"Content-Type":  "Application/json",
	}
	statusCode, body, err := utils.TestRequest(
		url, method, data, headers, nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, statusCode, 200)

	// 요청 결과 확인
	succResp := controllers.SuccessResponse{}
	if err := json.Unmarshal(body, &succResp); err != nil {
		log.Fatal(err)
	}

	authTokenResp := controllers.AuthTokenResponse{}
	err = json.Unmarshal(succResp.Data, &authTokenResp)
	if err != nil {
		log.Fatalln(err)
	}

	assert.Equal(t, authTokenResp.Code, 200)
	assert.Equal(t, authTokenResp.Message, "success")
}

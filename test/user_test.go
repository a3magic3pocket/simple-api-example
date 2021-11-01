package test

import (
	"encoding/json"
	"log"
	"simple-api-example/controllers"
	"simple-api-example/models"
	"simple-api-example/utils"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestCreateUser(t *testing.T) {
	TI.Reset()
	url := TI.Server.URL + "/user"
	method := "POST"

	// 정상 요청
	data := map[string]string{
		"UserName": "new_normal",
		"Password": "new_password1",
		"Group":    "normal",
	}
	headers := map[string]string{
		"Content-Type": "Application/json",
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

	succMsg := ""
	if err := json.Unmarshal(succResp.Data, &succMsg); err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, succMsg, "success")

	// 새로 등록한 유저가 있는지 확인
	tmpUser := models.User{}
	err = tmpUser.Get("new_normal")
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, tmpUser.Name, "new_normal")

}

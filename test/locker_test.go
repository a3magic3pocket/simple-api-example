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

var lockersMock = map[string]models.Lockers{
	"normal1": []models.Locker{
		{Location: "A"},
		{Location: "B"},
		{Location: "C"},
		{Location: "D"},
		{Location: "E"},
	},
}

// setLockersMock : locker mock 데이터 삽입
func setLockersMock() {
	url := TI.Server.URL + "/lockers"
	method := "POST"

	// 정상 요청
	data := lockersMock["normal1"]
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

	if statusCode != 200 {
		log.Fatalf("failed to setLockersMock. body: %s", body)
	}
}

func TestCreateLockers(t *testing.T) {
	TI.Reset()
	url := TI.Server.URL + "/lockers"
	method := "POST"

	// 정상 요청
	data := []map[string]string{
		{"Location": "A"},
		{"Location": "B"},
		{"Location": "C"},
		{"Location": "D"},
		{"Location": "E"},
	}
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

	succMsg := ""
	if err := json.Unmarshal(succResp.Data, &succMsg); err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, succMsg, "success")
}

func TestRetreiveLockers(t *testing.T) {
	TI.Reset()
	setLockersMock()
	url := TI.Server.URL + "/lockers"
	method := "GET"

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

	lockers := models.Lockers{}
	if err := json.Unmarshal(succResp.Data, &lockers); err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, len(lockersMock["normal1"]), len(lockers))

	exists := false
	for _, locker := range lockers {
		if locker.Location == lockersMock["normal1"][0].Location {
			exists = true
			break
		}
	}
	assert.Equal(t, exists, true)
}

func TestRetreiveLocker(t *testing.T) {
	TI.Reset()
	setLockersMock()
	// 정상 요청
	url := TI.Server.URL + "/locker/1"
	method := "GET"

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

	locker := models.Locker{}
	if err := json.Unmarshal(succResp.Data, &locker); err != nil {
		log.Fatal(err)
	}

	assert.NotEqual(t, locker.ID, 0)
	assert.Equal(t, locker.ID, 1)

	// 없는 ID로 요청
	url = TI.Server.URL + "/locker/99"
	method = "GET"

	data = map[string]string{}
	headers = map[string]string{
		"Authorization": "Bearer " + TI.Bearers["normal1"],
		"Content-Type":  "Application/json",
	}

	statusCode, body, err = utils.TestRequest(
		url, method, data, headers, nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, statusCode, 404)

	// 요청 결과 확인
	failResp := controllers.FailResponse{}
	if err := json.Unmarshal(body, &failResp); err != nil {
		log.Fatal(err)
	}

	errMsg := ""
	if err := json.Unmarshal(failResp.Error, &errMsg); err != nil {
		log.Fatal(err)
	}

	assert.NotEqual(t, errMsg, "")
}

func TestUpdateLocker(t *testing.T) {
	TI.Reset()
	setLockersMock()
	// 정상 요청
	url := TI.Server.URL + "/locker"
	method := "PATCH"

	data := models.Locker{
		ID:       1,
		Location: "AAA",
	}
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

	succMsg := ""
	if err := json.Unmarshal(succResp.Data, &succMsg); err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, succMsg, "success")

	// AAA가 존재하는지 확인
	url = TI.Server.URL + "/locker/1"
	method = "GET"

	emptyData := map[string]string{}
	headers = map[string]string{
		"Authorization": "Bearer " + TI.Bearers["normal1"],
		"Content-Type":  "Application/json",
	}

	statusCode, body, err = utils.TestRequest(
		url, method, emptyData, headers, nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, statusCode, 200)

	// 요청 결과 확인
	succResp = controllers.SuccessResponse{}
	if err := json.Unmarshal(body, &succResp); err != nil {
		log.Fatal(err)
	}

	locker := models.Locker{}
	if err := json.Unmarshal(succResp.Data, &locker); err != nil {
		log.Fatal(err)
	}

	assert.NotEqual(t, locker.ID, 0)
	assert.Equal(t, locker.ID, 1)
	assert.Equal(t, locker.Location, "AAA")

	// 없는 ID 갱신 시도
	url = TI.Server.URL + "/locker"
	method = "PATCH"

	data = models.Locker{
		ID:       999,
		Location: "Not exists",
	}
	headers = map[string]string{
		"Authorization": "Bearer " + TI.Bearers["normal1"],
		"Content-Type":  "Application/json",
	}

	statusCode, body, err = utils.TestRequest(
		url, method, data, headers, nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, statusCode, 200)

	// 요청 결과 확인
	failResp := controllers.SuccessResponse{}
	if err := json.Unmarshal(body, &failResp); err != nil {
		log.Fatal(err)
	}

	succMsg = ""
	if err := json.Unmarshal(failResp.Data, &succMsg); err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, succMsg, "success")
}

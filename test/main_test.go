package test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"simple-api-example/auth"
	"simple-api-example/controllers"
	"simple-api-example/database"
	"simple-api-example/models"
	"simple-api-example/router"
	"simple-api-example/utils"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 테스트 기본 정보
var TI TestInfo

// testDB : 테스트 DB Connection
var testDB *gorm.DB

// TestUserInput : 테스트 유저 인풋
type TestUserInput struct {
	models.User
	Password string `json:"Password"`
}

// TestInfo : 테스트 기본 정보 구조체
type TestInfo struct {
	Bearers    map[string]string
	Server     *httptest.Server
	UserInputs []TestUserInput
}

// Init : 테스트 정보 초기화
func (ti *TestInfo) Init() {
	utils.LoadEnv()
	initTestDB()
	ti.setTestServer()
	ti.createTestUsers()
	ti.setBearers()
}

// Reset : 테스트 정보 리셋
func (ti *TestInfo) Reset() {
	models.DeleteAllTables(testDB)
	ti.createTestUsers()
	ti.setBearers()
}

// initTestDB : testDB 초기화
func initTestDB() {
	var (
		_, b, _, _     = runtime.Caller(0)
		currentDirPath = filepath.Dir(b)
	)
	testDBPath := fmt.Sprintf("%s/%s", currentDirPath, "test.db")

	// test.db가 존재할 시 삭제
	if _, err := os.Stat(testDBPath); err == nil {
		err := os.Remove(testDBPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	db, err := gorm.Open(sqlite.Open(testDBPath), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	database.DB = db
	testDB = db

	models.CreateAllTablesIfNotExists(testDB)
}

// setTestServer : 테스트 서버 생성
func (ti *TestInfo) setTestServer() {
	r := router.SetupRouter()
	ti.Server = httptest.NewServer(r)
}

//createTestUser : 테스트 유저 생성
func (ti *TestInfo) createTestUsers() {
	ti.UserInputs = []TestUserInput{
		{User: models.User{Name: "normal1", Group: "normal"}, Password: "normalpassword1"},
		{User: models.User{Name: "admin1", Group: "admin"}, Password: "adminpassword1"},
	}

	users := models.Users{}
	for i, userInput := range ti.UserInputs {
		hash, err := utils.HashAndSalt(userInput.Password)
		if err != nil {
			log.Fatal(err.Error())
		}
		ti.UserInputs[i].User.SecretKey = hash
		users = append(users, ti.UserInputs[i].User)
	}

	err := users.Create()
	if err != nil {
		log.Fatalf(err.Error())
	}
}

// setBearers : 테스트 유저의 auth token 발급
func (ti *TestInfo) setBearers() {
	TI.Bearers = map[string]string{}

	for _, userInput := range ti.UserInputs {
		encoded := auth.GetBasicAuth(userInput.User.Name, userInput.Password)
		headers := map[string]string{
			"Authorization": "Basic " + encoded,
			"Content-Type":  "Application/json",
		}
		statusCode, body, err := utils.TestRequest(ti.Server.URL+"/login/basic", "POST", map[string]string{}, headers, nil)
		if err != nil {
			log.Fatal(err)
		}
		if statusCode != 200 {
			log.Fatalf("failed to get jwt. body: %s", body)
		}

		succResp := controllers.SuccessResponse{}
		err = json.Unmarshal(body, &succResp)
		if err != nil {
			log.Fatalln(err)
		}

		authTokenResp := controllers.AuthTokenResponse{}
		err = json.Unmarshal(succResp.Data, &authTokenResp)
		if err != nil {
			log.Fatalln(err)
		}

		ti.Bearers[userInput.User.Name] = authTokenResp.Token
	}
}

func TestMain(m *testing.M) {
	TI.Init()
	code := m.Run()
	defer TI.Server.Close()
	os.Exit(code)
}

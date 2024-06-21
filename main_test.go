package main_test

import (
	"bytes"
	"devdiaries/models"
	"devdiaries/payload/request"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/lpernett/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var test_user models.User
var ss string

func setup() {
	godotenv.Load(".env.test.local")

	db, err := gorm.Open(mysql.Open(os.Getenv("DB_URL")), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	DB = db

	result := DB.Where(&models.User{Email: os.Getenv("TEST_USER_EMAIL")}).First(&test_user)

	if result.Error != nil {
		panic(result.Error)
	}
}

func shutdown() {

}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func TestLogin(t *testing.T) {

	body, err := json.Marshal(request.LoginBody{Email: test_user.Email, Password: os.Getenv("TEST_USER_PASSWORD")})

	if err != nil {
		t.Errorf("Error creating POST request body: %s", err.Error())
		return
	}

	url, err := url.JoinPath(os.Getenv("SERVER_HOST"), "/login")

	if err != nil {
		t.Errorf("Error building url: %s", err.Error())
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))

	if err != nil {
		t.Errorf("Error sending POST request: %s", err.Error())
		return
	}

	if resp.StatusCode != 200 {
		t.Errorf("Login failed with status code: %d", resp.StatusCode)
		return
	}

	t.Logf("Successfull login, cookie: %s", resp.Cookies()[0].Value)
	ss = resp.Cookies()[0].Value

}

func TestMalformedAuth(t *testing.T) {

	url, err := url.JoinPath(os.Getenv("SERVER_HOST"), "/user/1")

	if err != nil {
		t.Errorf("Error building url: %s", err.Error())
		return
	}

	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		t.Errorf("Error creating request object: %s", err.Error())
		return
	}

	request.Header.Add("Authorization", "Bearer abcdef1234")

	client := http.Client{}
	resp, err := client.Do(request)

	if err != nil {
		t.Errorf("Error making GET request: %s", err.Error())
		return
	}

	if resp.StatusCode != 400 {
		t.Errorf("Expected 400 BAD request instead got %s", resp.Status)
	}

}

func TestValidAuth(t *testing.T) {
	url, err := url.JoinPath(os.Getenv("SERVER_HOST"), "/user/1")

	if err != nil {
		t.Errorf("Error building url: %s", err.Error())
		return
	}

	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		t.Errorf("Error creating request object: %s", err.Error())
		return
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ss))

	client := http.Client{}
	resp, err := client.Do(request)

	if err != nil {
		t.Errorf("Error making GET request: %s", err.Error())
		return
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected 200 OK instead got %s", resp.Status)
	}
}

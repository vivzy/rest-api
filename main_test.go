package main

import (
  "bytes"
  "github.com/gin-gonic/gin"
  "net/http"
  "net/http/httptest"
  "testing"
	"math/rand"
	"github.com/gin-gonic/gin/json"
	"log"
	"fmt"
	"rest-api/app/models"
	"github.com/rs/xid"
)

var username = xid.New().String()
var password = RandString(6)
var firstname = RandString(10)
var lastname = RandString(5)
var joindate = "2018-01-01"
var mcPostBody = map[string]interface{}{
"username": username,
"password": password,
"firstname": firstname,
"lastname": lastname,
"joindate": joindate,
}

var tokenResponse *TokenResponse
var createdUser *models.User

type TokenResponse struct {
	Code int `json:"code"`
	Response string `json:"response"`
	Token string  `json:"token"`
}

func getToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testRouter := SetupRouter()
	postBody := map[string]interface{}{
		"username": username,
		"password": password,
	}
	body, _ := json.Marshal(postBody)
	req, err := http.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Errorf("Post hearteat failed with error %d.", err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	if resp.Code != 200 {
		t.Errorf("/login failed with error code %d.", resp.Code)
	}
	respRaw := bytes.NewReader(resp.Body.Bytes())
	decoder := json.NewDecoder(respRaw)

	decodeErr := decoder.Decode(&tokenResponse)

	if decodeErr != nil {
		log.Printf("decode error")
		log.Fatal(err)
	}
}

func testGetUser(t *testing.T) {
  gin.SetMode(gin.TestMode)
  testRouter := SetupRouter()
  log.Printf("fetching user with id %d", createdUser.Id)
  req, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%d", createdUser.Id), nil)
  req.Header.Set("Authorization", "Bearer " + tokenResponse.Token)
  if err != nil {
    t.Errorf("Get hearteat failed with error %d.", err)
  }

  resp := httptest.NewRecorder()
  testRouter.ServeHTTP(resp, req)

  if resp.Code != 200 {
    t.Errorf("/api/v1/users failed with error code %d.", resp.Code)
  }
}

func testGetUsers(t *testing.T) {
  gin.SetMode(gin.TestMode)
  testRouter := SetupRouter()

  req, err := http.NewRequest("GET", "/api/v1/users", nil)
  log.Printf(" auth token %s", tokenResponse.Token )
  req.Header.Set("Authorization", "Bearer " + tokenResponse.Token)
  if err != nil {
    t.Errorf("Get hearteat failed with error %d.", err)
  }

  resp := httptest.NewRecorder()
  testRouter.ServeHTTP(resp, req)

  if resp.Code != 200 {
    t.Errorf("/api/v1/users failed with error code %d.", resp.Code)
  }
}

func testPostUser(t *testing.T) {
  gin.SetMode(gin.TestMode)
  testRouter := SetupRouter()

  body, _ := json.Marshal(mcPostBody)
  req, err := http.NewRequest("POST", "/adduser", bytes.NewReader(body))
  req.Header.Set("Content-Type", "application/json")
  if err != nil {
    t.Errorf("Post hearteat failed with error %d.", err)
  }

  resp := httptest.NewRecorder()
  testRouter.ServeHTTP(resp, req)

  if resp.Code != 201 {
    t.Errorf("/api/v1/users failed with error code %d.", resp.Code)
  }
  fmt.Printf(resp.Body.String())
  respRaw := bytes.NewReader(resp.Body.Bytes())
  decoder := json.NewDecoder(respRaw)

  decodeErr := decoder.Decode(&createdUser)
  fmt.Printf("created user with id %d", createdUser.Id)
  if decodeErr != nil {
	  log.Printf("decode error")
	  log.Fatal(err)
  }
}

func testPostSameUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testRouter := SetupRouter()

	body, _ := json.Marshal(mcPostBody)
	req, err := http.NewRequest("POST", "/adduser", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Errorf("Post hearteat failed with error %d.", err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	if resp.Code != 500 {
		t.Errorf("/api/v1/users failed with error code %d.", resp.Code)
	}
}

func testPutUser(t *testing.T) {
  gin.SetMode(gin.TestMode)
  testRouter := SetupRouter()

  body := bytes.NewBuffer([]byte("{\"firstname\": \"updatedname\", \"lastname\": \"updatedlastname\"}"))

  req, err := http.NewRequest("PUT", fmt.Sprintf("/api/v1/users/%d", createdUser.Id), body)
	req.Header.Set("Authorization", "Bearer " + tokenResponse.Token)
  req.Header.Set("Content-Type", "application/json")
  if err != nil {
    t.Errorf("Put hearteat failed with error %d.", err)
  }

  resp := httptest.NewRecorder()
  testRouter.ServeHTTP(resp, req)

  if resp.Code >= 500 {
    t.Errorf("/api/v1/users failed with error code %d.", resp.Code)
  }
}

func TestAllCases(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("adding a new user", testPostUser)
	t.Run("adding the same user ", testPostSameUser)
	t.Run("login as created user", getToken)
	t.Run("get all users", testGetUsers)
	t.Run("get a user", testGetUser)
	t.Run("update a user", testPutUser)

}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// A helper function to generate random string
func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}


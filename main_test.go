package main

import (
  "bytes"
  "github.com/gin-gonic/gin"
  "net/http"
  "net/http/httptest"
  "testing"
	"math/rand"
	"github.com/gin-gonic/gin/json"
)

func getToken() (u, v string){
	//return map[string][]string{"Authorization":{"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MjE1NzE5NzgsImlkIjoiYWdlbnRodWdzMyIsIm9yaWdfaWF0IjoxNTIxNTY4Mzc4fQ.vAsxl5ttH59g_eNBUhWPSgl5Fooml91pwc-CxHrUEGo"}}
	return "Authorization","Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MjE1NzE5NzgsImlkIjoiYWdlbnRodWdzMyIsIm9yaWdfaWF0IjoxNTIxNTY4Mzc4fQ.vAsxl5ttH59g_eNBUhWPSgl5Fooml91pwc-CxHrUEGo"
}

func TestGetUser(t *testing.T) {
  gin.SetMode(gin.TestMode)
  testRouter := SetupRouter()

  req, err := http.NewRequest("GET", "/api/v1/users/5", nil)
  req.Header.Set(getToken())
  if err != nil {
    t.Errorf("Get hearteat failed with error %d.", err)
  }

  resp := httptest.NewRecorder()
  testRouter.ServeHTTP(resp, req)

  if resp.Code != 200 {
    t.Errorf("/api/v1/users failed with error code %d.", resp.Code)
  }
}

func TestGetUsers(t *testing.T) {
  gin.SetMode(gin.TestMode)
  testRouter := SetupRouter()

  req, err := http.NewRequest("GET", "/api/v1/users", nil)
  req.Header.Set(getToken())
  if err != nil {
    t.Errorf("Get hearteat failed with error %d.", err)
  }

  resp := httptest.NewRecorder()
  testRouter.ServeHTTP(resp, req)

  if resp.Code != 200 {
    t.Errorf("/api/v1/users failed with error code %d.", resp.Code)
  }
}

func TestPostUser(t *testing.T) {
  gin.SetMode(gin.TestMode)
  testRouter := SetupRouter()
  username := RandString(10)
  password := RandString(6)
  firstname := RandString(10)
  lastname := RandString(5)
  joindate := "2018-01-01"
  //println("username : "+ username)
  //println("password : "+ password)
  mcPostBody := map[string]interface{}{
  	"username": username,
  	"password": password,
  	"firstname": firstname,
  	"lastname": lastname,
  	"joindate": joindate,
  }
  body, _ := json.Marshal(mcPostBody)
  //body := bytes.NewBuffer([]byte("{\"username\": \"" + username+ "\", \"password\": \""+ password + "\", \firstname\": \""+ firstname+ "\"" +
  	//"\"lastname\": \""+lastname+"\",\"joindate\": \""+joindate+"\"}"))

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
}

func TestPostSameUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testRouter := SetupRouter()
	username := "agenthugs1"
	password := RandString(6)
	firstname := RandString(10)
	lastname := RandString(5)
	joindate := "2018-01-01"
	mcPostBody := map[string]interface{}{
		"username": username,
		"password": password,
		"firstname": firstname,
		"lastname": lastname,
		"joindate": joindate,
	}
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

func TestPutUser(t *testing.T) {
  gin.SetMode(gin.TestMode)
  testRouter := SetupRouter()

  body := bytes.NewBuffer([]byte("{\"firstname\": \""+RandString(10)+"\", \"lastname\": \""+RandString(5)+"\"}"))

  req, err := http.NewRequest("PUT", "/api/v1/users/5", body)
  req.Header.Set(getToken())
  req.Header.Set("Content-Type", "application/json")
  if err != nil {
    t.Errorf("Put hearteat failed with error %d.", err)
  }

  resp := httptest.NewRecorder()
  testRouter.ServeHTTP(resp, req)

  if resp.Code != 200 {
    t.Errorf("/api/v1/users failed with error code %d.", resp.Code)
  }
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


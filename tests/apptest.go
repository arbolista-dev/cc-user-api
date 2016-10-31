package tests

import (
	"encoding/json"
	"github.com/revel/revel/testing"
	"io"
	"log"
	"net/http"
	"strings"
	"strconv"
)

var token, userID string

type AppTest struct {
	testing.TestSuite
}

type apiResult struct {
	Success bool        `json:"success"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func myVERB(verb, path string, contentType string, reader io.Reader, token string, t *AppTest) *http.Request {
	var err error
	var req *http.Request
	switch verb {
	case "POST":
		req, err = http.NewRequest("POST", t.BaseUrl()+path, reader)
		req.Header.Set("Content-Type", contentType)
	case "PUT":
		req, err = http.NewRequest("PUT", t.BaseUrl()+path, reader)
		req.Header.Set("Content-Type", contentType)
	case "GET":
		req, err = http.NewRequest("GET", t.BaseUrl()+path, nil)
	case "DELETE":
		req, err = http.NewRequest("DELETE", t.BaseUrl()+path, nil)
	}
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", token)
	return req
}

func testSuccess(t *AppTest, pass bool, errMessage string) {
	var result apiResult
	err := json.Unmarshal(t.ResponseBody, &result)
	t.AssertEqual(err, nil)
	t.AssertEqual(result.Success, pass)
	if pass == false {
		t.AssertEqual(result.Error, errMessage)
	}
}

// --------------- TEST FUNCTIONS -------------

func (t *AppTest) TestA_Add_SUCCESS() {
	t.Post("/user", "application/json; charset=utf-8", strings.NewReader(userBody))
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
	log.Println(string(t.ResponseBody))
	testSuccess(t, true, "")
}

func (t *AppTest) TestB_Add_ERROR_DuplicateEmail() {
	t.Post("/user", "application/json; charset=utf-8", strings.NewReader(userBody))
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
	log.Println(string(t.ResponseBody))
	testSuccess(t, false, `{"email": "non-unique"}`)
}

func (t *AppTest) TestC_Login_SUCCESS() {
	t.Post("/user/login", "application/json; charset=utf-8", strings.NewReader(loginBody))
	buf := t.ResponseBody
	var logRes apiResult
	err := json.Unmarshal(buf, &logRes)
	t.AssertEqual(err, nil)
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
	if logRes.Data != nil {
		token = logRes.Data.(map[string]interface{})["token"].(string)
		log.Println(string(t.ResponseBody))
		log.Println("Setting TOKEN to: " + token)
	}
	testSuccess(t, true, "")
}

func (t *AppTest) TestD_Login_ERROR_BadPassword() {
	t.Post("/user/login", "application/json; charset=utf-8", strings.NewReader(loginBody_badPassword))
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
	log.Println(string(t.ResponseBody))
	testSuccess(t, false, `{"password": "incorrect"}`)
}

func (t *AppTest) TestD_Login_ERROR_BadEmail() {
	t.Post("/user/login", "application/json; charset=utf-8", strings.NewReader(loginBody_badEmail))
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
	log.Println(string(t.ResponseBody))
	testSuccess(t, false, `{"email": "non-existent"}`)
}

func (t *AppTest) TestE1_Update_SUCCESS() {
	req := myVERB("PUT", "/user", "application/json; charset=utf-8", strings.NewReader(userBody_update), token, t)
	t.NewTestRequest(req).Send()
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
	log.Println(string(t.ResponseBody))
	testSuccess(t, true, "")
}

func (t *AppTest) TestE2_UpdateAnswers_SUCCESS() {
	req := myVERB("PUT", "/user/answers", "application/json; charset=utf-8", strings.NewReader(answers_update), token, t)
	t.NewTestRequest(req).Send()
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
	log.Println(string(t.ResponseBody))
	testSuccess(t, true, "")
}

func (t *AppTest) TestE3_SetLocation_SUCCESS() {
	req := myVERB("PUT", "/user/location", "application/json; charset=utf-8", strings.NewReader(location_set), token, t)
	t.NewTestRequest(req).Send()
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
	log.Println(string(t.ResponseBody))
	testSuccess(t, true, "")
}

func (t *AppTest) TestE_UserLogout_SUCCESS() {
	req := myVERB("GET", "/user/logout", "", nil, token, t)
	t.NewTestRequest(req).Send()
	log.Println(string(t.ResponseBody))
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
	testSuccess(t, true, "")
}

func (t *AppTest) TestF1_ListLeaders_SUCCESS() {
	req := myVERB("GET", "/user/leaders", "", nil, "", t)
	t.NewTestRequest(req).Send()
	buf := t.ResponseBody
	var listRes apiResult
	err := json.Unmarshal(buf, &listRes)
	t.AssertEqual(err, nil)
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
	if listRes.Data != nil {
		_userID := int(listRes.Data.(map[string]interface{})["list"].([]interface{})[0].(map[string]interface{})["user_id"].(float64))
		userID = strconv.Itoa(_userID)
		log.Println(string(t.ResponseBody))
		log.Printf("Setting userID to: ", userID)
	}
	testSuccess(t, true, "")
}

func (t *AppTest) TestF2_ShowUser_SUCCESS() {
	userPath := "/profile/" + userID
	req := myVERB("GET", userPath, "", nil, "", t)
	t.NewTestRequest(req).Send()
	log.Println(string(t.ResponseBody))
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
	testSuccess(t, true, profile)
}

func (t *AppTest) TestF3_ShowUser_ERROR_NonExistent() {
	userPath := "/profile/999999999"
	req := myVERB("GET", userPath, "", nil, "", t)
	t.NewTestRequest(req).Send()
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
	log.Println(string(t.ResponseBody))
	testSuccess(t, false, `{"profile": "non-existent"}`)
}

func (t *AppTest) TestG_UserLogout_ERROR_NoSession() {
	req := myVERB("GET", "/user/logout", "", nil, token, t)
	t.NewTestRequest(req).Send()
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
	log.Println(string(t.ResponseBody))
	testSuccess(t, false, `{"session": "non-existent"}`)
}

func (t *AppTest) TestH_UserLogin_SUCCESS() {
	t.TestC_Login_SUCCESS()
}

func (t *AppTest) TestI_Delete_SUCCESS() {
	req := myVERB("DELETE", "/user", "", nil, token, t)
	t.NewTestRequest(req).Send()
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
	log.Println(string(t.ResponseBody))
	testSuccess(t, true, "")
}

func (t *AppTest) Before() {
	log.Println("+++++++++++++++++++++++++++++++++++++++++++++++++")
}

func (t *AppTest) After() {
	log.Println("-------------------------------------------------")
}

package tests

import (
	"encoding/json"
	"github.com/revel/revel/testing"
	"io"
	"os"
	"log"
	"net/http"
	"net/textproto"
	"mime/multipart"
	"bytes"
	"strings"
	"github.com/arbolista-dev/cc-user-api/app/services"
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

func fileUploadRequest(path string, filepath string, token string, t *AppTest) *http.Request {
	var err error
	var req *http.Request

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	log.Println("file", file)

	h := textproto.MIMEHeader{}
	h.Set("Content-Type", "multipart/form-data")

	part, err := writer.CreatePart(h)
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		panic(err)
	}

	err = writer.Close()
	t.AssertEqual(nil, err)

	req, err = http.NewRequest("POST", t.BaseUrl()+path, body)
	req.Header.Set("Authorization", token)
	if err != nil {
		panic(err)
	}
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

func (t *AppTest) TestE4_UpdateUserGoals_SUCCESS() {
	req := myVERB("PUT", "/user/goals", "application/json; charset=utf-8", strings.NewReader(user_goals_update), token, t)
	t.NewTestRequest(req).Send()
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
	log.Println(string(t.ResponseBody))
	testSuccess(t, true, "")
}

// func (t *AppTest) TestE5_SetPhoto_SUCCESS() {
// 	filePath, _ := os.Getwd()
// 	filePath += "/tests/profile-photo.jpg"
// 	log.Println("filepath", filePath)
// 	req := fileUploadRequest("/user/photo", filePath, token, t)
// 	t.NewTestRequest(req).Send()
// 	t.AssertOk()
// 	log.Println(string(t.ResponseBody))
// 	testSuccess(t, true, "")
// }

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
		log.Printf("Setting userID for retrieving profile to: ", userID)
	}
	testSuccess(t, true, "")
}

func (t *AppTest) TestF2_ShowUserProfile_SUCCESS() {
	userPath := "/user/" + userID + "/profile"
	req := myVERB("GET", userPath, "", nil, "", t)
	t.NewTestRequest(req).Send()
	log.Println(string(t.ResponseBody))
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
	testSuccess(t, true, profile)
}

func (t *AppTest) TestF3_ShowUserProfile_ERROR_NonExistent() {
	userPath := "/user/999999999/profile"
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

func (t *AppTest) TestI_RetrieveUserGoals_SUCCESS() {
	req := myVERB("GET", "/user/goals", "", nil, token, t)
	t.NewTestRequest(req).Send()
	buf := t.ResponseBody
	var listRes apiResult
	err := json.Unmarshal(buf, &listRes)
	t.AssertEqual(err, nil)
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
  log.Printf("Retrieve user goals: ", listRes.Data)
	if listRes.Data != nil {
		status := listRes.Data.(map[string]interface{})["list"].([]interface{})[0].(map[string]interface{})["status"].(string)
		t.AssertEqual(status, "pledged")
	}
	testSuccess(t, true, "")
}

func (t *AppTest) TestJ_Delete_SUCCESS() {
	req := myVERB("DELETE", "/user", "", nil, token, t)
	t.NewTestRequest(req).Send()
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
	log.Println(string(t.ResponseBody))
	testSuccess(t, true, "")
}

// func (t *AppTest) TestK1_ConfirmMail(){
// 	data := map[string]string{"-name-": "prueba", "-link-":"http://www.google.com"}
// 	err := services.SendMail("confirm", "test@sink.sendgrid.net", data)
// 	if err != nil {
// 		log.Print(err)
// 	}
// 	t.AssertEqual(err, nil)
// }

func (t *AppTest) TestK2_ResetMail(){
	data := map[string]string{"-name-": "prueba", "-link-":"http://www.google.com"}
	err := services.SendMail("reset", "test@sink.sendgrid.net", data)
	if err != nil {
		log.Print(err)
	}
	t.AssertEqual(err, nil)
}


func (t *AppTest) Before() {
	log.Println("+++++++++++++++++++++++++++++++++++++++++++++++++")
}

func (t *AppTest) After() {
	log.Println("-------------------------------------------------")
}

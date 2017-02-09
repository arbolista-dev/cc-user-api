package controllers

import (
	"encoding/json"
	"github.com/arbolista-dev/cc-user-api/app/ds"
	"github.com/arbolista-dev/cc-user-api/app/models"
	"github.com/arbolista-dev/cc-user-api/app/services"
	"github.com/revel/revel"
	"io/ioutil"
	"net/url"
	"strconv"
	"math"
	"strings"
)

type Users struct {
	App
}

func (c Users) LoginFacebook() revel.Result {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return c.Error(err)
	}
	var logRequest models.UserFacebook
	err = json.Unmarshal(body, &logRequest)
	if err != nil {
		return c.Error(err)
	}
	login, err := ds.LoginFacebook(logRequest)
	if err != nil {
		return c.Error(err)
	}
	return c.Data(login)

}

func (c Users) Login() revel.Result {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return c.Error(err)
	}
	var logRequest models.UserLogin
	err = json.Unmarshal(body, &logRequest)
	if err != nil {
		return c.Error(err)
	}
	login, err := ds.Login(logRequest)
	if err != nil {
		return c.Error(err)
	}

	return c.Data(login)
}

func (c Users) Logout() revel.Result {
	userID, jti, err := c.GetSession()
	if err != nil {
		return c.Error(err)
	}

	err = ds.Logout(userID, jti)
	if err != nil {
		return c.Error(err)
	}
	return c.OK()
}

func (c Users) LogoutAll() revel.Result {
	userID, _, err := c.GetSession()
	if err != nil {
		return c.Error(err)
	}

	err = ds.LogoutAll(userID)
	if err != nil {
		return c.Error(err)
	}
	return c.OK()
}

func (c Users) ListLeaders() revel.Result {

	limit, offset, state, household_size := 10, 0, "", 0

	if len(c.Params.Values) != 0 {
		if value, ok := c.Params.Values["limit"]; ok {
			v, err := strconv.ParseInt(value[0], 10, 32)
			if err != nil {
				return c.Error(err)
			}
			limit = int(v)
		}
		if value, ok := c.Params.Values["offset"]; ok {
			v, err := strconv.ParseInt(value[0], 10, 32)
			if err != nil {
				return c.Error(err)
			}
			offset = int(v)
		}
		if value, ok := c.Params.Values["state"]; ok {
			state = value[0]
		}
		if value, ok := c.Params.Values["household_size"]; ok {
			v, err := strconv.ParseInt(value[0], 10, 32)
			if err != nil {
				return c.Error(err)
			}
			household_size = int(v)
		}
	}

	leaders, err := ds.ListLeaders(limit, offset, state, household_size)
	if err != nil {
		return c.Error(err)
	}

	return c.Data(leaders)
}

func (c Users) ListLocations() revel.Result {

	locations, err := ds.ListLocations()
	if err != nil {
		return c.Error(err)
	}

	return c.Data(locations)
}

func (c Users) Add() revel.Result {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return c.Error(err)
	}
	var newUser models.User
	err = json.Unmarshal(body, &newUser)
	if err != nil {
		return c.Error(err)
	}
	newUser.Validate(c.Validation)
	if c.Validation.HasErrors() {
		errors := c.Validation.ErrorMap()
		return c.ErrorData(errors)
	}

	login, userID, err := ds.Add(newUser)
	if err != nil {
		return c.Error(err)
	}
	token, err := ds.ConfirmRequest(newUser.Email)
	if err != nil {
		return c.Error(err)
	}
	err = SendUserMail(newUser.FirstName, userID, token, newUser.Email);
	if err != nil {
		return c.Error(err)
	}
	return c.Data(login)
}



func (c Users) Delete() revel.Result {
	userID, _, err := c.GetSession()
	if err != nil {
		return c.Error(err)
	}

	err = ds.Delete(userID)
	if err != nil {
		return c.Error(err)
	}
	return c.OK()
}

func (c Users) UpdateAnswers() revel.Result {
	userID, _, err := c.GetSession()
	if err != nil {
		return c.Error(err)
	}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return c.Error(err)
	}
	var bodyAnswers models.Answers
	err = json.Unmarshal(body, &bodyAnswers)
	if err != nil {
		return c.Error(err)
	}

	var answersMap map[string]interface{}
	err = json.Unmarshal([]byte(body), &answersMap)
	if err != nil {
		return c.Error(err)
	}

	footprintsMap := map[string]interface{}{
		"result_food_total": 0,
		"result_housing_total": 0,
		"result_shopping_total": 0,
		"result_transport_total": 0,
		"result_grand_total": 0,
	}

	for name, _ := range footprintsMap {
		if name == "result_shopping_total" {
			footprintsMap[name] = FootprintAnswerToUint("result_services_total", answersMap) + FootprintAnswerToUint("result_goods_total", answersMap)
		} else {
		  footprintsMap[name] = FootprintAnswerToUint(name, answersMap)
		}
	}

	var footprint models.TotalFootprint
	footprint.TotalFootprint, err = json.Marshal(footprintsMap)
	if err != nil {
		return c.Error(err)
	}

	var userAnswers models.AnswersUpdate

	userAnswers.Answers = bodyAnswers.Answers
	userAnswers.TotalFootprint = footprint.TotalFootprint

	input_size, ok := answersMap["answers"].(map[string]interface{})["input_size"].(string)
	if ok == true {
		householdSize, err := strconv.Atoi(input_size)
		if err != nil {
			return c.Error(err)
		}
		userAnswers.HouseholdSize = householdSize
	}

	err = ds.UpdateAnswers(userID, userAnswers)
	if err != nil {
		return c.Error(err)
	}
	return c.OK()
}

func FootprintAnswerToUint(name string, answersMap map[string]interface{}) (footprintAmount uint) {
	amount, ok := answersMap["answers"].(map[string]interface{})[name].(string)
	if ok == true {
		amountFloat, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			return
		}
		return uint(math.Floor(amountFloat + .5))
	} else {
		return
	}
}

func (c Users) SetLocation() revel.Result {
	userID, _, err := c.GetSession()
	if err != nil {
		return c.Error(err)
	}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return c.Error(err)
	}

	var bodyLocation models.Location
	err = json.Unmarshal(body, &bodyLocation)
	if err != nil {
		return c.Error(err)
	}

	err = ds.SetLocation(userID, bodyLocation)
	if err != nil {
		return c.Error(err)
	}
	return c.OK()
}

func (c Users) Update() revel.Result {
	userID, _, err := c.GetSession()
	if err != nil {
		return c.Error(err)
	}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return c.Error(err)
	}
	var user models.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		return c.Error(err)
	}

	err = ds.Update(userID, user)
	if err != nil {
		return c.Error(err)
	}
	return c.OK()
}


func (c Users) PassResetRequest() revel.Result {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return c.Error(err)
	}

	var email models.Email
	err = json.Unmarshal(body, &email)
	if err != nil {
		return c.Error(err)
	}

	userID, token, name, err := ds.PassResetRequest(email.Email)
	if err != nil {
		return c.Error(err)
	}

	data := map[string]string{"-link-": PasswordResetURL(userID, token),"-name-": name}
	err = services.SendMail("reset",email.Email, data)
	if err != nil {
		return c.Error(err)
	}
	return c.OK()
}

func (c Users) PassResetConfirm() revel.Result {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return c.Error(err)
	}

	var reset models.PasswordReset
	err = json.Unmarshal(body, &reset)
	if err != nil {
		return c.Error(err)
	}
	err = ds.PassResetConfirm(reset.Id, reset.Token, reset.Password)
	if err != nil {
		return c.Error(err)
	}else {
		return c.OK()
	}
}

func (c Users) PasswordReset(id uint, token string) revel.Result {
	// req, err = http.NewRequest("POST", c.BaseUrl()mt+"/user/reset", c.RenderJson(reset))
	// req.Header.Set("Content-Type", contentType)
	// fmt.Println("EL pass: ", reset.Password)
	return c.Render(id, token)
}

func (c Users) Confirm(id uint,token string) revel.Result {
	// req, err = http.NewRequest("POST", c.BaseUrl()mt+"/user/reset", c.RenderJson(reset))
	// req.Header.Set("Content-Type", contentType)
	// fmt.Println("EL pass: ", reset.Password)
	err := ds.ConfirmEmail(id, token)
	host := revel.Config.StringDefault("server.confirm.host","127.0.0.1:3000")
	if err != nil {
		return c.Redirect("http://"+host+"/en/settings?type=confirm&error="+err.Error())
	}else {
		return c.Redirect("http://"+host+"/en/settings?type=confirm")
	}
}
func (c Users) NeedActivate() revel.Result {
	userID, _, err := c.GetSession()
	if err != nil {
		return c.Error(err)
	}
	activate, err := ds.NeedActivate(userID)
	if err != nil {
		return c.Error(err)
	}
	return c.Data(activate)
}

func (c Users) SendActivate() revel.Result {
	userID, _, err := c.GetSession()
	if err != nil {
		return c.Error(err)
	}
	activate, err := ds.NeedActivate(userID)
	if err != nil {
		return c.Error(err)
	}
	if activate.Need {
		token, err := ds.ConfirmRequest(activate.Email)
		if err != nil {
			return c.Error(err)
		}
		err = SendUserMail(activate.Name, userID, token, activate.Email )
	}
	return c.Data(activate)
}

func SendUserMail(name string, userID uint, token string, email string ) (err error) {
	data := map[string]string{"-name-": name, "-link-": ConfirmURL( userID, token)}
	err = services.SendMail("confirm", email, data)
	return
}

func ConfirmURL(userID uint, token string) (uri string) {
	var host string
	var scheme string
	//Get the host and port from configuration
	host = revel.Config.StringDefault("http.addr","127.0.0.1") + revel.Config.StringDefault("http.port","9000")
	//If we set server.host this override the host name in url
	host = revel.Config.StringDefault("server.reset.host",host)
	if revel.Config.BoolDefault("http.ssl",false) {
		scheme = "https"
	} else {
		scheme = "http"
	}
	u := url.URL{}
	u.Scheme =  scheme
	u.Host = host
	u.Path = "/user/confirm"
	q := u.Query()
	q.Set("id", strconv.Itoa(int(userID)))
	q.Set("token", token)
	u.RawQuery = q.Encode()
	return u.String()
}

func PasswordResetURL(userID uint, token string) (uri string) {
	var host string
	var scheme string
	host = revel.Config.StringDefault("http.addr","127.0.0.1")
	if strings.TrimSpace(host) == "" {
		host = "127.0.0.1"
	}
	//Get the host and port from configuration
	host = host +":"+ revel.Config.StringDefault("http.port","9000")
	//If we set server.host this override the host name in url
	host = revel.Config.StringDefault("server.reset.host",host)
	if revel.Config.BoolDefault("http.ssl",false) {
		scheme = "https"
	}else {
		scheme = "http"
	}
	u := url.URL{}
	u.Scheme =  scheme
	u.Host = host
	u.Path = "/en/settings"
	q := u.Query()
	q.Set("type","reset")
	q.Set("id", strconv.Itoa(int(userID)))
	q.Set("token", token)
	u.RawQuery = q.Encode()
	return u.String()
}

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

	limit, offset, category, city, state := 10, 0, "", "", ""

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
		if value, ok := c.Params.Values["category"]; ok {
			category = value[0]
		}
		if value, ok := c.Params.Values["city"]; ok {
			city = value[0]
		}
		if value, ok := c.Params.Values["state"]; ok {
			state = value[0]
		}
	}

	leaders, err := ds.ListLeaders(limit, offset, category, city, state)
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

	login, err := ds.Add(newUser)
	if err != nil {
		return c.Error(err)
	}

	// Disable send email notification
	// data := map[string]string{"name": newUser.FirstName}
	// err = services.SendMail("new-user-beta", newUser.Email, data)
	// if err != nil {
	// 	return c.Error(err)
	// }

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

	var updatedUser models.User

	updatedUser.TotalFootprint = footprint.TotalFootprint
	updatedUser.Answers = bodyAnswers.Answers

	input_size, ok := answersMap["answers"].(map[string]interface{})["input_size"].(string)
	if ok {
		householdSize, err := strconv.Atoi(input_size)
		if err != nil {
			return c.Error(err)
		}
		updatedUser.HouseholdSize = householdSize
	}

	err = ds.Update(userID, updatedUser)
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

	userID, token, err := ds.PassResetRequest(email.Email)
	if err != nil {
		return c.Error(err)
	}

	data := map[string]string{"link": PasswordResetURL(userID, token)}
	err = services.SendMail("passwords-reset", email.Email, data)
	if err != nil {
		return c.Error(err)
	}
	return c.OK()
}

func (c Users) PassResetConfirm(userID uint, token, password string) revel.Result {
	err := ds.PassResetConfirm(userID, token, password)
	if err != nil {
		return c.Error(err)
	}
	return c.OK()
}

func PasswordResetURL(userID uint, token string) (uri string) {
	var host string
	if revel.Server.Addr[0] == ':' {
		host = "127.0.0.1" + revel.Server.Addr
	} else {
		host = revel.Server.Addr
	}
	u := url.URL{}
	u.Scheme = "http"
	u.Host = host
	u.Path = "/page/passreset"
	q := u.Query()
	q.Set("id", strconv.Itoa(int(userID)))
	q.Set("token", token)
	u.RawQuery = q.Encode()
	return u.String()
}

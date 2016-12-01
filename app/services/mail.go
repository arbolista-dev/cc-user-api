// PASWD = mailrobot
// ___________________________________________________________
// |  SPARKPOST KEY: 672c40cdb9bb75b6ccc81a9a080624877b516ca3 |
// ------------------------------------------------------------

package services

import (
	"errors"
	"io/ioutil"
	"log"
	"encoding/json"
	"bytes"
	"net/http"
	"strings"
	"github.com/revel/revel"

)

var client = &http.Client{}

var apiKey string
var sendMail string
var templateId map[string]string

type recipientAPI struct {
	Address string `json:"email"`
	Name 	string `json:"name"`
}


type substitutionAPI struct {
	Name 	string 	`json:"name"`
}

type personalizationAPI struct {
	Recipients       []recipientAPI    	`json:"to"`
	Substitution     map[string]string  `json:"substitutions"`

}


type transmissionAPI struct {
	TemplateId 		 string				`json:"template_id"`
	Senders       	 recipientAPI    	`json:"from"`
	Personalization  []personalizationAPI `json:"personalizations"`
}

func getConfig(key string) (result string) {
	err := false
	result, err = revel.Config.String(key)
	if !err {
		log.Print(key, " not configured");
	}
	return
}

func readConfig() {
	templateId = make(map[string]string)
	apiKey = getConfig("sendgrid.apikey")
	templateId["confirm"] = getConfig("sendgrid.template.confirm")
	templateId["reset"] = getConfig("sendgrid.template.reset")
	sendMail= getConfig("sendgrid.mail")
	// 	apiKey = "672c40cdb9bb75b6ccc81a9a080624877b516ca3"
}

func templateMail(template string, address string,  data map[string]string) (result []byte, err error) {
	recipients := []recipientAPI{
		recipientAPI{
			Address: address,
		},
	}
	sender := recipientAPI{
		Address: sendMail,
	}

 	personalization :=  []personalizationAPI{
 		personalizationAPI{
 			Recipients: recipients,
 			Substitution: data,
 		},
 	}

	request := transmissionAPI{
		TemplateId: templateId[template],
		Personalization: personalization,
		Senders: sender,
	}
	result, err = json.Marshal(request)
	return
}

func SendMail(template string, address string,  data map[string]string) (err error) {
	readConfig()
	jsonValue, _ := templateMail(template, address, data)
	req, err := http.NewRequest("POST", "https://api.sendgrid.com/v3/mail/send", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Print("Send email Error: ", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+ apiKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Print("Send email Error: ", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if strings.TrimSpace(string(body)) != "" {
		err = errors.New(`{"email": "failed"}`)
	}
	return
}

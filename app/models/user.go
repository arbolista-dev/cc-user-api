package models

import (
	"bytes"
	"encoding/gob"
	"github.com/jmoiron/sqlx/types"
	"github.com/revel/revel"
	"time"
  "strconv"
)

type User struct {
	UserID          uint           `json:"user_id" db:"user_id,omitempty"`
	FirstName       string         `json:"first_name" db:"first_name"`
	LastName        string         `json:"last_name" db:"last_name"`
	FacebookID      string         `json:"facebook_id" db:"facebook_id"`
	Password        string         `json:"password" db:"-"`
	Hash            []byte         `json:"-" db:"hash"`
	Salt            []byte         `json:"-" db:"salt"`
	Email           string         `json:"email" db:"email"`
	ValidJTIs       []string       `json:"-" db:"-"`
	ValidJTI        []byte         `json:"-" db:"valid_jti"`
	Answers         types.JSONText `json:"answers" db:"answers"`
	Public					bool				 	 `json:"public" db:"public"`
	City						string 				 `json:"city" db:"city"`
	State						string 				 `json:"state" db:"state"`
	County					string 				 `json:"county" db:"county"`
	Country					string 				 `json:"country" db:"country"`
	HouseholdSize		int 				 	 `json:"household_size" db:"household_size"`
	TotalFootprint	types.JSONText `json:"total_footprint" db:"total_footprint"`
	PhotoUrl				string				 `json:"photo_url" db:"photo_url"`
	ProfileData			types.JSONText `json:"profile_data" db:"profile_data"`
	ResetHash       []byte         `json:"-" db:"reset_hash"`
	ResetExpiration time.Time      `json:"-" db:"reset_expiration"`
	// EmailHash       []byte         `json:"-" db:"email_hash"`
	// EmailExpiration	time.Time      `json:"-" db:"email_expiration"`
}

type UserUpdate struct {
	FirstName       string         `json:"first_name" db:"first_name"`
	LastName        string         `json:"last_name" db:"last_name"`
	Email           string         `json:"email" db:"email"`
	Public					string				 `json:"public"`
	TotalFootprint	types.JSONText `json:"total_footprint" db:"total_footprint"`
	ProfileData			types.JSONText `json:"profile_data" db:"profile_data"`
}

type UserLogin struct {
	Email    		string	`json:"email"`
	Password		string	`json:"password"`
}

type UserFacebook struct {
	FacebookID			string	`json:"facebookID"`
	FacebookToken		string	`json:"facebookToken"`
	Answers       	types.JSONText `json:"answers" db:"answers"`
}
type FacebookToken struct {
	FacebookID		string	`json:"id"`
	FirstName   	string	`json:"first_name"`
	LastName   		string	`json:"last_name"`
	Email   			string	`json:"email"`
	Error       	types.JSONText	`json:"error"`
}

type Answers struct {
	Answers		types.JSONText	`json:"answers"`
}

type AnswersUpdate struct {
	Answers		types.JSONText	`json:"answers"`
	HouseholdSize		int 				 	 `json:"household_size" db:"household_size"`
	TotalFootprint	types.JSONText `json:"total_footprint" db:"total_footprint"`
}

type Location struct {
	City			string	`json:"city" db:"city"`
	State			string 	`json:"state" db:"state"`
	County		string  `json:"county" db:"county"`
	Country		string 	`json:"country" db:"country"`
}

type TotalFootprint struct {
	TotalFootprint		types.JSONText	`json:"total_footprint"`
}

type Email struct {
	FirstName       string         `json:"first_name" db:"first_name"`
	Email		string	`json:"email"`
}

type PaginatedLeaders struct {
	TotalCount		uint64	`json:"total_count"`
	List					[]Leader	`json:"list"`
}

type ProfileData struct {
  Facebook  string `json:"facebook"`
  Twitter   string `json:"twitter"`
  Instagram string `json:"instagram"`
	LinkedIn  string `json:"linkedin"`
	Medium    string `json:"medium"`
	Intro     string `json:"intro"`
}

type Leader struct {
	UserID          	uint           `json:"user_id" db:"user_id"`
	FirstName       	string         `json:"first_name" db:"first_name"`
	LastName        	string         `json:"last_name" db:"last_name"`
	City							string 				 `json:"city" db:"city"`
	State							string 				 `json:"state" db:"state"`
	County						string 				 `json:"county" db:"county"`
	HouseholdSize			int 				 	 `json:"household_size" db:"household_size"`
	TotalFootprint		types.JSONText `json:"total_footprint" db:"total_footprint"`
	PhotoUrl					string				 `json:"photo_url" db:"photo_url"`
	ProfileData				types.JSONText `json:"profile_data" db:"profile_data"`
	Public						bool				 	 `json:"public" db:"public"`
}

// type NeedActivate struct {
// 	Need 			 bool			 `json:"need"`
// 	Name 			 string		 `json:"-"`
// 	Email 		 string		 `json:"-"`
// }

type PasswordReset struct {
	Id 					uint			`json:"id"`
	Token 			string	  `json:"token"`
	Password 		string	  `json:"password"`
}

func (user *User) Validate(v *revel.Validation) {
	v.Required(user.FirstName)
	v.MinSize(user.FirstName, 4)
	v.Required(user.LastName)
	v.MinSize(user.LastName, 4)
	v.Required(user.Password)
	v.MinSize(user.Password, 4)
	v.Required(user.Email)
	v.Email(user.Email)
}

// func (c *User) MarshalJSON() ([]byte, error) {
//     type Alias User
//     return json.Marshal(&struct {
//         *Alias
// 		CreatedAt    string `json:"created_at"`
// 		UpdatedAt    string `json:"updated_at"`
//     }{
//         Alias: (*Alias)(c),
// 		CreatedAt: c.CreatedAt.Format("2006-01-02 15:04:05"),
// 		UpdatedAt: c.UpdatedAt.Format("2006-01-02 15:04:05"),
//     })
// }

func (u User) ContainsJTI(jti string) bool {
	for _, i := range u.ValidJTIs {
		if i == jti {
			return true
		}
	}
	return false
}

func (u *User) AddJTI(jti string) {
	if len(u.ValidJTIs) > 4 {
		u.ValidJTIs = append(u.ValidJTIs[1:], jti)
	} else {
		u.ValidJTIs = append(u.ValidJTIs, jti)
	}
}

func (u *User) ClearAllJTI() {
	u.ValidJTIs = []string{}
}

func (u *User) RemoveJTI(jti string) {
	for j, i := range u.ValidJTIs {
		if i == jti {
			u.ValidJTIs = append(u.ValidJTIs[:j], u.ValidJTIs[j+1:]...)
			break
		}
	}
}

func (u *User) MarshalDB() {
	buffer := &bytes.Buffer{}
	gob.NewEncoder(buffer).Encode(u.ValidJTIs)
	u.ValidJTI = buffer.Bytes()
	if u.Answers == nil {
		u.Answers = types.JSONText("{}")
	}
	if u.TotalFootprint == nil {
		u.TotalFootprint = types.JSONText("{}")
	}
	if u.ProfileData == nil {
		u.ProfileData = types.JSONText("{}")
	}
}

func (u *User) UnmarshalDB() {
	buffer := bytes.NewReader(u.ValidJTI)
	s := []string{}
	gob.NewDecoder(buffer).Decode(&s)
	u.ValidJTIs = s
}

func (u *User) Update(n UserUpdate) {
	if n.FirstName != "" {
		u.FirstName = n.FirstName
	}
	if n.LastName != "" {
		u.LastName = n.LastName
	}
	if n.Email != "" {
		u.Email = n.Email
	}
  // updating Public with string value because when using bool a non-given public parameter sets it to false
	if n.Public != "" {
    b, err := strconv.ParseBool(n.Public)
    if err != nil {
      return
    }
		u.Public = b
	}
  if n.TotalFootprint != nil {
    u.TotalFootprint = n.TotalFootprint
  }
	if n.ProfileData != nil {
		u.ProfileData = n.ProfileData
	}
}

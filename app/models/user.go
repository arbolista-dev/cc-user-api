package models

import (
	"github.com/revel/revel"
	"encoding/gob"
	"bytes"
	"github.com/jmoiron/sqlx/types"
)

type User struct {
	UserID         uint           `json:"user_id" db:"user_id,omitempty"`
	FirstName      string         `json:"first_name" db:"fist_name"`
	LastName	   string         `json:"last_name" db:"last_name"`
	Password       string         `json:"password" db:"-"`
	Hash           []byte         `json:"-" db:"hash"`
	Salt           []byte         `json:"-" db:"salt"`
	Email          string         `json:"email" db:"email"`
	ValidJTIs      []string       `json:"-" db:"-"`
	ValidJTI       []byte         `json:"-" db:"valid_jti"`
	Answers        types.JSONText `json:"answers" db:"answers"`
}

type UserLogin struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
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
}

func (u *User) UnmarshalDB() {
	buffer := &bytes.Buffer{}
	gob.NewDecoder(buffer).Decode(u.ValidJTIs)
}
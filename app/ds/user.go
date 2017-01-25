package ds

import (
	"crypto/rand"
	"errors"
	"github.com/arbolista-dev/cc-user-api/app/models"
	"github.com/arbolista-dev/cc-user-api/app/utils"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
	"upper.io/db.v2"
	"github.com/lib/pq"
	"net/http"
	"encoding/json"
	"io/ioutil"
)

func GetSession(token string) (userID uint, jti string, err error) {
	claims, err := ValidateToken(token)
	if err != nil {
		return
	}

	userID = uint(claims["id"].(float64))
	var user models.User
	err = userSource.Find(db.Cond{"user_id": userID}).One(&user)
	if err != nil {
		return
	}

	user.UnmarshalDB()

	if time.Now().Unix() > int64(claims["exp"].(float64)) {
		err = errors.New(`{"session": "expired"}`)
		user.RemoveJTI(claims["jti"].(string))
		user.MarshalDB()
		_ = userSource.Find(db.Cond{"user_id": userID}).Update(user)
		return
	}

	if !user.ContainsJTI(claims["jti"].(string)) {
		err = errors.New(`{"session": "non-existent"}`)
		return
	}

	jti = claims["jti"].(string)
	return
}

func Add(user models.User) (login map[string]interface{}, err error) {
	hashPassword(&user)

	user.MarshalDB()
	temp, err := userSource.Insert(user)
	if err != nil {
		pqErr := err.(*pq.Error)
		if pqErr != nil {
			if pqErr.Code == "23505" {
				err = errors.New(`{"email": "non-unique"}`)
				return
			}
		}
		return
	}
	userID := uint(temp.(int64))

	token, err := newToken(userID)
	if err != nil {
		return
	}

	sToken, err := token.SignedString([]byte(os.Getenv("CC_JWTSIGN")))
	if err != nil {
		return
	}
	user.AddJTI(token.Claims["jti"].(string))

	user.MarshalDB()
	err = userSource.Find("user_id", userID).Update(user)
	if err != nil {
		return
	}

	login = map[string]interface{}{
		"user_id": userID,
		"name":    user.FirstName,
		"token":   sToken,
		"answers": user.Answers.String(),
	}
	return
}

func Delete(userID uint) (err error) {
  err = DeleteUserActions(userID)
	err = userSource.Find(db.Cond{"user_id": userID}).Delete()
	return
}

func ValidateFacebookToken(logRequest models.UserFacebook) (facebookData models.FacebookToken,err error) {
	resp, err := http.Get("https://graph.facebook.com/v2.5/"+logRequest.FacebookID+"?fields=id,first_name,last_name,email&access_token="+logRequest.FacebookToken);
	if err != nil {
  	return
  }

  body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &facebookData)
  if facebookData.Error !=nil || facebookData.FacebookID != logRequest.FacebookID {
  	err =  errors.New("Token invalid");
  	return
  }
  return
}

func LoginFacebook(logRequest models.UserFacebook) (login map[string]interface{}, err error) {
	var user models.User
	var facebookData models.FacebookToken
	facebookData, err =ValidateFacebookToken(logRequest)
	if err != nil {
		err = errors.New(`{"facebook":"token-invalid"}`);
	}
	if err!=nil {
		return
	}
	err = userSource.Find("facebook_id", logRequest.FacebookID).One(&user)
	if err != nil {
		err = nil;
		user.FacebookID =  facebookData.FacebookID;
		user.Email =  facebookData.Email;
		user.FirstName = facebookData.FirstName;
		user.LastName = facebookData.LastName;
		user.Answers = logRequest.Answers;
		login, err = Add(user);
		return
	}

	token, err := newToken(user.UserID)
	if err != nil {
		return
	}

	sToken, err := token.SignedString([]byte(os.Getenv("CC_JWTSIGN")))
	if err != nil {
		return
	}

	user.AddJTI(token.Claims["jti"].(string))

	user.MarshalDB()
	err = userSource.Find("user_id", user.UserID).Update(user)
	if err != nil {
		return
	}

	login = map[string]interface{}{
		"user_id": user.UserID,
		"name":    user.FirstName,
		"token":   sToken,
		"answers": user.Answers.String(),
	}
	return
}

func Login(logRequest models.UserLogin) (login map[string]interface{}, err error) {
	var user models.User
	err = userSource.Find("email", logRequest.Email).One(&user)
	if err != nil {
		err = errors.New(`{"email": "non-existent"}`)
		return
	}
	user.UnmarshalDB()

	err = bcrypt.CompareHashAndPassword(user.Hash, append([]byte(logRequest.Password), user.Salt...))
	if err != nil {
		err = errors.New(`{"password": "incorrect"}`)
		return
	}
	token, err := newToken(user.UserID)
	if err != nil {
		return
	}

	sToken, err := token.SignedString([]byte(os.Getenv("CC_JWTSIGN")))
	if err != nil {
		return
	}
	user.AddJTI(token.Claims["jti"].(string))

	user.MarshalDB()
	err = userSource.Find("user_id", user.UserID).Update(user)
	if err != nil {
		return
	}

	login = map[string]interface{}{
		"user_id": user.UserID,
		"name":    user.FirstName,
		"token":   sToken,
		"answers": user.Answers.String(),
	}
	return
}

func Logout(userID uint, jti string) (err error) {
	var user models.User
	err = userSource.Find(db.Cond{"user_id": userID}).One(&user)
	if err != nil {
		return
	}
	user.UnmarshalDB()

	user.RemoveJTI(jti)

	user.MarshalDB()
	err = userSource.Find(db.Cond{"user_id": userID}).Update(user)
	return
}

func LogoutAll(userID uint) (err error) {
	var user models.User
	err = userSource.Find(db.Cond{"user_id": userID}).One(&user)
	if err != nil {
		return
	}
	user.UnmarshalDB()

	user.ClearAllJTI()

	user.MarshalDB()
	err = userSource.Find(db.Cond{"user_id": userID}).Update(user)
	return
}

func Show(userID uint, auth bool) (profile map[string]interface{}, err error) {
	var user models.Leader
	err = userSource.Find(db.Cond{"user_id": userID}).One(&user)
	if err != nil {
		err = errors.New(`{"profile": "non-existent"}`)
		return
	}

	if user.Public == false && auth == false {
		err = errors.New(`{"profile": "not-public"}`)
		return
	}

  actions, err := RetrieveUserActions(userID)
  if err != nil {
		return
	}

	var profileData models.ProfileData
	err = json.Unmarshal(user.ProfileData, &profileData)
	if err != nil {
		return
	}

	profile = map[string]interface{}{
		"user_id":    			user.UserID,
		"first_name": 			user.FirstName,
		"last_name":  			user.LastName,
		"city":  						user.City,
		"state":  					user.State,
		"county":  					user.County,
		"household_size": 	user.HouseholdSize,
		"total_footprint": 	user.TotalFootprint.String(),
		"photo_url": 				user.PhotoUrl,
		"profile_data":			profileData,
		"public":						user.Public,
    "actions":          actions.List,
	}
	return
}

func SetLocation(userID uint, location models.Location) (err error) {
	var user models.User
	err = userSource.Find(db.Cond{"user_id": userID}).One(&user)
	if err != nil {
		return
	}
	user.City = location.City
	user.State = location.State
	user.County = location.County
	err = userSource.Find(db.Cond{"user_id": userID}).Update(user)
	return
}

func SetPhoto(userID uint, photo_url string) (photo_set map[string]interface{}, err error) {
	var user models.User
	err = userSource.Find(db.Cond{"user_id": userID}).One(&user)
	if err != nil {
		return
	}
	user.PhotoUrl = photo_url
	err = userSource.Find(db.Cond{"user_id": userID}).Update(user)
	photo_set = map[string]interface{}{
		"photo_url": photo_url,
	}
	return
}

func Update(userID uint, userNew models.UserUpdate) (err error) {
	var user models.User
	err = userSource.Find(db.Cond{"user_id": userID}).One(&user)
	if err != nil {
		return
	}
	user.Update(userNew)
	err = userSource.Find(db.Cond{"user_id": userID}).Update(user)
	return
}

func UpdateAnswers(userID uint, userAnswers models.AnswersUpdate) (err error) {
	var user models.User
	err = userSource.Find(db.Cond{"user_id": userID}).One(&user)
	if err != nil {
		return
	}
	user.Answers = userAnswers.Answers
	user.HouseholdSize = userAnswers.HouseholdSize
	user.TotalFootprint = userAnswers.TotalFootprint
	err = userSource.Find(db.Cond{"user_id": userID}).Update(user)
	return
}

func PassResetRequest(email string) (userID uint, token string, err error) {
	var user models.User
	err = userSource.Find(db.Cond{"email": email}).One(&user)
	if err != nil {
		return
	}
	token = hashReset(&user)
	userID = user.UserID
	err = userSource.Find(db.Cond{"user_id": user.UserID}).Update(user)
	return
}

func PassResetConfirm(userID uint, token, password string) (err error) {
	var user models.User
	err = userSource.Find(db.Cond{"user_id": userID}).One(&user)
	if err != nil {
		return
	}
	user.UnmarshalDB()
	if user.ResetExpiration.After(time.Now()) {
		user.ResetHash = []byte{}
		user.ResetExpiration = time.Time{}
		err = errors.New(`{"password-reset": "expired"}`)
		return
	}

	err = bcrypt.CompareHashAndPassword(user.ResetHash, []byte(token))
	if err != nil {
		err = errors.New(`{"reset-token": "corrupt"}`)
		return err
	}

	user.Password = password
	hashPassword(&user)

	user.ClearAllJTI()
	user.ResetHash = []byte{}
	user.ResetExpiration = time.Time{}

	user.MarshalDB()
	err = userSource.Find(db.Cond{"user_id": user.UserID}).Update(user)
	if err != nil {
		return
	}
	return
}

func ListLeaders(limit int, offset int, state string, household_size int) (leaders models.PaginatedLeaders, err error) {

	if household_size != -1 {
		if len(state) == 0 {
			query = leadersSource.Find(db.Cond{"household_size": household_size})
		} else if len(state) > 0 {
			query = leadersSource.Find(db.Cond{"household_size": household_size}, db.Cond{"state": state})
		}
	} else {
		if len(state) == 0 {
			query = leadersSource.Find()
		} else if len(state) > 0 {
			query = leadersSource.Find(db.Cond{"state": state})
		}
	}

	count, err := query.Count()
	if err != nil {
		return
	}
	leaders.TotalCount = count

	list := query.Limit(limit).Offset(offset)
	err = list.All(&leaders.List)
	if err != nil {
		return
	}

	return
}

func ListLocations() (locations []models.Location, err error) {

	q := leadersSource.Find().Select("city", "state", "county").Group("city", "state", "county")
	err = q.All(&locations)
	if err != nil {
		return
	}
	return
}

func hashPassword(user *models.User) {
	b := make([]byte, 10)
	_, err := rand.Read(b)
	user.Hash, err = bcrypt.GenerateFromPassword(append([]byte(user.Password), b...), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	user.Salt = b
}

func hashReset(user *models.User) (token string) {
	token = utils.RandString(10)
	var err error
	user.ResetHash, err = bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	user.ResetExpiration = time.Now().Add(time.Minute * 5)
	return
}

package ds

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/arbolista-dev/cc-user-api/app/utils"
	"os"
	"time"
)


func newToken(userID uint) (token *jwt.Token, err error) {
	token = jwt.New(jwt.SigningMethodHS256)
	token.Header["alg"] = "HS256"
	token.Header["typ"] = "JWT"
	token.Claims["id"] = userID
	token.Claims["iat"] = time.Now().Unix()
	token.Claims["exp"] = time.Now().Add(time.Second * 3600 * 24 * 14).Unix()
	token.Claims["jti"] = utils.RandString(5)
	return
}

func ValidateToken(sToken string) (claims map[string]interface{}, err error) {
	token, err := jwt.Parse(sToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if token.Header["alg"] != "HS256" {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// return []byte("pl8IKa8Wz5tu64JuV3ksSQ7YVyDDjet17jE5YXS37lIasCxjhYlHjYYGnNT9Gzs"), nil
		return []byte(os.Getenv("CC_JWTSIGN")), nil
	})

	if err == nil && token.Valid {
		claims = token.Claims
	}
	return
}

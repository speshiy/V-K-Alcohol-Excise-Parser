package common

import (
	"fmt"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

//TokenCustomClaims data from client in token
type TokenCustomClaims struct {
	jwt.StandardClaims
}

var mWO = "Hello user! You are an successfully authorized!"

//JWTCreate  - create a new token
func JWTCreate(userID uint, accessLevel string, hours time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &jwt.StandardClaims{
		Id:        fmt.Sprint(userID),
		Subject:   accessLevel,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute * hours).Unix(),
	})

	tokenstring, err := token.SignedString([]byte(mWO))
	if err != nil {
		return "", err
	}
	return tokenstring, nil
}

//JWTParse Parse token
func JWTParse(tokenstring string) (jwt.StandardClaims, error) {
	tokenCustomClaims := jwt.StandardClaims{}
	_, err := jwt.ParseWithClaims(tokenstring, &tokenCustomClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(mWO), nil
	})

	if err != nil {
		return tokenCustomClaims, err
	}

	return tokenCustomClaims, nil
}

//JWTParseStringID Parse token
func JWTParseStringID(tokenstring string) (string, error) {
	tokenCustomClaims, err := JWTParse(tokenstring)
	if err != nil {
		return "", err
	}

	return tokenCustomClaims.Id, nil
}

//JWTParseUintID Parse token
func JWTParseUintID(tokenstring string) (uint, error) {
	var userID uint64

	tokenCustomClaims, err := JWTParse(tokenstring)
	if err != nil {
		return 0, err
	}

	userID, _ = strconv.ParseUint(tokenCustomClaims.Id, 10, 64)

	return uint(userID), nil
}

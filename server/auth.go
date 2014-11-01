package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"code.google.com/p/go.crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
)

type Auth struct {
	privateKey []byte
	publicKey  []byte
}

func NewAuth() *Auth {
	// todo: configuration of paths (private/public key)
	// todo: error handling (when keys not present)
	privKey, err := ioutil.ReadFile("develop/jwt-keys/dev.rsa")
	if err != nil {
		log.Fatal(err)
	}

	pubKey, err := ioutil.ReadFile("develop/jwt-keys/dev.rsa.pub")
	if err != nil {
		log.Fatal(err)
	}

	return &Auth{privKey, pubKey}
}

func (this *Auth) Authenticate(user *User, pwd string) (string, error) {
	err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(pwd))
	if err != nil {
		return "", errors.New("invalid password")
	}

	token, err := this.createToken(user)
	if err != nil {
		return "", errors.New("token could not created")
	}

	return token, nil
}

func (this *Auth) createToken(user *User) (string, error) {
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims["user"] = user
	// todo: exp -> config
	token.Claims["exp"] = time.Now().Unix() + 3600
	return token.SignedString(this.privateKey)
}

func (this *Auth) getValidToken(r *http.Request) *jwt.Token {
	token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface{}, error) {
		return this.privateKey, nil
	})
	if err == nil && token.Valid {
		return token
	}
	return nil
}

func (this *Auth) ValidateToken(r *http.Request) bool {
	token := this.getValidToken(r)
	if token == nil {
		return false
	}
	return true
}

func (this *Auth) AuthenticateRequest(r *http.Request) (bool, *User) {
	token := this.getValidToken(r)
	if token != nil {
		user, err := this.extractUser(token.Raw)
		if err != nil {
			return false, nil
		}
		return true, user
	}
	return false, nil
}

func (this *Auth) extractUser(token string) (*User, error) {
	// todo: rethink: how to get user struct out of token?
	// this seems lika a ridiculous solution:
	// 1. get part of token, which contains claims (middle part):
	//    eyJhb[...]XVCJ9.eyJle[...]IifX0.4ZEVb[...]xc8ig -> eyJle[...]IifX0
	// 2. decote that segment using jwt -> []byte
	// 3. []byte to string:
	//    {"exp":1414463656,"user":{"id":1,"username":"admin","first_name":"","last_name":"","email":""}}
	// 4. extract user json part (using regexp):
	//    {"id":1,"username":"admin","first_name":"","last_name":"","email":""}
	// 5. unmarshal

	tokenParts := strings.Split(token, ".")
	if len(tokenParts) != 3 {
		return nil, errors.New("extract user from token: invalid raw token")
	}

	tokenClaimsPart := tokenParts[1]
	claimsData, err := jwt.DecodeSegment(tokenClaimsPart)
	if err != nil {
		return nil, err
	}

	jsonClaims := string(claimsData)
	re := regexp.MustCompile(`"user":(.*)?}`)
	matches := re.FindStringSubmatch(jsonClaims)
	if len(matches) != 2 {
		return nil, errors.New("extract user from token: invalid json data")
	}

	jsonUser := matches[1]
	var user *User
	err = json.Unmarshal([]byte(jsonUser), &user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (this *Auth) GeneratePasswordHash(pwd string) ([]byte, error) {
	const bcryptCost int = 13
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcryptCost)
	if err != nil {
		return nil, err
	}
	return pwdHash, nil
}

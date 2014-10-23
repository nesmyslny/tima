package services

import (
	"code.google.com/p/go.crypto/bcrypt"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"gnomon/dbaccess"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type AuthService struct {
	db         *DbAccess.Db
	privateKey []byte
	publicKey  []byte
}

func NewAuthService(db *DbAccess.Db) *AuthService {
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

	return &AuthService{db, privKey, pubKey}
}

func (this *AuthService) Authenticate(username string, pwd string) (string, error) {
	user := this.db.GetUserByName(username)
	if user == nil {
		return "", errors.New("invalid username")
	}

	err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(pwd))
	if err != nil {
		return "", errors.New("invalid password")
	}

	token, err := this.createToken(username)
	if err != nil {
		return "", errors.New("token could not created")
	}

	return token, nil
}

func (this *AuthService) createToken(username string) (string, error) {
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims["username"] = username
	// todo: exp -> config
	token.Claims["exp"] = time.Now().Unix() + 3600
	return token.SignedString(this.privateKey)
}

func (this *AuthService) ValidateToken(r *http.Request) bool {
	token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface{}, error) {
		return this.privateKey, nil
	})

	if err == nil && token.Valid {
		return true
	}
	return false
}

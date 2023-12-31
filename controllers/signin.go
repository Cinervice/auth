
package authservice

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/shadowshot-x/micro-product-go/authservice/data"
	"github.com/shadowshot-x/micro-product-go/authservice/jwt"
)

func getSignedToken() (string, error) {

	claimsMap := map[string]string{
		"aud": "frontend.knowsearch.ml",
		"iss": "knowsearch.ml",
		"exp": fmt.Sprint(time.Now().Add(time.Minute * 1).Unix()),
	}

	secret := "S0m3_R4n90m_sss"
	header := "HS256"
	tokenString, err := jwt.GenerateToken(header, claimsMap, secret)
	if err != nil {
		return tokenString, err
	}
	return tokenString, nil
}

func validateUser(email string, passwordHash string) (bool, error) {
	usr, exists := data.GetUserObject(email)
	if !exists {
		return false, errors.New("user does not exist")
	}
	passwordCheck := usr.ValidatePasswordHash(passwordHash)

	if !passwordCheck {
		return false, nil
	}
	return true, nil
}

func SigninHandler(rw http.ResponseWriter, r *http.Request) {

	if _, ok := r.Header["Email"]; !ok {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Email Missing"))
		return
	}
	if _, ok := r.Header["Passwordhash"]; !ok {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Passwordhash Missing"))
		return
	}
	valid, err := validateUser(r.Header["Email"][0], r.Header["Passwordhash"][0])
	if err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("User Does not Exist"))
		return
	}

	if !valid {
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Incorrect Password"))
		return
	}
	tokenString, err := getSignedToken()
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Internal Server Error"))
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(tokenString))
}

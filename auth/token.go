package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func CreateToken(user_id string, email string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

func TokenValid(r *http.Request) error {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, parseToken)

	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		Pretty(claims)
	}

	return nil
}

func ExtractToken(r *http.Request) string {
	keys := r.URL.Query()

	token := keys.Get("token")

	if token != "" {
		return token
	}

	bearerToken := r.Header.Get("Authorization")

	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}

	return ""
}

func ExtractTokenIDAndEmail(r *http.Request) (string, string, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, parseToken)

	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {

		id, okI := claims["user_id"].(string)
		email, okE := claims["email"].(string)

		if (len(id) == 0 || !okI) && (len(email) == 0 || !okE) {
			return "", "", errors.New("No User Id in claims")
		}

		return id, email, nil
	}

	return "", "", nil
}

func Pretty(data interface{}) string {
	b, err := json.MarshalIndent(data, "", " ")

	if err != nil {
		log.Println(err)
		return ""
	}

	return string(b)
}

func parseToken(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}
	return []byte(os.Getenv("API_SECRET")), nil
}

package auth

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/ash3798/Social-Network/config"
	jwt "github.com/dgrijalva/jwt-go"
)

//LoginResponse is response user will get on successful login
type LoginResponse struct {
	Token string `json:"token"`
}

//CreateToken will create jwt token with username in it
func CreateToken(userName string) ([]byte, error) {
	atClaims := jwt.MapClaims{}
	atClaims["username"] = userName
	atClaims["exp"] = time.Now().Add(time.Second * time.Duration(config.Manager.TokenExpireTimeSec)).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(config.Manager.AccessSecret))
	if err != nil {
		return []byte(""), err
	}

	response := LoginResponse{}
	response.Token = token

	res, err := json.Marshal(response)
	if err != nil {
		log.Println("error while marshelling token response. Error :", err)
	}
	return res, nil
}

//extractToken will extract token string from Authorization header
func extractToken(authToken string) (string, error) {
	extractedToken := strings.Split(authToken, "Bearer ")
	if len(extractedToken) == 2 {
		return extractedToken[1], nil
	}

	return "", errors.New("incorrect format of Authorization Token")
}

//ValidateToken will check if authorization token is valid or not
func ValidateToken(authToken string) (jwt.MapClaims, error) {
	signedToken, err := extractToken(authToken)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Manager.AccessSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return nil, err
	}

	isValid := claims.VerifyExpiresAt(time.Now().Unix(), true)

	if !isValid {
		return nil, errors.New("JWT is expired")
	}

	return claims, nil
}

package auth

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ash3798/Social-Network/config"
	"github.com/stretchr/testify/assert"
)

func TestCreateToken(t *testing.T) {
	username := "ashish"
	config.InitEnv()

	response, err := CreateToken(username)
	if err != nil {
		t.Error("error while generating token")
	}
	if len(response) <= 0 {
		t.Errorf("Token didn't get generated")
	}

	t.Log("token generated : ", string(response))
}

func TestVerifyToken(t *testing.T) {
	username := "ashish"
	config.InitEnv()

	config.Manager.TokenExpireTimeSec = 3

	res, err := CreateToken(username)
	if err != nil {
		t.Error("error while generating token")
	}

	loginRes := LoginResponse{}
	err = json.Unmarshal(res, &loginRes)
	if err != nil {
		t.Error("error while unmarshelling response of create Token, Error : ", err.Error())
	}

	authToken := "Bearer " + loginRes.Token
	claims, err := ValidateToken(authToken)
	if err != nil {
		t.Error("correct token is not getting validated, Error: ", err.Error())
	}
	tknUsername := claims["username"].(string)
	assert.Equal(t, username, tknUsername) //check if token claim has provided username in it or not

	//test for validating expire token
	time.Sleep(5 * time.Second)

	_, err = ValidateToken(authToken)
	if err == nil {
		t.Error("error is expected for expired token. Validation should fail")
	}
}

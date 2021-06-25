package auth

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ash3798/Social-Network/config"
)

func TestCreateToken(t *testing.T) {
	username := "ashish"
	config.InitEnv()

	token, err := CreateToken(username)
	if err != nil {
		t.Error("error while generating token")
	}
	if len(token) <= 0 {
		t.Errorf("Token didn't get generated")
	}

	t.Log("token generated : ", string(token))
}

func TestVerifyToken(t *testing.T) {
	username := "ashish"
	config.InitEnv()

	config.Manager.TokenExpireTimeSec = 3

	res, err := CreateToken(username)
	if err != nil {
		t.Error("error while generating token")
	}

	loginRes := loginResponse{}
	err = json.Unmarshal(res, &loginRes)
	if err != nil {
		t.Error("error while unmarshelling response of create Token, Error : ", err.Error())
	}

	authToken := "Bearer " + loginRes.Token
	err = ValidateToken(authToken)
	if err != nil {
		t.Error("correct token is not getting validated, Error: ", err.Error())
	}

	//test for validating expire token
	time.Sleep(5 * time.Second)

	err = ValidateToken(authToken)
	if err == nil {
		t.Error("error is expected for expired token. Validation should fail")
	}
}

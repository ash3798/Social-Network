package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ash3798/Social-Network/auth"
	"github.com/ash3798/Social-Network/config"
	"github.com/stretchr/testify/assert"
)

func TestIsAuthorized(t *testing.T) {
	username := "ashish"
	config.InitEnv()

	res, err := auth.CreateToken(username)
	if err != nil {
		t.Error("error while generating token")
	}

	loginRes := auth.LoginResponse{}
	err = json.Unmarshal(res, &loginRes)
	if err != nil {
		t.Error("error while unmarshelling response of create Token, Error : ", err.Error())
	}

	token := loginRes.Token

	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/comment", bytes.NewBuffer([]byte("")))
	if err != nil {
		t.Error("error while creating request for test , Error : ", err.Error())
	}

	//adding token to header
	req.Header.Add("Authorization", "Bearer "+token)

	tknUsername, ok := isAuthorized(recorder, req)
	assert.Equal(t, true, ok, "token should get authorized")

	assert.Equal(t, username, tknUsername, "it should return the same username with which token was created")

}

func TestIsAuthorizedWithInvalidToken(t *testing.T) {
	config.InitEnv()

	token := "invalid token"

	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/comment", bytes.NewBuffer([]byte("")))
	if err != nil {
		t.Error("error while creating request for test , Error : ", err.Error())
	}

	//adding token to header
	req.Header.Add("Authorization", "Bearer "+token)

	tknUsername, ok := isAuthorized(recorder, req)
	assert.Equal(t, false, ok, "token should not get authorized")
	assert.Equal(t, "", tknUsername)

	assert.Equal(t, 401, recorder.Code, "status should be unauthorized")
}

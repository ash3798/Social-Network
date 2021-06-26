//containes test for all app functionalities
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/ash3798/Social-Network/auth"
	"github.com/ash3798/Social-Network/config"
	"github.com/ash3798/Social-Network/database"
	"github.com/ash3798/Social-Network/server"
	"github.com/ash3798/Social-Network/structures"
	"github.com/stretchr/testify/assert"
)

//mockDatabase mocks the responses from database for unit tests
type mockDatabase struct{}

func (m mockDatabase) PrepareDatabase() error                                        { return nil }
func (m mockDatabase) InsertUser(userInfo structures.User) (int, error)              { return 1, nil }
func (m mockDatabase) InsertComment(commentInfo structures.CommentInfo) (int, error) { return 1, nil }
func (m mockDatabase) InsertReaction(reactionInfo structures.ReactionInfo) (int, error) {
	return 1, nil
}
func (m mockDatabase) DeleteComment(comment_id int, username string) error      { return nil }
func (m mockDatabase) ValidateLoginCreds(loginCreds structures.LoginCred) error { return nil }
func (m mockDatabase) GetComments(username string) ([]structures.WallUnit, error) {
	arr := []structures.WallUnit{}
	return arr, nil
}
func (m mockDatabase) GetReactionCount(commentID int) (map[string]int, error) {
	mp := make(map[string]int)
	return mp, nil
}
func (m mockDatabase) CloseDBConnection() {}

func (m mockDatabase) GetCommentByID(commentID int) (string, error) { return "ashish", nil }

//preconfig does initializations and put mocks before tests
func preconfig() {
	config.InitEnv()
	config.InitReactions()
	config.Manager.AuthEnabled = false

	database.Action = mockDatabase{}
}

func generateToken() (string, error) {
	username := "ashish"
	config.InitEnv()

	res, err := auth.CreateToken(username)
	if err != nil {
		return "", errors.New("error while generating token")
	}

	loginRes := auth.LoginResponse{}
	err = json.Unmarshal(res, &loginRes)
	if err != nil {
		return "", errors.New("error while unmarshelling response of create Token, Error : " + err.Error())
	}

	return loginRes.Token, nil
}

func TestCreateUser(t *testing.T) {
	preconfig()
	body := []byte(`{
		"username" : "ash123",
		"name" : "ash",
		"password" : "pass123"
		}`)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/createuser", bytes.NewBuffer(body))

	server.HandleCreateUser(recorder, req)

	assert.Equal(t, 200, recorder.Code)
}

func TestCreateUserWrongMethod(t *testing.T) {
	preconfig()
	body := []byte(`{
		"username" : "ash123",
		"name" : "ash",
		"password" : "pass123"
		}`)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/createuser", bytes.NewBuffer(body))

	server.HandleCreateUser(recorder, req)

	assert.Equal(t, 405, recorder.Code, "status is expected to 405 (MethodNotAllowed)")
}

func TestCreateUserInvalidPayload(t *testing.T) {
	preconfig()
	body := []byte(`{
		"username" : "ash123",
		"name" : "ash",
		`)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/createuser", bytes.NewBuffer(body))

	server.HandleCreateUser(recorder, req)

	assert.Equal(t, 400, recorder.Code)
}

func TestComment(t *testing.T) {
	preconfig()
	token, err := generateToken()
	if err != nil {
		t.Errorf(err.Error())
	}

	body := []byte(`{
		"comment_text" : "retest comment 1",
		"parent_comment_id" : 0,
		"sender_username" : "ash",
		"receiver_username":"nit"
	}`)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/comment", bytes.NewBuffer(body))
	req.Header.Add("Authorization", "Bearer "+token)
	server.HandleComment(recorder, req)

	assert.Equal(t, 200, recorder.Code, recorder.Body)
}

func TestCommentWrongMethod(t *testing.T) {
	preconfig()
	token, err := generateToken()
	if err != nil {
		t.Errorf(err.Error())
	}

	body := []byte(`{
		"comment_text" : "retest comment 1",
		"parent_comment_id" : 0,
		"sender_username" : "ash",
		"receiver_username":"nit"
	}`)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/comment", bytes.NewBuffer(body))
	req.Header.Add("Authorization", "Bearer "+token)
	server.HandleComment(recorder, req)

	assert.Equal(t, 405, recorder.Code, "status is expected to 405 (MethodNotAllowed)")
}

func TestCommentInvalidPayload(t *testing.T) {
	preconfig()
	token, err := generateToken()
	if err != nil {
		t.Errorf(err.Error())
	}

	body := []byte(`{`)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/createuser", bytes.NewBuffer(body))
	req.Header.Add("Authorization", "Bearer "+token)
	server.HandleCreateUser(recorder, req)

	assert.Equal(t, 400, recorder.Code)
}

func TestDeleteComment(t *testing.T) {
	preconfig()
	token, err := generateToken()
	if err != nil {
		t.Errorf(err.Error())
	}

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/comment?id=1", bytes.NewBuffer([]byte("")))
	req.Header.Add("Authorization", "Bearer "+token)
	server.HandleComment(recorder, req)

	assert.Equal(t, 200, recorder.Code)
}

func TestDeleteCommentWithNoQueryParams(t *testing.T) {
	preconfig()
	token, err := generateToken()
	if err != nil {
		t.Errorf(err.Error())
	}

	//without id query param
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/comment", bytes.NewBuffer([]byte("")))
	req.Header.Add("Authorization", "Bearer "+token)
	server.HandleComment(recorder, req)
	assert.Equal(t, 404, recorder.Code)

}

func TestCreateReaction(t *testing.T) {
	preconfig()
	token, err := generateToken()
	if err != nil {
		t.Errorf(err.Error())
	}

	body := []byte(`{
		"comment_id" : 8,
		"reaction" : "dislike"
	}`)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/reaction", bytes.NewBuffer(body))
	req.Header.Add("Authorization", "Bearer "+token)
	server.HandleCreateReaction(recorder, req)

	assert.Equal(t, 200, recorder.Code, recorder.Body)
}

func TestCreateReactionWrongMethod(t *testing.T) {
	preconfig()
	token, err := generateToken()
	if err != nil {
		t.Errorf(err.Error())
	}

	body := []byte(`{
		"comment_id" : 8,
		"reaction" : "dislike"
	}`)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/reaction", bytes.NewBuffer(body))
	req.Header.Add("Authorization", "Bearer "+token)
	server.HandleCreateReaction(recorder, req)
	assert.Equal(t, 405, recorder.Code, recorder.Body)
}

func TestCreateReactionInvalidPayload(t *testing.T) {
	preconfig()
	body := []byte(`{`)
	token, err := generateToken()
	if err != nil {
		t.Errorf(err.Error())
	}

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/reaction", bytes.NewBuffer(body))
	req.Header.Add("Authorization", "Bearer "+token)
	server.HandleCreateReaction(recorder, req)

	assert.Equal(t, 400, recorder.Code)
}

func TestCreateReactionInvalidReaction(t *testing.T) {
	preconfig()
	token, err := generateToken()
	if err != nil {
		t.Errorf(err.Error())
	}

	body := []byte(`{
		"comment_id" : 8,
		"reaction" : "invalid"
	}`)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/reaction", bytes.NewBuffer(body))
	req.Header.Add("Authorization", "Bearer "+token)
	server.HandleCreateReaction(recorder, req)

	assert.Equal(t, 400, recorder.Code, recorder.Body)
}

func TestGetWall(t *testing.T) {
	preconfig()
	token, err := generateToken()
	if err != nil {
		t.Errorf(err.Error())
	}

	body := []byte(``)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/wall", bytes.NewBuffer(body))
	req.Header.Add("Authorization", "Bearer "+token)
	server.HandleGetWall(recorder, req)

	assert.Equal(t, 200, recorder.Code, recorder.Body)
}

func TestGetWallWithWrongMethod(t *testing.T) {
	preconfig()
	token, err := generateToken()
	if err != nil {
		t.Errorf(err.Error())
	}

	body := []byte(``)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/wall", bytes.NewBuffer(body))
	req.Header.Add("Authorization", "Bearer "+token)
	server.HandleGetWall(recorder, req)

	assert.Equal(t, 405, recorder.Code, recorder.Body)
}

func TestCreateSubcomment(t *testing.T) {
	preconfig()
	token, err := generateToken()
	if err != nil {
		t.Errorf(err.Error())
	}

	body := []byte(`{
		"comment_text" : "ash subcommented",
		"parent_comment_id" : 2
	}`)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/subcomment", bytes.NewBuffer(body))
	req.Header.Add("Authorization", "Bearer "+token)

	server.HandleCreateSubcomment(recorder, req)

	assert.Equal(t, 200, recorder.Code, recorder.Body)

}

func TestCreateSubcommentWithInvalidID(t *testing.T) {
	preconfig()
	token, err := generateToken()
	if err != nil {
		t.Errorf(err.Error())
	}

	body := []byte(`{
		"comment_text" : "ash subcommented",
		"parent_comment_id" : 0
	}`)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/subcomment", bytes.NewBuffer(body))
	req.Header.Add("Authorization", "Bearer "+token)

	server.HandleCreateSubcomment(recorder, req)

	assert.Equal(t, 400, recorder.Code, recorder.Body)

}

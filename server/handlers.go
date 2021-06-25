package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/ash3798/Social-Network/auth"
	"github.com/ash3798/Social-Network/config"
	"github.com/ash3798/Social-Network/task"
)

//HandleCreateUser handles the request for creating a new user
func HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Wrong method used. Please use POST method", http.StatusMethodNotAllowed)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read the request body", http.StatusBadRequest)
		return
	}

	//log.Println(string(data))
	id, err := task.Action.CreateUser(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("User Created with ID :" + strconv.Itoa(id)))
}

//HandleLogin handles the login request
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Wrong method used. Please use POST method", http.StatusMethodNotAllowed)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read the request body", http.StatusBadRequest)
		return
	}

	//log.Println(string(data))
	username, err := auth.CheckLoginCreds(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response, err := auth.CreateToken(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

//HandleCreateComment handles the request to create comment
func HandleComment(w http.ResponseWriter, r *http.Request) {
	if !isAuthorized(w, r) {
		return
	}

	if r.Method == "POST" {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Could not read the request body", http.StatusBadRequest)
			return
		}

		//log.Println(string(data))
		id, err := task.Action.CreateComment(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Write([]byte("Comment Created with ID : " + strconv.Itoa(id)))
		return
	}

	if r.Method == "DELETE" {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			http.Error(w, "no comment with given id found", http.StatusNotFound)
			return
		}

		username := r.URL.Query().Get("username")
		if len(username) <= 0 || len(username) > 50 {
			http.Error(w, "username is either empty or too big", http.StatusBadRequest)
			return
		}

		err = task.Action.DeleteCmt(id, username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Write([]byte(fmt.Sprintf("Comment with id=%d and all reactions on it deleted", id)))
		return
	}

	http.Error(w, "Wrong method used. Please use POST or DELETE method", http.StatusMethodNotAllowed)
}

//HandleCreateReaction handles the request to create reaction on comment
func HandleCreateReaction(w http.ResponseWriter, r *http.Request) {
	if !isAuthorized(w, r) {
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Wrong method used. Please use POST method", http.StatusMethodNotAllowed)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read the request body", http.StatusBadRequest)
		return
	}

	//log.Println("payload received : ", string(data))
	id, err := task.Action.CreateReaction(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("Reaction created with ID : " + strconv.Itoa(id)))
}

//HandleGetWall handles the request to get wall
func HandleGetWall(w http.ResponseWriter, r *http.Request) {
	if !isAuthorized(w, r) {
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Wrong method used. Please use GET method", http.StatusMethodNotAllowed)
	}

	//log.Println(string(data))
	username := r.URL.Query().Get("username")
	if len(username) <= 0 || len(username) > 50 {
		http.Error(w, "username is either empty or too big", http.StatusBadRequest)
		return
	}

	fmt.Printf("Generating wall for user %s", username)
	wall, err := task.Action.GenerateWall(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//log.Println(string(wall))
	w.Write(wall)
}

//isAuthorized check if user is authorized or not
func isAuthorized(w http.ResponseWriter, r *http.Request) bool {
	if !config.Manager.AuthEnabled {
		//only for unit testing.  AuthEnable should always be true
		return true
	}

	authToken := r.Header.Get("Authorization")
	if authToken == "" {
		http.Error(w, "No authorization token provided", http.StatusUnauthorized)
		return false
	}

	err := auth.ValidateToken(authToken)
	if err != nil {
		http.Error(w, "Invalid authorization token provided", http.StatusUnauthorized)
		return false
	}

	return true
}

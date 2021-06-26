package auth

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/ash3798/Social-Network/database"
	"github.com/ash3798/Social-Network/structures"
)

func CheckLoginCreds(data []byte) (string, error) {
	loginCred := structures.LoginCred{}

	err := json.Unmarshal(data, &loginCred)
	if err != nil {
		log.Println("error while unmarshalling the login Creds. Error : ", err.Error())
		return "", errors.New("could not login. Use correct json structure for login info")
	}

	if loginCred.Username == "" || loginCred.Password == "" {
		log.Println("provided username or password field provided are empty")
		return "", errors.New("could not login. empty fields are provided for username or password")
	}

	err = database.Action.ValidateLoginCreds(loginCred)
	if err == nil {
		return loginCred.Username, nil
	}
	return "", err
}

package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/ash3798/Social-Network/structures"
)

//ValidateLoginCreds checks if login credentials are there in database or not
func (d database) ValidateLoginCreds(loginCreds structures.LoginCred) error {
	userSql := fmt.Sprintf(`select username from %s where username=$1 and password=$2`, usersTableName)

	result := ""
	err := d.db.QueryRow(userSql, loginCreds.Username, loginCreds.Password).Scan(&result)
	if err == sql.ErrNoRows {
		log.Println("Incorrect username or password")
		return errors.New("could not login.Incorrect username or password")
	}
	if err != nil {
		log.Println("error while checking the database. Error :", err.Error())
		return errors.New("could not login.Enter correct username and password")
	}

	return nil
}

//GetComments will get the comments from database received by user
func (d database) GetComments(username string) ([]structures.WallUnit, error) {
	wall := []structures.WallUnit{}
	userSql := fmt.Sprintf(`select id ,comment_text , sender_username , comment_time from %s where receiver_username = $1 ;`, commentTableName)

	rows, err := d.db.Query(userSql, username)
	if err != nil {
		log.Println("error while querying for comments. Error : ", err.Error())
		return nil, errors.New("could not generate wall")
	}
	defer rows.Close()

	for rows.Next() {
		tmp := structures.WallUnit{}

		err = rows.Scan(&tmp.CommentID, &tmp.CommentText, &tmp.SenderUsername, &tmp.CommentTime)
		if err != nil {
			log.Println("Error while reading the comment fields from database, Error :", err.Error())
			return nil, errors.New("could not generate wall")
		}

		wall = append(wall, tmp)
	}

	err = rows.Err()
	if err != nil {
		log.Println("Error while reading the comment fields from database, Error :", err.Error())
		return nil, errors.New("could not generate wall")
	}

	return wall, nil
}

//GetReactionCount gives the count of the reactions on the comment with given commentID
func (d database) GetReactionCount(commentID int) (map[string]int, error) {
	m := make(map[string]int)

	userSql := fmt.Sprintf(`select reaction , count(reaction) from %s where comment_id = $1 group by reaction`, reactionsTableName)

	rows, err := d.db.Query(userSql, commentID)
	if err != nil {
		log.Println("error while querying for reactions. Error : ", err.Error())
		return m, errors.New("could not generate wall")
	}
	defer rows.Close()

	for rows.Next() {
		var rc string
		var count int

		err = rows.Scan(&rc, &count)
		if err != nil {
			log.Println("Error while reading the reaction fields from database, Error :", err.Error())
			return m, errors.New("could not generate wall")
		}

		m[rc] = count
	}

	err = rows.Err()
	if err != nil {
		log.Println("Error while reading the reaction fields from database, Error :", err.Error())
		return m, errors.New("could not generate wall")
	}

	return m, nil
}

//GetCommentByID gets comment from database which matches the given id
func (d database) GetCommentByID(commentID int) (string, error) {
	userSql := fmt.Sprintf(`select receiver_username from %s where id = $1`, commentTableName)

	receiverUsername := ""
	err := d.db.QueryRow(userSql, commentID).Scan(&receiverUsername)
	if err == sql.ErrNoRows {
		log.Println("error finding the comment with the mentioned commentID , Error :", err.Error())
		return "", errors.New("could not find the comment with mentioned parent_comment_id")
	}

	if err != nil {
		log.Println("error while checking the database. Error :", err.Error())
		return "", errors.New("error finding the comment with mentioned parent_comment_id")
	}

	return receiverUsername, nil
}

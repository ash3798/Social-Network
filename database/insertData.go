package database

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ash3798/Social-Network/structures"
)

//InsertUser inserts the user into database
func (d database) InsertUser(userInfo structures.User) (int, error) {
	userSql := fmt.Sprintf(`insert into %s(username , name ,password) 
	values ($1 , $2 , $3) returning id`, usersTableName)

	lastInsertId := 0
	err := d.db.QueryRow(userSql, userInfo.Username, userInfo.Name, userInfo.Password).Scan(&lastInsertId)
	if err != nil {
		log.Println("Error while inserting user to DB. Error :", err.Error())
		return 0, err
	}

	log.Printf("user record inserted to database successfully. ID: %d", lastInsertId)
	return lastInsertId, nil
}

//InsertComment insert comment into database based on info given by user
func (d database) InsertComment(commentInfo structures.CommentInfo) (int, error) {
	userSql := ""

	lastInsertId := 0
	var err error
	if commentInfo.ParentCommentID != 0 {
		userSql = fmt.Sprintf(`insert into %s (comment_text , parent_comment_id , sender_username , receiver_username,comment_time)
		values($1, $2 , $3 , $4 , $5) returning id`, commentTableName)
		err = d.db.QueryRow(userSql, commentInfo.CommentText, commentInfo.ParentCommentID, commentInfo.SenderUsername, commentInfo.ReceiverUsername, time.Now().Unix()).Scan(&lastInsertId)
	} else {
		userSql = fmt.Sprintf(`insert into %s (comment_text , sender_username , receiver_username,comment_time)
		values($1, $2 , $3 , $4 ) returning id`, commentTableName)
		err = d.db.QueryRow(userSql, commentInfo.CommentText, commentInfo.SenderUsername, commentInfo.ReceiverUsername, time.Now().Unix()).Scan(&lastInsertId)
	}

	if err != nil {
		log.Println("Error while inserting comment to DB. Error :", err.Error())
		return 0, err
	}

	log.Printf("comment inserted to database successfully. ID: %d", lastInsertId)
	return lastInsertId, nil
}

//InsertReaction inserts the reaction into database based on the info given
func (d database) InsertReaction(reactionInfo structures.ReactionInfo) (int, error) {
	userSql := fmt.Sprintf(`insert into %s (comment_id , reaction)
	values($1 , $2) returning id`, reactionsTableName)

	lastInsertId := 0
	err := d.db.QueryRow(userSql, reactionInfo.CommentID, reactionInfo.Reaction).Scan(&lastInsertId)
	if err != nil {
		log.Println("Error while inserting reaction to DB. Error :", err.Error())
		return 0, err
	}

	log.Printf("reaction inserted to database successfully. ID: %d", lastInsertId)
	return lastInsertId, nil
}

//DeleteComment deletes the comment from database for given comment id and username
func (d database) DeleteComment(comment_id int, username string) error {
	userSql := fmt.Sprintf(`delete from %s where id = $1 and sender_username = $2`, commentTableName)

	_, err := d.db.Exec(userSql, comment_id, username)
	if err != nil {
		log.Println("Error while deleting comment from database. Error :", err.Error())
		return errors.New("could not delete comment")
	}

	return nil
}

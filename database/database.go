package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/ash3798/Social-Network/config"
	"github.com/ash3798/Social-Network/structures"
	_ "github.com/lib/pq"
)

const (
	usersTableName     = "users"
	commentTableName   = "comments"
	reactionsTableName = "reactions"
)

type DATABASE interface {
	PrepareDatabase() error
	InsertUser(userInfo structures.User) (int, error)
	InsertComment(commentInfo structures.CommentInfo) (int, error)
	InsertReaction(reactionInfo structures.ReactionInfo) (int, error)
	DeleteComment(comment_id int, username string) error
	ValidateLoginCreds(loginCreds structures.LoginCred) error
	GetComments(username string) ([]structures.WallUnit, error)
	GetReactionCount(commentID int) (map[string]int, error)
	GetCommentByID(commentID int) (string, error)
	CloseDBConnection()
}

type database struct {
	DB_DSN string
	db     *sql.DB
}

var (
	Action DATABASE = database{}
)

func InitDatabase() error {
	database := database{}
	DB_DSN := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", config.Manager.DatabaseUsername, config.Manager.DatabasePassword, config.Manager.Hostname, config.Manager.DatabasePort, config.Manager.DatabaseName)
	log.Println("db dsn :", DB_DSN)
	//create db connection
	var err error
	db, err := sql.Open("postgres", DB_DSN)
	if err != nil {
		log.Println("Failed to open DB connection.Error : ", err.Error())
		return err
	}

	err = db.Ping()
	if err != nil {
		log.Println("Not able to ping DB for testing connection. Error :", err.Error())
		return errors.New("not connected to DB")
	}

	database.DB_DSN = DB_DSN
	database.db = db

	Action = database
	log.Println("connected to DB")
	return nil
}

func (d database) CloseDBConnection() {
	d.db.Close()
}

func (d database) PrepareDatabase() error {
	userSql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id serial PRIMARY KEY,username VARCHAR(50) UNIQUE NOT NULL,name VARCHAR(50) NOT NULL,password VARCHAR(50) NOT NULL);", usersTableName)

	_, err := d.db.Exec(userSql)
	if err != nil {
		log.Printf("Error while creating %s table in database.Error : %s", usersTableName, err.Error())
		return err
	} else {
		log.Printf("%s table created successfully in DB\n", usersTableName)
	}

	//create comments table
	userSql = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
		id serial PRIMARY KEY,
		comment_text VARCHAR(255) NOT NULL,
		sender_username VARCHAR(50),
		receiver_username VARCHAR(50) ,
		parent_comment_id INT ,
		comment_time BIGINT NOT NULL ,
		CONSTRAINT fk_sender_username
			FOREIGN KEY (sender_username)
				REFERENCES users(username) ,
		CONSTRAINT fk_receiver_username
			FOREIGN KEY (receiver_username)
				REFERENCES users(username)        
	);`, commentTableName)

	_, err = d.db.Exec(userSql)
	if err != nil {
		log.Printf("Error while creating %s table in database.Error : %s", commentTableName, err.Error())
		return err
	} else {
		log.Printf("%s table created successfully in DB\n", commentTableName)
	}

	//create reactions Table
	userSql = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
		id serial PRIMARY KEY,
		reaction VARCHAR(10) CHECK (reaction = 'like' or reaction = 'dislike' or reaction = '+1'),
		comment_id INT NOT NULL,
		CONSTRAINT fk_comment_id
			FOREIGN KEY (comment_id)
				REFERENCES comments(id)
				ON DELETE CASCADE
	);`, reactionsTableName)

	_, err = d.db.Exec(userSql)
	if err != nil {
		log.Printf("Error while creating %s table in database.Error : %s", reactionsTableName, err.Error())
		return err
	} else {
		log.Printf("%s table created successfully in DB\n", reactionsTableName)
	}

	return nil
}

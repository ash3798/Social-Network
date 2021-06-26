package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/ash3798/Social-Network/config"
	"github.com/ash3798/Social-Network/database"
	"github.com/ash3798/Social-Network/structures"
)

type TASK interface {
	CreateUser(data []byte) (int, error)
	CreateComment(username string, data []byte) (int, error)
	DeleteCmt(commentID int, username string) error
	GenerateWall(username string) ([]byte, error)
	CreateReaction(data []byte) (int, error)
	CreateSubcomment(username string, data []byte) (int, error)
}

type task struct{}

var (
	Action TASK = task{}
)

//validateUserData func validates the user info passed against various constraints
func validateUserData(user structures.User) error {
	if user.Username == "" || len(user.Username) > 50 {
		return errors.New("username field is empty or too big in userdata provided")
	}

	if user.Name == "" || len(user.Name) > 50 {
		return errors.New("name field is empty or too big in userdata provided")
	}

	if user.Password == "" || len(user.Password) > 50 {
		return errors.New("password field is empty or too big in userdata provided")
	}

	return nil
}

//CreateUser creates the user according to the user info received
func (t task) CreateUser(data []byte) (int, error) {
	user := structures.User{}

	err := json.Unmarshal(data, &user)
	if err != nil {
		log.Printf("Error while unmarshelling the user structure. Error : %s", err.Error())
		return 0, errors.New("could not create user. Make sure that structure of json info is correct")
	}

	err = validateUserData(user)
	if err != nil {
		return 0, err
	}

	//log.Printf("user :  %+v \n", user)
	id, err := database.Action.InsertUser(user)
	if err != nil {
		return 0, errors.New("user not created .could not insert user to database")
	}

	return id, nil
}

//validateCommentInfo validates the commentInfo against various constraints
func validateCommentInfo(commentInfo structures.CommentInfo) error {
	if commentInfo.CommentText == "" || len(commentInfo.CommentText) > 255 {
		return errors.New("comment_text field is empty or too big in comment info")
	}
	if commentInfo.ReceiverUsername == "" || len(commentInfo.ReceiverUsername) > 50 {
		return errors.New("receiver_username field is empty or too big in comment info")
	}
	if commentInfo.SenderUsername == "" || len(commentInfo.SenderUsername) > 50 {
		return errors.New("sender_username field empty or too big in comment info")
	}
	return nil
}

//CreateComment creates the comment according to comment info received
func (t task) CreateComment(username string, data []byte) (int, error) {
	commentInfo := structures.CommentInfo{}

	err := json.Unmarshal(data, &commentInfo)
	if err != nil {
		log.Printf("Error while unmarshalling the comment . Error : %s", err.Error())
		return 0, errors.New("could not add comment. Make sure structure of the json provided is correct")
	}

	commentInfo.SenderUsername = username
	err = validateCommentInfo(commentInfo)
	if err != nil {
		return 0, err
	}

	//log.Printf("comment is : %+v \n", commentInfo)
	id, err := database.Action.InsertComment(commentInfo)
	if err != nil {
		return 0, errors.New("could not insert comment to database. check receiver_username is valid username.Also check if parent_comment_id is valid if you are doing subcomment")
	}
	return id, nil
}

func (t task) CreateSubcomment(username string, data []byte) (int, error) {
	commentInfo := structures.CommentInfo{}

	err := json.Unmarshal(data, &commentInfo)
	if err != nil {
		log.Println("error while unmarshelling the subcomment info , error :", err.Error())
		return 0, errors.New("check the structure of payload sent")
	}

	if commentInfo.ParentCommentID < 1 {
		return 0, errors.New("Invalid parent_comment_id")
	}
	receiverUsername, err := database.Action.GetCommentByID(commentInfo.ParentCommentID)
	if err != nil {
		return 0, err
	}

	commentInfo.SenderUsername = username
	commentInfo.ReceiverUsername = receiverUsername

	err = validateCommentInfo(commentInfo)
	if err != nil {
		return 0, err
	}

	//log.Printf("comment is : %+v \n", commentInfo)
	id, err := database.Action.InsertComment(commentInfo)
	if err != nil {
		return 0, errors.New("could not add subcomment")
	}
	return id, nil
}

//DeleteCmt deletes comment
func (t task) DeleteCmt(commentID int, username string) error {
	err := database.Action.DeleteComment(commentID, username)
	return err
}

//GenerateWall generates wall for user
func (t task) GenerateWall(username string) ([]byte, error) {
	wall, err := database.Action.GetComments(username)
	if err != nil {
		return []byte(""), err
	}

	///fetch reaction counts for comments on wall
	for idx, cmt := range wall {
		reactionMap, err := database.Action.GetReactionCount(cmt.CommentID)
		if err != nil {
			return []byte(""), err
		}

		wall[idx].Reaction = reactionMap
	}

	wallJson, err := json.Marshal(wall)
	if err != nil {
		log.Println("Error while marshalling json. Error :", err.Error())
		return []byte(""), errors.New("could not generate the wall")
	}

	return wallJson, nil
}

//validateReactionInfo validates the reactionInfo against various constraints
func validateReactionInfo(reactionInfo structures.ReactionInfo) error {
	if reactionInfo.CommentID == 0 {
		return errors.New("comment_id field is empty in reaction info")
	}

	_, ok := config.ReactionMap[strings.ToLower(reactionInfo.Reaction)]
	if !ok {
		errMsg := fmt.Sprintf("Not a valid reaction. Use one of the valid reactions : %v", config.AllowedReactions)
		return errors.New(errMsg)
	}
	return nil
}

//CreateReaction creates the reaction according to the reaction info passed
func (t task) CreateReaction(data []byte) (int, error) {
	reactionInfo := structures.ReactionInfo{}

	err := json.Unmarshal(data, &reactionInfo)
	if err != nil {
		log.Printf("Error while unmarshalling the reaction . Error : %s", err.Error())
		return 0, errors.New("could not add reaction. error while unmarshalling the reaction")
	}

	err = validateReactionInfo(reactionInfo)
	if err != nil {
		return 0, err
	}

	//log.Printf("reaction is : %+v \n", reactionInfo)
	id, err := database.Action.InsertReaction(reactionInfo)
	if err != nil {
		return 0, errors.New("could not insert reaction to database. Check if comment_id is valid")
	}
	return id, nil
}

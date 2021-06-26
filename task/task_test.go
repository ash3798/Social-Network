package task

import (
	"testing"

	"github.com/ash3798/Social-Network/config"
	"github.com/ash3798/Social-Network/structures"
	"github.com/stretchr/testify/assert"
)

func TestValidateUserData(t *testing.T) {
	//TODO
	userInfo := structures.User{Username: "ash", Password: "pass123", Name: "ashish"}

	err := validateUserData(userInfo)
	assert.Nil(t, err)
}

func TestValidateUserDataForInvalidUsername(t *testing.T) {
	//TODO
	userInfo := structures.User{Username: "", Password: "pass123", Name: "ashish"}
	err := validateUserData(userInfo)
	assert.NotNil(t, err)

	//testcase for username bigger than 50 chars
	userInfo = structures.User{Username: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		Password: "pass123", Name: "ashish"}
	err = validateUserData(userInfo)
	assert.NotNil(t, err)
}

func TestValidateUserDataForInvalidPassword(t *testing.T) {
	//testcase for password empty
	userInfo := structures.User{Username: "ash", Password: "", Name: "ashish"}
	err := validateUserData(userInfo)
	assert.NotNil(t, err)

	//testcase for password bigger than 50 chars
	userInfo = structures.User{Username: "ash",
		Password: "pass1233333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333",
		Name:     "ashish"}
	err = validateUserData(userInfo)
	assert.NotNil(t, err)
}

func TestValidateUserDataForInvalidName(t *testing.T) {
	//testcase for name field empty
	userInfo := structures.User{Username: "ash", Password: "pass123", Name: ""}
	err := validateUserData(userInfo)
	assert.NotNil(t, err)

	//testcase for name bigger than 50 chars
	userInfo = structures.User{Username: "ash", Password: "pass12333",
		Name: "ashishhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh"}
	err = validateUserData(userInfo)
	assert.NotNil(t, err)
}

func TestValidateCommentInfo(t *testing.T) {
	commentInfo := structures.CommentInfo{CommentText: "test comment", SenderUsername: "ash", ReceiverUsername: "nit", ParentCommentID: 0}
	err := validateCommentInfo(commentInfo)
	assert.Nil(t, err)
}

func TestValidateCommentInfoWithInvalidCommentText(t *testing.T) {
	//testcase for comment text being empty
	commentInfo := structures.CommentInfo{CommentText: "", SenderUsername: "ash", ReceiverUsername: "nit", ParentCommentID: 0}
	err := validateCommentInfo(commentInfo)
	assert.NotNil(t, err)

	//testcase for comment text being too big
	commentInfo = structures.CommentInfo{CommentText: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		SenderUsername: "ash", ReceiverUsername: "nit", ParentCommentID: 0}
	err = validateCommentInfo(commentInfo)
	assert.NotNil(t, err)
}

func TestValidateCommentInfoWithInvalidUsername(t *testing.T) {
	//testcase for sender username empty
	commentInfo := structures.CommentInfo{CommentText: "test", SenderUsername: "", ReceiverUsername: "nit", ParentCommentID: 0}
	err := validateCommentInfo(commentInfo)
	assert.NotNil(t, err)

	//testcase for receiver username empty
	commentInfo = structures.CommentInfo{CommentText: "test", SenderUsername: "ash", ReceiverUsername: "", ParentCommentID: 0}
	err = validateCommentInfo(commentInfo)
	assert.NotNil(t, err)
}

func TestValidateReactionInfo(t *testing.T) {
	config.InitReactions()
	reactionInfo := structures.ReactionInfo{CommentID: 2, Reaction: config.AllowedReactions[0]}
	err := validateReactionInfo(reactionInfo)
	assert.Nil(t, err)
}

func TestValidateReactionInfoWithInvalidCommentID(t *testing.T) {
	config.InitReactions()
	reactionInfo := structures.ReactionInfo{CommentID: 0, Reaction: config.AllowedReactions[0]}
	err := validateReactionInfo(reactionInfo)
	assert.NotNil(t, err)
}

func TestValidateReactionInfoWithInvalidReaction(t *testing.T) {
	config.InitReactions()
	reactionInfo := structures.ReactionInfo{CommentID: 1, Reaction: "invalid reaction"}
	err := validateReactionInfo(reactionInfo)
	assert.NotNil(t, err)
}

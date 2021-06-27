package structures

//created seperate package for common structures to remove import cycle

//User struct contains user information
type User struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

//CommentInfo containes comment related info
type CommentInfo struct {
	CommentText      string `json:"comment_text"`
	ParentCommentID  int    `json:"parent_comment_id"`
	SenderUsername   string `json:"sender_username"`
	ReceiverUsername string `json:"receiver_username"`
}

//ReactionInfo containes reaction related info
type ReactionInfo struct {
	CommentID int    `json:"comment_id"`
	Reaction  string `json:"reaction"`
}

//LoginCred contains the login credentials
type LoginCred struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//WallUnit is the structure of generated wall
type WallUnit struct {
	CommentID      int            `json:"comment_id"`
	CommentText    string         `json:"comment_text"`
	SenderUsername string         `json:"sender_username"`
	CommentTime    int64          `json:"timestamp"`
	Reaction       map[string]int `json:"reactions"`
}

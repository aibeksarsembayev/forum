package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no suitable record was found")

type Post struct {
	PostID       int
	Title        string
	Content      string
	UserID       int
	CategoryName string
	Username     string
	Created      time.Time
	Votes        *VoteCount
}

type User struct {
	UserID   int
	Username string
	Password string
	Confirm  string
	Email    string
	Created  time.Time
}

type Comment struct {
	CommentID   int
	PostID      int
	UserID      int
	CommentBody string
	Created     time.Time
	Username    string
	Votes       *VoteCountComment
}

type Category struct {
	CategoryID   int
	CategoryName string
}

type Vote struct {
	ID     int
	UserID int
	PostID int
	Value  bool
}

type VoteCount struct {
	Likes    uint
	Dislikes uint
}

type VoteCountComment struct {
	Likes    uint
	Dislikes uint
}

type VoteComment struct {
	ID        int
	UserID    int
	PostID    int
	CommentID int
	Value     bool
}

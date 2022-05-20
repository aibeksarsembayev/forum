package sqlite

import (
	"time"

	"git.01.alem.school/quazar/forum/pkg/models"
)

func (m *ForumModel) GetComments(id int) (*[]models.Comment, error) {
	rows, err := m.DB.Query("SELECT comment_id, post_id, user_id, comment_body, created FROM comments ORDER BY created DESC")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var comments []models.Comment

	for rows.Next() {
		c := &models.Comment{}

		// time layout format
		dateString := "2006-01-02 15:04:05 -0700 MST"
		var timeTemp string

		err = rows.Scan(&c.CommentID, &c.PostID, &c.UserID, &c.CommentBody, &timeTemp)

		c.Created, _ = time.Parse(dateString, timeTemp)

		// Count comment votes for post
		likesComment, err := m.CountVotesComment(c.CommentID, true)
		if err != nil {
			return nil, err
		}
		dislikesComment, err := m.CountVotesComment(c.CommentID, false)
		if err != nil {
			return nil, err
		}

		vc := models.VoteCountComment{
			Likes:    uint(likesComment),
			Dislikes: uint(dislikesComment),
		}
		c.Votes = &vc

		// get username
		user, err := m.GetUserByID(c.UserID)
		if err != nil {
			return nil, err
		}
		c.Username = user.Username
		if id == c.PostID {
			comments = append(comments, *c)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &comments, nil
}

func (m *ForumModel) CreateComment(post_id, user_id int, comment_body, created string) (int, error) {
	stmt, _ := m.DB.Prepare("INSERT INTO comments (post_id, user_id, comment_body, created) VALUES(?, ?, ?, ?)")

	result, err := stmt.Exec(post_id, user_id, comment_body, created)
	if err != nil {

		return 0, err
	}

	comment_id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	return int(comment_id), nil
}

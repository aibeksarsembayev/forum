package sqlite

import (
	"database/sql"
	"errors"
	"time"

	"git.01.alem.school/quazar/forum/pkg/models"
)

type ForumModel struct {
	DB *sql.DB
}

func (m *ForumModel) CreatePost(title, content string, user_id int, category_name, created string) (int, error) {
	stmt, _ := m.DB.Prepare("INSERT INTO posts (title, content, user_id, category_name, created) VALUES(?, ?, ?, ?, ?)")

	result, err := stmt.Exec(title, content, user_id, category_name, created)
	if err != nil {
		return 0, err
	}

	post_id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	return int(post_id), nil
}

func (m *ForumModel) GetPost(post_id int) (*models.Post, error) {
	stmt, _ := m.DB.Prepare("SELECT post_id, title, content, user_id, category_name, created FROM posts WHERE post_id = ?")

	row := stmt.QueryRow(post_id)

	s := &models.Post{}

	// time layout format
	dateString := "2006-01-02 15:04:05 -0700 MST"
	var timeTemp string

	err := row.Scan(&s.PostID, &s.Title, &s.Content, &s.UserID, &s.CategoryName, &timeTemp)
	//convert time string to time.Time
	s.Created, _ = time.Parse(dateString, timeTemp)
	//check if err is here required
	// if err != nil {
	// 	log.Fatal(err)
	// }

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	// get username
	user, err := m.GetUserByID(s.UserID)
	if err != nil {
		return nil, err
	}
	s.Username = user.Username

	return s, nil
}

func (m *ForumModel) LatestPost() (*[]models.Post, error) {
	rows, err := m.DB.Query("SELECT post_id, title, content, user_id, category_name, created FROM posts ORDER BY post_id DESC LIMIT 10")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		s := &models.Post{}

		// time layout format
		dateString := "2006-01-02 15:04:05 -0700 MST"
		var timeTemp string

		err = rows.Scan(&s.PostID, &s.Title, &s.Content, &s.UserID, &s.CategoryName, &timeTemp)
		s.Created, _ = time.Parse(dateString, timeTemp)
		//check if err is here required
		// if err != nil {
		// 	return nil, err
		// }

		// get username
		user, err := m.GetUserByID(s.UserID)
		if err != nil {
			return nil, err
		}
		s.Username = user.Username

		// temp Votes values
		v := &models.VoteCount{Likes: 0, Dislikes: 0}
		s.Votes = v

		posts = append(posts, *s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &posts, nil
}

func (m *ForumModel) InsertUser(username, password, email, created string) (int, error) {
	stmt, _ := m.DB.Prepare("INSERT INTO users (username, password, email, created) VALUES (?, ?, ?, ?)")

	result, err := stmt.Exec(username, password, email, created)
	if err != nil {
		return 0, err
	}

	user_id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	return int(user_id), nil
}

func (m *ForumModel) GetUser(username string) (*models.User, error) {
	stmt, _ := m.DB.Prepare("SELECT * FROM users WHERE username = ?")

	row := stmt.QueryRow(username)

	u := &models.User{}
	// time layout format
	dateString := "2006-01-02 15:04:05 -0700 MST"
	var timeTemp string

	err := row.Scan(&u.UserID, &u.Username, &u.Password, &u.Email, &timeTemp)
	//convert time string to time.Time
	u.Created, _ = time.Parse(dateString, timeTemp)
	//check if err is here required
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(timeTemp)
	// fmt.Println(u.Created)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return u, nil
}

func (m *ForumModel) GetUserByID(userID int) (*models.User, error) {
	stmt, _ := m.DB.Prepare("SELECT * FROM users WHERE user_id = ?")

	row := stmt.QueryRow(userID)

	u := &models.User{}

	// time layout format
	dateString := "2006-01-02 15:04:05 -0700 MST"
	var timeTemp string

	err := row.Scan(&u.UserID, &u.Username, &u.Password, &u.Email, &timeTemp)
	//convert time string to time.Time
	u.Created, _ = time.Parse(dateString, timeTemp)
	//check if err is here required
	// if err != nil {
	// 	log.Fatal(err)
	// }

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return u, nil
}

func (m *ForumModel) AllCategory() (*[]models.Category, error) {
	rows, err := m.DB.Query("SELECT category_id, category_name FROM category ORDER BY category_id ASC")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var categories []models.Category

	for rows.Next() {
		c := &models.Category{}

		err = rows.Scan(&c.CategoryID, &c.CategoryName)

		// no error usage? need to check
		categories = append(categories, *c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &categories, nil
}

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

func (m *ForumModel) CreateVote(post_id, user_id int, vote bool) (int, error) {
	v, err := m.GetVote(post_id, user_id)

	// check if empty row or other issues
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// create row in the table
			stmt, _ := m.DB.Prepare("INSERT INTO votes (post_id, user_id, value) VALUES (?, ?, ?)")
			result, err := stmt.Exec(post_id, user_id, vote)
			if err != nil {
				return 0, err
			}
			vote_id, err := result.LastInsertId()
			if err != nil {
				return 0, nil
			}
			return int(vote_id), nil
		} else {
			return 0, err
		}
	}

	if v.Value == vote {
		err = m.DeleteVote(v.ID)
		return 0, err
	}

	err = m.UpdateVote(v.ID, vote)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrNoRecord
		} else {
			return 0, err
		}
	}
	return v.ID, nil
}

func (m *ForumModel) UpdateVote(id int, vote bool) error {
	stmt, _ := m.DB.Prepare("UPDATE votes SET value = ? WHERE id = ?")
	_, err := stmt.Exec(vote, id)

	if err != nil {
		return err
	}
	return nil
}

func (m *ForumModel) GetVote(post_id, user_id int) (*models.Vote, error) {
	stmt, _ := m.DB.Prepare("SELECT id, value FROM votes WHERE post_id = ? AND user_id = ?")

	row := stmt.QueryRow(post_id, user_id)

	v := &models.Vote{}
	err := row.Scan(&v.ID, &v.Value)

	return v, err
}

func (m *ForumModel) DeleteVote(id int) error {
	stmt, _ := m.DB.Prepare("DELETE FROM votes WHERE id = ?")

	_, err := stmt.Exec(id)

	if err != nil {
		return err
	}
	return nil
}

func (m *ForumModel) CountVotes(post_id int, vote bool) (int, error) {
	stmt, _ := m.DB.Prepare("SELECT COUNT(value) FROM votes WHERE post_id = ? AND value = ?")

	row := stmt.QueryRow(post_id, vote)

	var countnumber int
	err := row.Scan(&countnumber)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		} else {
			return 0, err
		}
	}
	return countnumber, nil
}

func (m *ForumModel) CreateVoteComment(post_id, user_id, comment_id int, vote bool) (int, error) {
	v, err := m.GetVoteComment(post_id, user_id, comment_id)

	// check if empty row or other issues
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// create row in the table
			stmt, _ := m.DB.Prepare("INSERT INTO vote_comment (post_id, user_id, comment_id, value) VALUES (?, ?, ?, ?)")
			result, err := stmt.Exec(post_id, user_id, comment_id, vote)
			if err != nil {
				return 0, err
			}
			vote_id, err := result.LastInsertId()
			if err != nil {
				return 0, nil
			}
			return int(vote_id), nil
		} else {
			return 0, err
		}
	}

	if v.Value == vote {
		err = m.DeleteVoteComment(v.ID)
		return 0, err
	}

	err = m.UpdateVoteComment(v.ID, vote)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrNoRecord
		} else {
			return 0, err
		}
	}
	return v.ID, nil
}

func (m *ForumModel) GetVoteComment(post_id, user_id, comment_id int) (*models.VoteComment, error) {
	stmt, _ := m.DB.Prepare("SELECT id, value FROM vote_comment WHERE post_id = ? AND user_id = ? AND comment_id = ?")

	row := stmt.QueryRow(post_id, user_id, comment_id)

	v := &models.VoteComment{}
	err := row.Scan(&v.ID, &v.Value)

	return v, err
}

func (m *ForumModel) UpdateVoteComment(id int, vote bool) error {
	stmt, _ := m.DB.Prepare("UPDATE vote_comment SET value = ? WHERE id = ?")
	_, err := stmt.Exec(vote, id)

	if err != nil {
		return err
	}
	return nil
}

func (m *ForumModel) DeleteVoteComment(id int) error {
	stmt, _ := m.DB.Prepare("DELETE FROM vote_comment WHERE id = ?")

	_, err := stmt.Exec(id)

	if err != nil {
		return err
	}
	return nil
}

func (m *ForumModel) CountVotesComment(comment_id int, vote bool) (int, error) {
	stmt, _ := m.DB.Prepare("SELECT COUNT(value) FROM vote_comment WHERE comment_id = ? AND value = ?")

	row := stmt.QueryRow(comment_id, vote)

	var countnumber int
	err := row.Scan(&countnumber)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		} else {
			return 0, err
		}
	}
	return countnumber, nil
}

func (m *ForumModel) GetPostsByID(id int) (*[]models.Post, error) {
	stmt, _ := m.DB.Prepare("SELECT post_id, title, content, user_id, category_name, created FROM posts WHERE user_id = ? ORDER BY post_id DESC")
	rows, err := stmt.Query(id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		s := &models.Post{}

		// time layout format
		dateString := "2006-01-02 15:04:05 -0700 MST"
		var timeTemp string

		err = rows.Scan(&s.PostID, &s.Title, &s.Content, &s.UserID, &s.CategoryName, &timeTemp)
		s.Created, _ = time.Parse(dateString, timeTemp)

		// get username
		user, err := m.GetUserByID(s.UserID)
		if err != nil {
			return nil, err
		}
		s.Username = user.Username

		// temp Votes values
		v := &models.VoteCount{Likes: 0, Dislikes: 0}
		s.Votes = v

		posts = append(posts, *s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &posts, nil
}

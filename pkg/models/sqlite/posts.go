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

func (m *ForumModel) GetPostbyCat(category_name string) (*[]models.Post, error) {
	stmt, _ := m.DB.Prepare("SELECT post_id, title, content, user_id, category_name, created FROM posts WHERE category_name = ? ORDER BY post_id DESC")
	rows, err := stmt.Query(category_name)

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

func (m *ForumModel) GetPostbyUserVote(user_id int) (*[]models.Post, error) {
	stmt, _ := m.DB.Prepare("SELECT posts.post_id, posts.title, posts.content, posts.user_id, posts.category_name, posts.created FROM posts INNER JOIN votes ON posts.post_id = votes.post_id AND votes.user_id = ? AND votes.value = 1")
	rows, err := stmt.Query(user_id)

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

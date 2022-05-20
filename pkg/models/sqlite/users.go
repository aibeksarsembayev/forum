package sqlite

import (
	"database/sql"
	"errors"
	"time"

	"git.01.alem.school/quazar/forum/pkg/models"
)

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

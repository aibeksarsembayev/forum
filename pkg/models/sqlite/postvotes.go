package sqlite

import (
	"database/sql"
	"errors"

	"git.01.alem.school/quazar/forum/pkg/models"
)

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

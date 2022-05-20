package sqlite

import (
	"database/sql"
	"errors"

	"git.01.alem.school/quazar/forum/pkg/models"
)

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

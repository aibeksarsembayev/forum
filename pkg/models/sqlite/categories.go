package sqlite

import (
	"git.01.alem.school/quazar/forum/pkg/models"
)

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

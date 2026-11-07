package db

import (
	"database/sql"
	"idea_app/models"
)

func CreateIdea(idea *models.Idea) error {
	_, err := DB.Exec("INSERT INTO ideas (id, title, description, icon_path, is_important) VALUES (?, ?, ?, ?, ?)",
		idea.ID, idea.Title, idea.Description, idea.IconPath, idea.IsImportant)
	return err
}

func GetAllIdeas() ([]*models.Idea, error) {
	rows, err := DB.Query("SELECT id, title, description, icon_path, is_important FROM ideas ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ideas []*models.Idea
	for rows.Next() {
		i := &models.Idea{}
		if err := rows.Scan(&i.ID, &i.Title, &i.Description, &i.IconPath, &i.IsImportant); err != nil {
			return nil, err
		}
		ideas = append(ideas, i)
	}
	return ideas, nil
}

func GetIdeaByID(id string) (*models.Idea, error) {
	row := DB.QueryRow("SELECT id, title, description, icon_path, is_important FROM ideas WHERE id = ?", id)
	i := &models.Idea{}
	if err := row.Scan(&i.ID, &i.Title, &i.Description, &i.IconPath, &i.IsImportant); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	metaData, _ := GetMetaDataForParent(i.ID, "Idea")
	i.MetaData = metaData
	return i, nil
}

func UpdateIdea(idea *models.Idea) error {
	_, err := DB.Exec("UPDATE ideas SET title=?, description=?, icon_path=?, is_important=? WHERE id=?",
		idea.Title, idea.Description, idea.IconPath, idea.IsImportant, idea.ID)
	return err
}

func DeleteIdea(id string) error {
	_, err := DB.Exec("DELETE FROM ideas WHERE id = ?", id)
	return err
}

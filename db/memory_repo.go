package db

import (
	"database/sql"
	"idea_app/models"
	"time"
)

func CreateMemory(m *models.Memory) error {
	_, err := DB.Exec("INSERT INTO memories (id, title, description, memory_date, is_important) VALUES (?, ?, ?, ?, ?)",
		m.ID, m.Title, m.Description, m.Date, m.IsImportant)
	return err
}

func GetAllMemories() ([]*models.Memory, error) {
	// Dynamically sorted for UI Grouping
	rows, err := DB.Query("SELECT id, title, description, memory_date, is_important FROM memories ORDER BY strftime('%Y', memory_date) DESC, strftime('%m', memory_date) DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memories []*models.Memory
	for rows.Next() {
		m := &models.Memory{}
		var d time.Time
		if err := rows.Scan(&m.ID, &m.Title, &m.Description, &d, &m.IsImportant); err != nil {
			return nil, err
		}
		m.Date = d
		memories = append(memories, m)
	}
	return memories, nil
}

func GetMemoryByID(id string) (*models.Memory, error) {
	row := DB.QueryRow("SELECT id, title, description, memory_date, is_important FROM memories WHERE id = ?", id)
	m := &models.Memory{}
	var d time.Time
	if err := row.Scan(&m.ID, &m.Title, &m.Description, &d, &m.IsImportant); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	m.Date = d
	metaData, _ := GetMetaDataForParent(m.ID, "Memory")
	m.MetaData = metaData
	return m, nil
}

func UpdateMemory(m *models.Memory) error {
	_, err := DB.Exec("UPDATE memories SET title=?, description=?, memory_date=?, is_important=? WHERE id=?",
		m.Title, m.Description, m.Date, m.IsImportant, m.ID)
	return err
}

func DeleteMemory(id string) error {
	_, err := DB.Exec("DELETE FROM memories WHERE id = ?", id)
	return err
}

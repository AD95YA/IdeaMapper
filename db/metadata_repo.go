package db

import (
	"database/sql"
	"fmt"
	"idea_app/models"
)

func CreateMetaData(m *models.MetaData) error {
	if m.Level > 6 {
		return fmt.Errorf("maximum metadata nesting level of 6 exceeded")
	}
	_, err := DB.Exec("INSERT INTO metaData (id, parent_id, parent_type, level, title, description) VALUES (?, ?, ?, ?, ?, ?)",
		m.ID, m.ParentID, m.ParentType, m.Level, m.Title, m.Description)
	return err
}

// CreateNote is an alias for CreateMetaData for backward compatibility
func CreateNote(m *models.MetaData) error {
	return CreateMetaData(m)
}

func GetMetaDataForParent(parentID string, parentType string) ([]*models.MetaData, error) {
	rows, err := DB.Query("SELECT id, parent_id, parent_type, level, title, description FROM metaData WHERE parent_id = ? AND parent_type = ? ORDER BY created_at ASC", parentID, parentType)
	if err != nil {
		return nil, err
	}

	var metaDataList []*models.MetaData
	for rows.Next() {
		m := &models.MetaData{}
		if err := rows.Scan(&m.ID, &m.ParentID, &m.ParentType, &m.Level, &m.Title, &m.Description); err != nil {
			rows.Close()
			return nil, err
		}
		metaDataList = append(metaDataList, m)
	}
	rows.Close()

	for _, m := range metaDataList {
		// Fetch Attachments
		attRows, err := DB.Query("SELECT file_path, is_photo FROM attachments WHERE metadata_id = ?", m.ID)
		if err == nil {
			for attRows.Next() {
				var path string
				var isPhoto bool
				if err := attRows.Scan(&path, &isPhoto); err == nil {
					if isPhoto {
						m.PhotoPaths = append(m.PhotoPaths, path)
					} else {
						m.FilePaths = append(m.FilePaths, path)
					}
				}
			}
			attRows.Close()
		}

		// Recursively fetch submetaData if we are below max level
		if m.Level < 6 {
			subMetaData, err := GetMetaDataForParent(m.ID, "MetaData")
			if err != nil {
				return nil, err
			}
			m.SubMetaData = subMetaData
		}
	}
	return metaDataList, nil
}

func GetMetaDataByID(id string) (*models.MetaData, error) {
	row := DB.QueryRow("SELECT id, parent_id, parent_type, level, title, description FROM metaData WHERE id = ?", id)
	m := &models.MetaData{}
	if err := row.Scan(&m.ID, &m.ParentID, &m.ParentType, &m.Level, &m.Title, &m.Description); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return m, nil
}

// GetNoteByID is an alias for GetMetaDataByID for backward compatibility
func GetNoteByID(id string) (*models.MetaData, error) {
	return GetMetaDataByID(id)
}

func UpdateMetaData(m *models.MetaData) error {
	_, err := DB.Exec("UPDATE metaData SET title = ?, description = ? WHERE id = ?",
		m.Title, m.Description, m.ID)
	return err
}

func DeleteMetaData(id string) error {
	_, err := DB.Exec("DELETE FROM metaData WHERE id = ?", id)
	return err // Cascade triggers will automatically delete sub-metaData and attachments
}

// DeleteNote is an alias for DeleteMetaData for backward compatibility
func DeleteNote(id string) error {
	return DeleteMetaData(id)
}

func AddAttachment(metadataID, path string, isPhoto bool, attachID string) error {
	_, err := DB.Exec("INSERT INTO attachments (id, metadata_id, file_path, is_photo) VALUES (?, ?, ?, ?)",
		attachID, metadataID, path, isPhoto)
	return err
}

package db

import (
	"idea_app/models"
)

func CreateInterLink(link *models.InterLink) error {
	_, err := DB.Exec("INSERT INTO interlinks (id, source_id, source_type, target_id, target_type, name) VALUES (?, ?, ?, ?, ?, ?)",
		link.ID, link.SourceID, link.SourceType, link.TargetID, link.TargetType, link.Name)
	return err
}

func GetInterLinksForSource(sourceID string) ([]*models.InterLink, error) {
	rows, err := DB.Query("SELECT id, source_id, source_type, target_id, target_type, name FROM interlinks WHERE source_id = ? OR target_id = ?", sourceID, sourceID)
	if err != nil {
		return nil, err
	}

	var links []*models.InterLink
	for rows.Next() {
		l := &models.InterLink{}
		if err := rows.Scan(&l.ID, &l.SourceID, &l.SourceType, &l.TargetID, &l.TargetType, &l.Name); err != nil {
			rows.Close()
			return nil, err
		}
		links = append(links, l)
	}
	rows.Close()

	for _, l := range links {
		// A Link only supports 1 singular top-level MetaData attached directly to it.
		metaData, _ := GetMetaDataForParent(l.ID, "InterLink")
		if len(metaData) > 0 {
			l.LinkMetaData = metaData[0]
		}
	}
	return links, nil
}

func GetAllInterLinks() ([]*models.InterLink, error) {
	rows, err := DB.Query("SELECT id, source_id, source_type, target_id, target_type, name FROM interlinks ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}

	var links []*models.InterLink
	for rows.Next() {
		l := &models.InterLink{}
		if err := rows.Scan(&l.ID, &l.SourceID, &l.SourceType, &l.TargetID, &l.TargetType, &l.Name); err != nil {
			rows.Close()
			return nil, err
		}
		links = append(links, l)
	}
	rows.Close()

	for _, l := range links {
		metaData, _ := GetMetaDataForParent(l.ID, "InterLink")
		if len(metaData) > 0 {
			l.LinkMetaData = metaData[0]
		}
	}
	return links, nil
}

func DeleteInterLink(id string) error {
	_, err := DB.Exec("DELETE FROM interlinks WHERE id = ?", id)
	return err
}

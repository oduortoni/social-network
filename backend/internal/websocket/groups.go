package websocket

import (
	"database/sql"
)

type DBGroupMemberFetcher struct {
	DB *sql.DB
}

func NewDBGroupMemberFetcher(db *sql.DB) *DBGroupMemberFetcher {
	return &DBGroupMemberFetcher{DB: db}
}

func (g *DBGroupMemberFetcher) GetGroupMemberIDs(groupID string) ([]int64, error) {
	rows, err := g.DB.Query(`
		SELECT user_id FROM Group_Members
		WHERE group_id = ? AND is_accepted = 1
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err == nil {
			ids = append(ids, id)
		}
	}
	return ids, nil
}

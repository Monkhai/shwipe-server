package db

import "github.com/jackc/pgx/v5"

func (db *DB) GetAllSessions(userID string) ([]*DBSession, error) {
	query, err := db.CreateQuery("queries/sessions/query_getAllSessions.sql", "getAllSessions", DB_Tables)
	if err != nil {
		return nil, err
	}

	rows, err := db.RunQuery(query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	sessions := make([]*DBSession, 0)
	for rows.Next() {
		session, err := pgx.RowToStructByName[DBSession](rows)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, &session)
	}

	return sessions, nil
}

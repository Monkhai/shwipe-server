package db

import "github.com/jackc/pgx/v5"

func (db *DB) GetOneSession(sessionID string) (*DBSession, error) {
	query, err := db.CreateQuery("queries/sessions/query_getOneSession.sql", "getOneSession", DB_Tables)
	if err != nil {
		return nil, err
	}

	rows, err := db.RunQuery(query, sessionID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	session, err := pgx.RowToStructByName[DBSession](rows)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

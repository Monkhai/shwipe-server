package db

import "time"

func (db *DB) InsertSession(sessionID string) error {
	query, err := db.CreateQuery("queries/sessions/query_insertSession.sql", "insertSession", DB_Tables)
	if err != nil {
		return err
	}

	rows, err := db.RunQuery(query, sessionID, time.Now())
	if err != nil {
		return err
	}

	defer rows.Close()

	return nil
}

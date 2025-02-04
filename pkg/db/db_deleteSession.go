package db

func (db *DB) DeleteSession(sessionID string) error {
	query, err := db.CreateQuery("queries/sessions/query_deleteSession.sql", "deleteSession", DB_Tables)
	if err != nil {
		return err
	}

	_, err = db.RunQuery(query, sessionID)
	if err != nil {
		return err
	}

	return nil
}

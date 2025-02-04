package db

func (db *DB) InsertSessionUser(sessionID string, userID string) error {
	query, err := db.CreateQuery("queries/sessions/query_insertSessionUser.sql", "insertSessionUser", DB_Tables)
	if err != nil {
		return err
	}

	err = db.ExecuteQuery(query, sessionID, userID)
	if err != nil {
		return err
	}

	return nil
}

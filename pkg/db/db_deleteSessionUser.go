package db

func (db *DB) DeleteSessionUser(sessionID string, userID string) error {
	query, err := db.CreateQuery("queries/session_users/query_deleteSessionUser.sql", "deleteSessionUser", DB_Tables)
	if err != nil {
		return err
	}

	err = db.ExecuteQuery(query, sessionID, userID)
	if err != nil {
		return err
	}

	return nil
}

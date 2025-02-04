package db

func (db *DB) DeleteAllSessions() error {
	query, err := db.CreateQuery("queries/sessions/query_deleteAllSessions.sql", "deleteAllSessions", DB_Tables)
	if err != nil {
		return err
	}
	rows, err := db.RunQuery(query)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

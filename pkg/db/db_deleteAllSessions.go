package db

import "context"

func (db *DB) DeleteAllSessions(ctx context.Context) error {
	query, err := db.CreateQuery("queries/sessions/query_deleteAllSessions.sql", "deleteAllSessions", DB_Tables)
	if err != nil {
		return err
	}
	_, err = db.pool.Exec(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

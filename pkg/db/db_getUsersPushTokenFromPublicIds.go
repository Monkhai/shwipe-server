package db

import "github.com/jackc/pgx/v5"

func (db *DB) GetUsersPushTokenFromPublicIds(publicIds []string) ([]string, error) {
	query, err := db.CreateQuery("queries/users/query_getUsersPushTokensFromPublicIds.sql", "getUsersPushTokenFromPublicIds", DB_Tables)
	if err != nil {
		return nil, err
	}

	rows, err := db.RunQuery(query, publicIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pushTokens, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, err
	}

	return pushTokens, nil
}

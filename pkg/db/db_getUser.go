package db

import (
	"github.com/jackc/pgx/v5"
)

func (db *DB) GetUser(id string) (*DBUser, error) {
	query, err := db.CreateQuery("queries/users/query_getUser.sql", "getUser", DB_Tables)
	if err != nil {
		return nil, err
	}

	rows, err := db.RunQuery(query, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[DBUser])

	if err != nil {
		return nil, err
	}

	return &user, nil
}

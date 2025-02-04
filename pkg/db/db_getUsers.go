package db

import (
	"log"

	"github.com/jackc/pgx/v5"
)

func (db *DB) GetUsers() ([]*DBUser, error) {
	query, err := db.CreateQuery("queries/users/query_getUsers.sql", "getUsers", DB_Tables)
	if err != nil {
		log.Println(err, "from db_getUsers.go")
		return nil, err
	}

	rows, err := db.RunQuery(query)
	if err != nil {
		log.Println(err, "from db_getUsers.go")
		return nil, err
	}
	defer rows.Close()

	users := make([]*DBUser, 0)
	for rows.Next() {
		u, err := pgx.RowToStructByName[DBUser](rows)
		if err != nil {
			log.Println(err, "from db_getUsers.go")
			return nil, err
		}
		users = append(users, &u)
	}

	return users, nil
}

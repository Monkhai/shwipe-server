package db

func (db *DB) GetGroupMembersPushTokens(groupId string) ([]string, error) {
	query, err := db.CreateQuery("queries/groups/query_getGroupMembersPushTokens.sql", "getGroupMembersPushTokens", DB_Tables)
	if err != nil {
		return nil, err
	}

	rows, err := db.RunQuery(query, groupId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	pushTokens := make([]string, 0)
	for rows.Next() {
		var pushToken string
		err := rows.Scan(&pushToken)
		if err != nil {
			return nil, err
		}
		pushTokens = append(pushTokens, pushToken)
	}
	return pushTokens, nil
}

DELETE FROM {{.DB_Tables.SessionUsers}}
WHERE session_id = $1 AND user_id = $2
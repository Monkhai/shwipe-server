DELETE FROM {{ .SessionUsers.TableName }}
WHERE session_id = $1 AND user_id = $2
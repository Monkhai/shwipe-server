SELECT u.{{ .UsersTable.ExpoPushToken }}
FROM {{ .UsersTable.TableName }} u
JOIN {{ .GroupMembersTable.TableName }} gm ON u.{{ .UsersTable.ID }} = gm.{{ .GroupMembersTable.UserID }}
WHERE gm.{{ .GroupMembersTable.GroupID }} = $1

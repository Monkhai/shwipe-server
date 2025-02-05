    SELECT
      u.{{ .UsersTable.ExpoPushToken }}
    FROM {{ .UsersTable.TableName }} u
    INNER JOIN {{ .PublicUserIdsTable.TableName }} pui ON pui.{{ .PublicUserIdsTable.ID }} = u.{{ .UsersTable.ID }}
    WHERE pui.{{ .PublicUserIdsTable.PublicID }} = ANY($1)
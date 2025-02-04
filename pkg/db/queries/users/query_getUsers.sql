SELECT
    u.{{.UsersTable.ID}},
    u.{{.UsersTable.DisplayName}},
    u.{{.UsersTable.PhotoURL}},
    u.{{.UsersTable.ExpoPushToken}},
    pui.{{.PublicUserIdsTable.PublicID}} AS public_id
FROM {{.UsersTable.TableName}} u
    JOIN {{.PublicUserIdsTable.TableName}} pui ON u.{{.UsersTable.ID}} = pui.{{.PublicUserIdsTable.ID}}
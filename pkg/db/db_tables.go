package db

var DB_Tables = struct {
	UsersTable         UsersTable
	PublicUserIdsTable PublicUserIdsTable
}{
	UsersTable:         UsersTableVar,
	PublicUserIdsTable: PublicUserIdsTableVar,
}

//===============================================================

type UsersTable struct {
	TableName     string
	ID            string
	DisplayName   string
	PhotoURL      string
	ExpoPushToken string
}

var UsersTableVar = UsersTable{
	TableName:     "users",
	ID:            "id",
	DisplayName:   "display_name",
	PhotoURL:      "photo_url",
	ExpoPushToken: "expo_push_token",
}

//===============================================================

type PublicUserIdsTable struct {
	TableName string
	ID        string
	PublicID  string
}

var PublicUserIdsTableVar = PublicUserIdsTable{
	TableName: "public_user_ids",
	ID:        "id",
	PublicID:  "public_id",
}

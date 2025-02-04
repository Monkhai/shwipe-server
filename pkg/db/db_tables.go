package db

var DB_Tables = struct {
	UsersTable         UsersTable
	PublicUserIdsTable PublicUserIdsTable
	SessionsTable      SessionsTable
	SessionUsersTable  SessionUsersTable
	FriendshipsTable   FriendshipsTable
}{
	UsersTable:         UsersTableVar,
	PublicUserIdsTable: PublicUserIdsTableVar,
	SessionsTable:      SessionsTableVar,
	SessionUsersTable:  SessionUsersTableVar,
	FriendshipsTable:   FriendshipsTableVar,
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

//===============================================================

type SessionsTable struct {
	TableName string
	ID        string
	CreatedAt string
}

var SessionsTableVar = SessionsTable{
	TableName: "sessions",
	ID:        "id",
	CreatedAt: "created_at",
}

// ===============================================================
type SessionUsersTable struct {
	TableName string
	ID        string
	SessionID string
	UserID    string
}

var SessionUsersTableVar = SessionUsersTable{
	TableName: "session_users",
	ID:        "id",
	SessionID: "session_id",
	UserID:    "user_id",
}

// ===============================================================

type FriendshipsTable struct {
	TableName string
	UserID1   string
	UserID2   string
}

var FriendshipsTableVar = FriendshipsTable{
	TableName: "friendships",
	UserID1:   "user_id_1",
	UserID2:   "user_id_2",
}

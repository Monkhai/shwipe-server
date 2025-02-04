package db

type DBUser struct {
	ID            string `db:"id"`
	DisplayName   string `db:"display_name"`
	PhotoURL      string `db:"photo_url"`
	ExpoPushToken string `db:"expo_push_token"`
	PublicID      string `db:"public_id"`
}

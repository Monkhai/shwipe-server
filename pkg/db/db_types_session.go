package db

const (
	SessionStatusActive   = "active"
	SessionStatusInactive = "inactive"
)

type DBSession struct {
	ID        string `db:"id"`
	CreatedAt string `db:"created_at"`
	Status    string `db:"status"`
}

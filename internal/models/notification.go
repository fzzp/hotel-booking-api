package models

type Notification struct {
	ID        uint   `db:"id" json:"id"`
	UserID    uint   `db:"user_id" json:"userId"`
	Message   string `db:"message" json:"message"`
	IsRead    int8   `db:"is_read" json:"isRead"`
	CreatedAt string `db:"created_at" json:"createdAt"`
	IsDeleted int8   `db:"is_deleted" json:"-"`
}

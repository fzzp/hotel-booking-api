package models

type Hotel struct {
	ID        uint   `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	Address   string `db:"address" json:"address"`
	Logo      string `db:"logo" json:"logo"`
	CreatedAt string `db:"created_at" json:"createdAt"`
	UpdatedAt string `db:"updated_at" json:"updatedAt"`
	IsDeleted int8   `db:"is_deleted" json:"-"`
}

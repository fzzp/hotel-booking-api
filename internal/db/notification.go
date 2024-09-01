package db

import "github.com/fzzp/hotel-booking-api/internal/models"

type NotificationRepo interface {
	InsertOne(hotel *models.Notification) (id uint, err error)
	GetOneByUq(map[string]string) (*models.Notification, error)
	UpdateOne(uid uint, hotel *models.Notification) error
}

var _ NotificationRepo = (*notificationRepo)(nil)

type notificationRepo struct {
	DB Queryable
}

// GetOneByUq implements NotificationRepo.
func (n *notificationRepo) GetOneByUq(map[string]string) (*models.Notification, error) {
	panic("unimplemented")
}

// InsertOne implements NotificationRepo.
func (n *notificationRepo) InsertOne(hotel *models.Notification) (id uint, err error) {
	panic("unimplemented")
}

// UpdateOne implements NotificationRepo.
func (n *notificationRepo) UpdateOne(uid uint, hotel *models.Notification) error {
	panic("unimplemented")
}

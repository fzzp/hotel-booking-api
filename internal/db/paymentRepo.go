package db

import "github.com/fzzp/hotel-booking-api/internal/models"

type PaymentRepo interface {
	InsertOne(hotel *models.Payment) (id uint, err error)
	GetOneByUq(map[string]string) (*models.Payment, error)
	UpdateOne(uid uint, hotel *models.Payment) error
}

var _ PaymentRepo = (*paymentRepo)(nil)

type paymentRepo struct {
	DB Queryable
}

// GetOneByUq implements PaymentRepo.
func (p *paymentRepo) GetOneByUq(map[string]string) (*models.Payment, error) {
	panic("unimplemented")
}

// InsertOne implements PaymentRepo.
func (p *paymentRepo) InsertOne(hotel *models.Payment) (id uint, err error) {
	panic("unimplemented")
}

// UpdateOne implements PaymentRepo.
func (p *paymentRepo) UpdateOne(uid uint, hotel *models.Payment) error {
	panic("unimplemented")
}

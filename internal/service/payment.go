package service

type PaymentService interface{}

type paymentService struct{}

var _ PaymentService = (*paymentService)(nil)

func NewPaymentService() PaymentService {
	ps := &paymentService{}

	return ps
}

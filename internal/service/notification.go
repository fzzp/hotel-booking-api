package service

type NotificationService interface{}

type notificationService struct{}

var _ NotificationService = (*notificationService)(nil)

func NewNotificationService() NotificationService {
	ps := &notificationService{}

	return ps
}

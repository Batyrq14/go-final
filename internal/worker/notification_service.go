package worker

import (
	"qasynda/internal/models"

	"github.com/google/uuid"
)

type NotificationService struct {
	worker *NotificationWorker
}

func NewNotificationService(worker *NotificationWorker) *NotificationService {
	return &NotificationService{
		worker: worker,
	}
}

func (s *NotificationService) NotifyBookingCreated(booking *models.Booking, providerEmail string) {
	// Notify provider about new booking
	s.worker.Enqueue(NotificationTask{
		UserID:  booking.ProviderID,
		Message: "You have a new booking request",
		Type:    NotificationEmail,
	})
}

func (s *NotificationService) NotifyBookingStatusChanged(booking *models.Booking, clientEmail string) {
	var message string
	switch booking.Status {
	case models.StatusAccepted:
		message = "Your booking has been accepted"
	case models.StatusRejected:
		message = "Your booking has been rejected"
	case models.StatusCompleted:
		message = "Your booking has been completed"
	case models.StatusCancelled:
		message = "Your booking has been cancelled"
	default:
		return
	}

	s.worker.Enqueue(NotificationTask{
		UserID:  booking.ClientID,
		Message: message,
		Type:    NotificationEmail,
	})
}

func (s *NotificationService) NotifyReviewCreated(providerID uuid.UUID) {
	s.worker.Enqueue(NotificationTask{
		UserID:  providerID,
		Message: "You received a new review",
		Type:    NotificationEmail,
	})
}

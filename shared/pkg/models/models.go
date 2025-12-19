package models

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
	Phone    string `json:"phone"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
	Phone    string `json:"phone"`
}

type AuthResponse struct {
	Token string        `json:"token"`
	User  *UserResponse `json:"user"`
}

type ValidateTokenRequest struct {
	Token string `json:"token"`
}

type GetUserRequest struct {
	UserID string `json:"user_id"`
}

type ListProvidersRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type ProviderResponse struct {
	User            *UserResponse `json:"user"`
	Location        string        `json:"location"`
	HourlyRate      float64       `json:"hourly_rate"`
	ExperienceYears int           `json:"experience_years"`
	Bio             string        `json:"bio"`
	IsAvailable     bool          `json:"is_available"`
	Rating          float64       `json:"rating"`
	ProviderID      string        `json:"provider_id"`
}

type ListProvidersResponse struct {
	Providers []*ProviderResponse `json:"providers"`
}

type CreateServiceRequest struct {
	UserID      string  `json:"user_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
}

type ServiceResponse struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ProviderID  string  `json:"provider_id"`
}

type GetServicesRequest struct {
	Category string `json:"category"`
}

type GetServicesResponse struct {
	Services []*ServiceResponse `json:"services"`
}

type CreateBookingRequest struct {
	ServiceID     string `json:"service_id"`
	UserID        string `json:"user_id"`
	ProviderID    string `json:"provider_id"`
	ScheduledTime string `json:"scheduled_time"`
}

type BookingResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type ListBookingsRequest struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

type BookingDetails struct {
	ID             string `json:"id"`
	ServiceID      string `json:"service_id"`
	ClientID       string `json:"client_id"`
	ProviderID     string `json:"provider_id"`
	Status         string `json:"status"`
	ScheduledTime  string `json:"scheduled_time"`
	ServiceTitle   string `json:"service_title"`
	OtherPartyName string `json:"other_party_name"`
}

type ListBookingsResponse struct {
	Bookings []*BookingDetails `json:"bookings"`
}

type UpdateBookingStatusRequest struct {
	BookingID string `json:"booking_id"`
	Status    string `json:"status"`
	UserID    string `json:"user_id"`
}

type GetHistoryRequest struct {
	UserID1 string `json:"user_id_1"`
	UserID2 string `json:"user_id_2"`
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
}

type Message struct {
	ID         string `json:"id"`
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	Content    string `json:"content"`
	Timestamp  string `json:"timestamp"`
}

type GetHistoryResponse struct {
	Messages []*Message `json:"messages"`
}

type SendMessageRequest struct {
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	Content    string `json:"content"`
}

type SendMessageResponse struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
}

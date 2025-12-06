package main

import (
	"context"
	"time"

	"qasynda/shared/pkg/logger"
	pb "qasynda/shared/proto"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedMarketplaceServiceServer
	store *Store
}

func NewServer(store *Store) *Server {
	return &Server{store: store}
}

func (s *Server) CreateService(ctx context.Context, req *pb.CreateServiceRequest) (*pb.ServiceResponse, error) {
	id := uuid.New()
	service := &Service{
		ID:          id,
		Name:        req.Title, // Mapping Title to Name
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.store.CreateService(ctx, service); err != nil {
		logger.Error("failed to create service", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.ServiceResponse{
		Id:          service.ID.String(),
		Title:       service.Name,
		Description: service.Description,
	}, nil
}

func (s *Server) GetServices(ctx context.Context, req *pb.GetServicesRequest) (*pb.GetServicesResponse, error) {
	services, err := s.store.ListServices(ctx)
	if err != nil {
		logger.Error("failed to list services", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	var pbServices []*pb.ServiceResponse
	for _, svc := range services {
		pbServices = append(pbServices, &pb.ServiceResponse{
			Id:          svc.ID.String(),
			Title:       svc.Name,
			Description: svc.Description,
		})
	}

	return &pb.GetServicesResponse{
		Services: pbServices,
	}, nil
}

func (s *Server) CreateBooking(ctx context.Context, req *pb.CreateBookingRequest) (*pb.BookingResponse, error) {
	serviceID, err := uuid.Parse(req.ServiceId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid service id")
	}
	clientID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	scheduledTime, err := time.Parse(time.RFC3339, req.ScheduledTime)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid scheduled time format (use ISO8601/RFC3339)")
	}

	providerID, err := uuid.Parse(req.ProviderId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid provider id")
	}

	id := uuid.New()
	booking := &Booking{
		ID:            id,
		ClientID:      clientID,
		ProviderID:    providerID,
		ServiceID:     serviceID,
		ScheduledDate: scheduledTime,
		Status:        "pending",
		DurationHours: 1.0, // Default duration
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.store.CreateBooking(ctx, booking); err != nil {
		logger.Error("failed to create booking", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.BookingResponse{
		Id:     booking.ID.String(),
		Status: booking.Status,
	}, nil
}

func (s *Server) ListBookings(ctx context.Context, req *pb.ListBookingsRequest) (*pb.ListBookingsResponse, error) {
	bookings, err := s.store.ListBookings(ctx, req.UserId, req.Role)
	if err != nil {
		logger.Error("failed to list bookings", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	var pbBookings []*pb.BookingDetails
	for _, b := range bookings {
		// TODO: Fetch Service Title and User Names. For now returning IDs.
		pbBookings = append(pbBookings, &pb.BookingDetails{
			Id:             b.ID.String(),
			ServiceId:      b.ServiceID.String(),
			ClientId:       b.ClientID.String(),
			ProviderId:     b.ProviderID.String(),
			Status:         b.Status,
			ScheduledTime:  b.ScheduledDate.Format(time.RFC3339),
			ServiceTitle:   "Service " + b.ServiceID.String(), // Placeholder
			OtherPartyName: "User " + b.ClientID.String(),     // Placeholder
		})
	}

	return &pb.ListBookingsResponse{
		Bookings: pbBookings,
	}, nil
}

func (s *Server) UpdateBookingStatus(ctx context.Context, req *pb.UpdateBookingStatusRequest) (*pb.BookingResponse, error) {

	err := s.store.UpdateBookingStatus(ctx, req.BookingId, req.Status)
	if err != nil {
		logger.Error("failed to update booking status", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.BookingResponse{
		Id:     req.BookingId,
		Status: req.Status,
	}, nil
}

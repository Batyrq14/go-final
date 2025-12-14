package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"qasynda/shared/pkg/config"
	"qasynda/shared/pkg/models"
)

type Clients struct {
	User        *UserClient
	Marketplace *MarketplaceClient
	Chat        *ChatClient
}

func InitClients(cfg *config.Config) *Clients {
	httpClient := &http.Client{}
	return &Clients{
		User:        NewUserClient(cfg.Services.UserUrl, httpClient),
		Marketplace: NewMarketplaceClient(cfg.Services.MarketplaceUrl, httpClient),
		Chat:        NewChatClient(cfg.Services.ChatUrl, httpClient),
	}
}

// User Client
type UserClient struct {
	BaseURL string
	Client  *http.Client
}

func NewUserClient(baseURL string, client *http.Client) *UserClient {
	return &UserClient{BaseURL: baseURL, Client: client}
}

func (c *UserClient) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	return doPost[models.RegisterRequest, models.AuthResponse](c.Client, c.BaseURL+"/register", req)
}

func (c *UserClient) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	return doPost[models.LoginRequest, models.AuthResponse](c.Client, c.BaseURL+"/login", req)
}

func (c *UserClient) ValidateToken(ctx context.Context, req *models.ValidateTokenRequest) (*models.UserResponse, error) {
	return doPost[models.ValidateTokenRequest, models.UserResponse](c.Client, c.BaseURL+"/validate", req)
}

func (c *UserClient) GetUser(ctx context.Context, req *models.GetUserRequest) (*models.UserResponse, error) {
	url := fmt.Sprintf("%s/users/%s", c.BaseURL, req.UserID)
	return doGet[models.UserResponse](c.Client, url)
}

func (c *UserClient) ListProviders(ctx context.Context, req *models.ListProvidersRequest) (*models.ListProvidersResponse, error) {
	url := fmt.Sprintf("%s/providers?limit=%d&offset=%d", c.BaseURL, req.Limit, req.Offset)
	return doGet[models.ListProvidersResponse](c.Client, url)
}

func (c *UserClient) UpdateProviderStatus(ctx context.Context, userID string, isAvailable bool) (*map[string]interface{}, error) {
	url := fmt.Sprintf("%s/providers/%s/status", c.BaseURL, userID)
	reqBody := map[string]bool{"is_available": isAvailable}
	return doPut[map[string]bool, map[string]interface{}](c.Client, url, &reqBody)
}

func (c *UserClient) GetProviderStatus(ctx context.Context, userID string) (*map[string]interface{}, error) {
	url := fmt.Sprintf("%s/providers/%s/status", c.BaseURL, userID)
	return doGet[map[string]interface{}](c.Client, url)
}

// Marketplace Client
type MarketplaceClient struct {
	BaseURL string
	Client  *http.Client
}

func NewMarketplaceClient(baseURL string, client *http.Client) *MarketplaceClient {
	return &MarketplaceClient{BaseURL: baseURL, Client: client}
}

func (c *MarketplaceClient) CreateService(ctx context.Context, req *models.CreateServiceRequest) (*models.ServiceResponse, error) {
	return doPost[models.CreateServiceRequest, models.ServiceResponse](c.Client, c.BaseURL+"/services", req)
}

func (c *MarketplaceClient) GetServices(ctx context.Context, req *models.GetServicesRequest) (*models.GetServicesResponse, error) {
	url := c.BaseURL + "/services"
	if req.Category != "" {
		url += "?category=" + req.Category
	}
	return doGet[models.GetServicesResponse](c.Client, url)
}

func (c *MarketplaceClient) CreateBooking(ctx context.Context, req *models.CreateBookingRequest) (*models.BookingResponse, error) {
	return doPost[models.CreateBookingRequest, models.BookingResponse](c.Client, c.BaseURL+"/bookings", req)
}

func (c *MarketplaceClient) ListBookings(ctx context.Context, req *models.ListBookingsRequest) (*models.ListBookingsResponse, error) {
	url := fmt.Sprintf("%s/bookings?user_id=%s&role=%s", c.BaseURL, req.UserID, req.Role)
	return doGet[models.ListBookingsResponse](c.Client, url)
}

func (c *MarketplaceClient) UpdateBookingStatus(ctx context.Context, req *models.UpdateBookingStatusRequest) (*models.BookingResponse, error) {
	url := fmt.Sprintf("%s/bookings/%s/status", c.BaseURL, req.BookingID)
	return doPut[models.UpdateBookingStatusRequest, models.BookingResponse](c.Client, url, req)
}

// Chat Client
type ChatClient struct {
	BaseURL string
	Client  *http.Client
}

func NewChatClient(baseURL string, client *http.Client) *ChatClient {
	return &ChatClient{BaseURL: baseURL, Client: client}
}

func (c *ChatClient) GetHistory(ctx context.Context, req *models.GetHistoryRequest) (*models.GetHistoryResponse, error) {
	url := fmt.Sprintf("%s/history?user_id_1=%s&user_id_2=%s&limit=%d&offset=%d",
		c.BaseURL, req.UserID1, req.UserID2, req.Limit, req.Offset)
	return doGet[models.GetHistoryResponse](c.Client, url)
}

// HTTP Helpers
func doPost[Req, Resp any](client *http.Client, url string, req *Req) (*Resp, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result Resp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func doGet[Resp any](client *http.Client, url string) (*Resp, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result Resp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func doPut[Req, Resp any](client *http.Client, url string, req *Req) (*Resp, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result Resp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

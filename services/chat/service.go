package main

import (
	"net/http"
	"strconv"

	"qasynda/shared/pkg/logger"
	"qasynda/shared/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Server struct {
	store *Store
}

func NewServer(store *Store) *Server {
	return &Server{store: store}
}

func (s *Server) GetHistory(c *gin.Context) {
	userID1 := c.Query("user_id_1")
	userID2 := c.Query("user_id_2")
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	if userID1 == "" || userID2 == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id_1 and user_id_2 are required"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}
	if limit > 200 {
		limit = 200
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	u1, err := uuid.Parse(userID1)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id_1"})
		return
	}
	u2, err := uuid.Parse(userID2)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id_2"})
		return
	}

	messages, err := s.store.GetHistory(c.Request.Context(), u1, u2, limit, offset)
	if err != nil {
		logger.Error("failed to get history", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	var respMessages []*models.Message
	for _, m := range messages {
		respMessages = append(respMessages, &models.Message{
			ID:         m.ID.String(),
			SenderID:   m.SenderID.String(),
			ReceiverID: m.ReceiverID.String(),
			Content:    m.Content,
			Timestamp:  m.CreatedAt.String(),
		})
	}

	c.JSON(http.StatusOK, &models.GetHistoryResponse{
		Messages: respMessages,
	})
}

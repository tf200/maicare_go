package ws

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"maicare_go/token"
	"maicare_go/util"
)

const WsTicketQueryKey = "ticket"

var (
	errInvalidWSTicket = fmt.Errorf("invalid websocket ticket")
	errExpiredWSTicket = fmt.Errorf("expired websocket ticket")
)

type TicketManager struct {
	redisClient *redis.Client
	ttl         time.Duration
	keyPrefix   string
}

type wsTicketRecord struct {
	PayloadID  string    `json:"payload_id"`
	UserID     string    `json:"user_id"`
	EmployeeID string    `json:"employee_id"`
	IssuedAt   time.Time `json:"issued_at"`
	ExpiresAt  time.Time `json:"expires_at"`
}

func NewTicketManager(config util.Config) (*TicketManager, error) {
	client := redis.NewClient(&redis.Options{
		Addr:      config.RedisHost,
		Password:  config.RedisPassword,
		Username:  "",
		TLSConfig: nil,
	})

	return &TicketManager{
		redisClient: client,
		ttl:         config.WsTicketTTL,
		keyPrefix:   "ws_ticket:",
	}, nil
}

func (m *TicketManager) Close() error {
	if m == nil || m.redisClient == nil {
		return nil
	}
	return m.redisClient.Close()
}

func (m *TicketManager) Issue(ctx context.Context, payload *token.Payload) (string, time.Time, error) {
	if m == nil || m.redisClient == nil {
		return "", time.Time{}, fmt.Errorf("websocket ticket manager unavailable")
	}
	if payload == nil {
		return "", time.Time{}, fmt.Errorf("missing auth payload")
	}

	now := time.Now()
	expiresAt := now.Add(m.ttl)
	if payload.ExpiresAt.Before(expiresAt) {
		expiresAt = payload.ExpiresAt
	}
	if !expiresAt.After(now) {
		return "", time.Time{}, errExpiredWSTicket
	}

	record := wsTicketRecord{
		PayloadID:  payload.ID.String(),
		UserID:     payload.UserId.String(),
		EmployeeID: payload.EmployeeID.String(),
		IssuedAt:   now,
		ExpiresAt:  expiresAt,
	}

	recordBytes, err := json.Marshal(record)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("marshal websocket ticket payload: %w", err)
	}

	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return "", time.Time{}, errExpiredWSTicket
	}

	ticket, err := generateTicket(32)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("generate websocket ticket: %w", err)
	}

	key := m.keyPrefix + ticket
	stored, err := m.redisClient.SetNX(ctx, key, recordBytes, ttl).Result()
	if err != nil {
		return "", time.Time{}, fmt.Errorf("store websocket ticket: %w", err)
	}
	if !stored {
		return "", time.Time{}, fmt.Errorf("failed to store websocket ticket")
	}

	return ticket, expiresAt, nil
}

func (m *TicketManager) Consume(ctx context.Context, ticketValue string) (*token.Payload, error) {
	if m == nil || m.redisClient == nil {
		return nil, fmt.Errorf("websocket ticket manager unavailable")
	}
	if ticketValue == "" {
		return nil, errInvalidWSTicket
	}

	key := m.keyPrefix + ticketValue
	recordRaw, err := m.redisClient.GetDel(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, errInvalidWSTicket
		}
		return nil, fmt.Errorf("consume websocket ticket: %w", err)
	}

	var record wsTicketRecord
	if err := json.Unmarshal(recordRaw, &record); err != nil {
		return nil, errInvalidWSTicket
	}

	if time.Now().After(record.ExpiresAt) {
		return nil, errExpiredWSTicket
	}

	payloadID, err := uuid.Parse(record.PayloadID)
	if err != nil {
		return nil, errInvalidWSTicket
	}

	userID, err := uuid.Parse(record.UserID)
	if err != nil {
		return nil, errInvalidWSTicket
	}

	employeeID, err := uuid.Parse(record.EmployeeID)
	if err != nil {
		return nil, errInvalidWSTicket
	}

	return &token.Payload{
		ID:         payloadID,
		UserId:     userID,
		EmployeeID: employeeID,
		TokenType:  token.AccessToken,
		IssuedAt:   record.IssuedAt,
		ExpiresAt:  record.ExpiresAt,
	}, nil
}

func generateTicket(size int) (string, error) {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

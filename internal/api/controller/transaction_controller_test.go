package controller

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type mockTransactionService struct {
	tripayWebhookFn         func(ctx context.Context, rawBody []byte, payload dto.TripayWebhookRequest, callbackSignature string, event string) (dto.TripayWebhookResponse, error)
	softDeleteTransactionFn func(ctx context.Context, id uuid.UUID) error
}

func (m *mockTransactionService) TripayWebhook(ctx context.Context, rawBody []byte, payload dto.TripayWebhookRequest, callbackSignature string, event string) (dto.TripayWebhookResponse, error) {
	if m.tripayWebhookFn != nil {
		return m.tripayWebhookFn(ctx, rawBody, payload, callbackSignature, event)
	}
	return dto.TripayWebhookResponse{Success: true}, nil
}

func (m *mockTransactionService) SoftDeleteTransaction(ctx context.Context, id uuid.UUID) error {
	if m.softDeleteTransactionFn != nil {
		return m.softDeleteTransactionFn(ctx, id)
	}
	return nil
}

func setupTransactionControllerTestContext(method, target string, body []byte) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest(method, target, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	return c, w
}

func TestTransactionControllerTripayWebhookInvalidJSON(t *testing.T) {
	controller := NewTransactionController(&mockTransactionService{})
	c, w := setupTransactionControllerTestContext(http.MethodPost, "/api/transaction/webhook/tripay", []byte("not-json"))

	controller.TripayWebhook(c)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestTransactionControllerTripayWebhookSuccess(t *testing.T) {
	called := false
	controller := NewTransactionController(&mockTransactionService{
		tripayWebhookFn: func(ctx context.Context, rawBody []byte, payload dto.TripayWebhookRequest, callbackSignature string, event string) (dto.TripayWebhookResponse, error) {
			called = true
			if payload.Reference != "REF-001" {
				t.Fatalf("unexpected reference: %s", payload.Reference)
			}
			if callbackSignature != "sig-1" {
				t.Fatalf("unexpected signature: %s", callbackSignature)
			}
			if event != "payment_status" {
				t.Fatalf("unexpected event: %s", event)
			}
			return dto.TripayWebhookResponse{Success: true}, nil
		},
	})

	body := []byte(`{"reference":"REF-001","is_closed_payment":1,"status":"PAID","total_amount":10000}`)
	c, w := setupTransactionControllerTestContext(http.MethodPost, "/api/transaction/webhook/tripay", body)
	c.Request.Header.Set("X-Callback-Signature", "sig-1")
	c.Request.Header.Set("X-Callback-Event", "payment_status")

	controller.TripayWebhook(c)

	if !called {
		t.Fatal("expected TripayWebhook service to be called")
	}
	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

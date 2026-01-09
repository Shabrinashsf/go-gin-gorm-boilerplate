package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"strings"

	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/dto"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	TransactionService interface {
		TripayWebhook(ctx context.Context, rawBody []byte, payload dto.TripayWebhookRequest, callbackSignature string, event string) (dto.TripayWebhookResponse, error)
		SoftDeleteTransaction(ctx context.Context, id uuid.UUID) error
	}

	transactionService struct {
		transactionRepo repository.TransactionRepository
		db              *gorm.DB
	}
)

func NewTransactionService(transactionRepo repository.TransactionRepository, db *gorm.DB) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
		db:              db,
	}
}

func (s *transactionService) TripayWebhook(ctx context.Context, rawBody []byte, payload dto.TripayWebhookRequest, callbackSignature string, event string) (dto.TripayWebhookResponse, error) {
	privateKey := os.Getenv("TRIPAY_PRIVATE_KEY")

	if event != "payment_status" {
		return dto.TripayWebhookResponse{}, dto.ErrUnrecognizedCallbackEvent
	}

	// hitung signature local
	mac := hmac.New(sha256.New, []byte(privateKey))
	mac.Write(rawBody)
	localSignature := hex.EncodeToString(mac.Sum(nil))

	if localSignature != callbackSignature {
		return dto.TripayWebhookResponse{}, dto.ErrInvalidSignature
	}

	// pastikan closed payment
	if payload.IsClosedPayment != 1 {
		return dto.TripayWebhookResponse{}, dto.ErrOnlyClosedPaymentSupported
	}

	// cari transaksi berdasarkan reference
	transaction, err := s.transactionRepo.GetTransactionByReference(ctx, nil, payload.Reference)
	if err != nil {
		return dto.TripayWebhookResponse{}, dto.ErrTransactionNotFound
	}

	// update status transaksi
	switch strings.ToUpper(payload.Status) {
	case "PAID":
		transaction.Status = "PAID"
		transaction.AmountPaid = payload.TotalAmount
		if err := s.transactionRepo.UpdateTransaction(ctx, nil, transaction); err != nil {
			return dto.TripayWebhookResponse{}, dto.ErrFailedToUpdateStatus
		}
	case "FAILED":
		transaction.Status = "FAILED"
		if err := s.transactionRepo.UpdateTransaction(ctx, nil, transaction); err != nil {
			return dto.TripayWebhookResponse{}, dto.ErrFailedToUpdateStatus
		}
	case "EXPIRED":
		if transaction.Status == "PAID" {
			// jika sudah PAID, jangan diubah ke EXPIRED
			return dto.TripayWebhookResponse{
				Success: true,
			}, nil
		}
		transaction.Status = "EXPIRED"
		if err := s.transactionRepo.UpdateTransaction(ctx, nil, transaction); err != nil {
			return dto.TripayWebhookResponse{}, dto.ErrFailedToUpdateStatus
		}

		if err := s.transactionRepo.SoftDeleteTransaction(ctx, nil, transaction.ID); err != nil {
			return dto.TripayWebhookResponse{}, dto.ErrFailedToSoftDeleteTransaction
		}
	case "REFUND":
		transaction.Status = "REFUND"
		if err := s.transactionRepo.UpdateTransaction(ctx, nil, transaction); err != nil {
			return dto.TripayWebhookResponse{}, dto.ErrFailedToUpdateStatus
		}
	default:
		return dto.TripayWebhookResponse{}, dto.ErrUnknownStatus
	}

	return dto.TripayWebhookResponse{
		Success: true,
	}, nil
}

func (s *transactionService) SoftDeleteTransaction(ctx context.Context, id uuid.UUID) error {
	return s.transactionRepo.SoftDeleteTransaction(ctx, nil, id)
}

package dto

import "errors"

var (
	ErrTransactionNotFound           = errors.New("transaction not found")
	ErrUnrecognizedCallbackEvent     = errors.New("unrecognized callback event")
	ErrInvalidSignature              = errors.New("invalid signature")
	ErrOnlyClosedPaymentSupported    = errors.New("only closed payment supported")
	ErrFailedToUpdateStatus          = errors.New("failed to update transaction status")
	ErrFailedToSoftDeleteTransaction = errors.New("failed to soft delete transaction")
	ErrUnknownStatus                 = errors.New("unknown transaction status")
)

type (
	TripayExpiredTime int

	TripayOrderRequest struct {
		Method        string                    `json:"method"`
		MerchantRef   string                    `json:"merchant_ref"`
		Amount        int                       `json:"amount"`
		CustomerName  string                    `json:"customer_name"`
		CustomerEmail string                    `json:"customer_email"`
		CustomerPhone string                    `json:"customer_phone"`
		OrderItems    []OrderItemPaymentRequest `json:"order_items"`
		ReturnURL     string                    `json:"return_url"`
		ExpiredTime   TripayExpiredTime         `json:"expired_time"`
		Signature     string                    `json:"signature"`
	}

	OrderItemPaymentRequest struct {
		SKU        string `json:"sku"`
		Name       string `json:"name"`
		Price      int    `json:"price"`
		Quantity   int    `json:"quantity"`
		ProductURL string `json:"product_url"`
		ImageURL   string `json:"image_url"`
	}

	TripayResponse struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Data    Data   `json:"data"`
	}

	Data struct {
		Reference   string `json:"reference"`
		MerchantRef string `json:"merchant_ref"`
		PaymentURL  string `json:"checkout_url"`
	}

	TripayWebhookRequest struct {
		Reference         string `json:"reference"`
		MerchantRef       string `json:"merchant_ref"`
		PaymentMethod     string `json:"payment_method"`
		PaymentMethodCode string `json:"payment_method_code"`
		TotalAmount       int    `json:"total_amount"`
		FeeMerchant       int    `json:"fee_merchant"`
		FeeCustomer       int    `json:"fee_customer"`
		TotalFee          int    `json:"total_fee"`
		AmountReceived    int    `json:"amount_received"`
		IsClosedPayment   int    `json:"is_closed_payment"`
		Status            string `json:"status"`
		PaidAt            int    `json:"paid_at"`
	}

	TripayWebhookResponse struct {
		Success bool `json:"success"`
	}
)

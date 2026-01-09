package tripay

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/dto"
)

type Client struct {
	MerchantCode string
	ApiKey       string
	PrivateKey   string
	Mode         string
	signature    Signature
}

func (c *Client) SetSignature(sig Signature) {
	c.signature = sig
}

func (c Client) BaseUrl() string {
	if c.Mode == "development" {
		return "https://tripay.co.id/api-sandbox"
	}
	return "https://tripay.co.id/api"
}

func (c *Client) CreateTransaction(ctx context.Context, req dto.TripayOrderRequest) (dto.TripayResponse, error) {
	if c.signature.MerchanReff == "" {
		return dto.TripayResponse{}, errors.New("signature not set")
	}

	// Tambah signature ke request
	requestBody := map[string]interface{}{
		"method":         req.Method,
		"merchant_ref":   req.MerchantRef,
		"amount":         req.Amount,
		"customer_name":  req.CustomerName,
		"customer_email": req.CustomerEmail,
		"customer_phone": req.CustomerPhone,
		"order_items":    req.OrderItems,
		"expired_time":   int64(req.ExpiredTime),
		"return_url":     req.ReturnURL,
		"signature":      c.signature.CreateSignature(),
	}

	jsonBody, _ := json.Marshal(requestBody)

	url := c.BaseUrl() + "/transaction/create"

	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.ApiKey)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return dto.TripayResponse{}, err
	}

	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)

	var parsed dto.TripayResponse
	if err := json.Unmarshal(bodyBytes, &parsed); err != nil {
		return dto.TripayResponse{}, err
	}

	if !parsed.Success {
		return dto.TripayResponse{}, errors.New(parsed.Message)
	}

	return parsed, nil
}

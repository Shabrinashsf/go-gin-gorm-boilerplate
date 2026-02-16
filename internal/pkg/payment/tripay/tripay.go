package tripay

import (
	"context"
	"os"

	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/dto"
)

func CreateTripayTransaction(ctx context.Context, invoice dto.TripayOrderRequest) (dto.TripayResponse, error) {
	// 1. Create Signature
	sig := Signature{
		Amount:       int64(invoice.Amount),
		PrivateKey:   os.Getenv("TRIPAY_PRIVATE_KEY"),
		MerchantCode: os.Getenv("TRIPAY_MERCHANT_CODE"),
		MerchanReff:  invoice.MerchantRef,
	}

	// 2. Create Client
	client := Client{
		MerchantCode: os.Getenv("TRIPAY_MERCHANT_CODE"),
		ApiKey:       os.Getenv("TRIPAY_API_KEY"),
		PrivateKey:   os.Getenv("TRIPAY_PRIVATE_KEY"),
		Mode:         os.Getenv("APP_ENV"),
	}

	// 3. Set Signature
	client.SetSignature(sig)

	// 4. Call Tripay API
	res, err := client.CreateTransaction(ctx, invoice)
	if err != nil {
		return dto.TripayResponse{}, err
	}

	// Return only checkout URL
	return res, nil
}

package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/internal/api/controller"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/internal/api/service"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type integrationMockUserService struct{}

func (m *integrationMockUserService) RegisterUser(ctx context.Context, req dto.UserRegistrationRequest) (dto.UserResponse, error) {
	return dto.UserResponse{ID: uuid.NewString(), Name: req.Name, Email: req.Email, Role: "user"}, nil
}

func (m *integrationMockUserService) Login(ctx context.Context, req dto.UserLoginRequest) (dto.UserLoginResponse, error) {
	return dto.UserLoginResponse{Token: "dummy-token", Role: "user"}, nil
}

func (m *integrationMockUserService) SendVerificationEmail(ctx context.Context, req dto.SendVerificationEmailRequest) error {
	return nil
}

func (m *integrationMockUserService) VerifyEmail(ctx context.Context, req dto.VerifyEmailRequest) (dto.VerifyEmailResponse, error) {
	return dto.VerifyEmailResponse{Email: "john@example.com", IsVerified: true}, nil
}

func (m *integrationMockUserService) ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error {
	return nil
}

func (m *integrationMockUserService) ResetPassword(ctx context.Context, token string, newPassword string) error {
	return nil
}

func (m *integrationMockUserService) GetUserByID(ctx context.Context, userId uuid.UUID) (dto.UserResponse, error) {
	return dto.UserResponse{ID: userId.String(), Name: "John", Email: "john@example.com", Role: "user"}, nil
}

func (m *integrationMockUserService) UpdateUser(ctx context.Context, userId uuid.UUID, req dto.UserUpdateRequest) (dto.UserResponse, error) {
	name := req.Name
	if name == "" {
		name = "John"
	}

	return dto.UserResponse{ID: userId.String(), Name: name, Email: "john@example.com", Role: "user"}, nil
}

type integrationMockTransactionService struct{}

func (m *integrationMockTransactionService) TripayWebhook(ctx context.Context, rawBody []byte, payload dto.TripayWebhookRequest, callbackSignature string, event string) (dto.TripayWebhookResponse, error) {
	return dto.TripayWebhookResponse{Success: true}, nil
}

func (m *integrationMockTransactionService) SoftDeleteTransaction(ctx context.Context, id uuid.UUID) error {
	return nil
}

type integrationMockJWTService struct{}

func (m *integrationMockJWTService) GenerateToken(userId string, role string) string {
	return ""
}

func (m *integrationMockJWTService) GenerateResetPasswordToken(email string) string {
	return ""
}

func (m *integrationMockJWTService) ValidateToken(token string) (*jwt.Token, error) {
	claims := jwt.MapClaims{
		"user_id": token,
		"role":    "user",
	}

	return &jwt.Token{
		Valid:  true,
		Claims: claims,
	}, nil
}

func (m *integrationMockJWTService) GetUserIDByToken(token string) (string, error) {
	return token, nil
}

func (m *integrationMockJWTService) GetEmailByToken(token string) (string, error) {
	return "john@example.com", nil
}

func (m *integrationMockJWTService) ValidateResetToken(token string) (string, error) {
	return "john@example.com", nil
}

func setupIntegrationRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	userController := controller.NewUserController(&integrationMockUserService{})
	transactionController := controller.NewTransactionController(&integrationMockTransactionService{})
	jwtService := &integrationMockJWTService{}

	User(r, userController, jwtService)
	Transaction(r, transactionController)

	return r
}

func TestRoutesAllEndpoints(t *testing.T) {
	r := setupIntegrationRouter()
	authUserID := uuid.NewString()

	testCases := []struct {
		name       string
		method     string
		target     string
		body       any
		headers    map[string]string
		expectCode int
	}{
		{
			name:   "register user",
			method: http.MethodPost,
			target: "/api/auth",
			body: map[string]any{
				"name":     "John",
				"email":    "john@example.com",
				"password": "secret",
				"instansi": "ACME",
				"no_telp":  "08123",
			},
			expectCode: http.StatusOK,
		},
		{
			name:   "login",
			method: http.MethodPost,
			target: "/api/auth/login",
			body: map[string]any{
				"email":    "john@example.com",
				"password": "secret",
			},
			expectCode: http.StatusOK,
		},
		{
			name:   "send verification email",
			method: http.MethodPost,
			target: "/api/auth/send-verification-email",
			body: map[string]any{
				"email": "john@example.com",
			},
			expectCode: http.StatusOK,
		},
		{
			name:       "verify email",
			method:     http.MethodGet,
			target:     "/api/auth/verify-email?token=sample-token",
			expectCode: http.StatusOK,
		},
		{
			name:   "forgot password",
			method: http.MethodPost,
			target: "/api/auth/forgot-password",
			body: map[string]any{
				"email": "john@example.com",
			},
			expectCode: http.StatusOK,
		},
		{
			name:   "reset password",
			method: http.MethodPost,
			target: "/api/auth/reset-password?token=sample-token",
			body: map[string]any{
				"password": "new-secret",
			},
			expectCode: http.StatusOK,
		},
		{
			name:   "get authenticated user",
			method: http.MethodGet,
			target: "/api/auth/me",
			headers: map[string]string{
				"Authorization": "Bearer " + authUserID,
			},
			expectCode: http.StatusOK,
		},
		{
			name:   "update authenticated user",
			method: http.MethodPatch,
			target: "/api/auth/update",
			body: map[string]any{
				"name": "Updated Name",
			},
			headers: map[string]string{
				"Authorization": "Bearer " + authUserID,
			},
			expectCode: http.StatusOK,
		},
		{
			name:   "tripay webhook",
			method: http.MethodPost,
			target: "/api/transaction/webhook/tripay",
			body: map[string]any{
				"reference":         "REF-001",
				"is_closed_payment": 1,
				"status":            "PAID",
				"total_amount":      10000,
			},
			headers: map[string]string{
				"X-Callback-Signature": "signature",
				"X-Callback-Event":     "payment_status",
			},
			expectCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var bodyBytes []byte
			if tc.body != nil {
				encoded, err := json.Marshal(tc.body)
				if err != nil {
					t.Fatalf("failed to marshal body: %v", err)
				}
				bodyBytes = encoded
			}

			req := httptest.NewRequest(tc.method, tc.target, bytes.NewBuffer(bodyBytes))
			if tc.body != nil {
				req.Header.Set("Content-Type", "application/json")
			}

			for k, v := range tc.headers {
				req.Header.Set(k, v)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectCode {
				t.Fatalf("expected status %d, got %d, body: %s", tc.expectCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestProtectedRouteUnauthorized(t *testing.T) {
	r := setupIntegrationRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

var _ service.JWTService = (*integrationMockJWTService)(nil)

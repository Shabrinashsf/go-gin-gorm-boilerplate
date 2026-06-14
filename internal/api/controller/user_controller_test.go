package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type mockUserService struct {
	registerUserFn          func(ctx context.Context, req dto.UserRegistrationRequest) (dto.UserResponse, error)
	loginFn                 func(ctx context.Context, req dto.UserLoginRequest) (dto.UserLoginResponse, error)
	sendVerificationEmailFn func(ctx context.Context, req dto.SendVerificationEmailRequest) error
	verifyEmailFn           func(ctx context.Context, req dto.VerifyEmailRequest) (dto.VerifyEmailResponse, error)
	forgotPasswordFn        func(ctx context.Context, req dto.ForgotPasswordRequest) error
	resetPasswordFn         func(ctx context.Context, token string, newPassword string) error
	getUserByIDFn           func(ctx context.Context, userId uuid.UUID) (dto.UserResponse, error)
	updateUserFn            func(ctx context.Context, userId uuid.UUID, req dto.UserUpdateRequest) (dto.UserResponse, error)
}

func (m *mockUserService) RegisterUser(ctx context.Context, req dto.UserRegistrationRequest) (dto.UserResponse, error) {
	if m.registerUserFn != nil {
		return m.registerUserFn(ctx, req)
	}
	return dto.UserResponse{}, nil
}

func (m *mockUserService) Login(ctx context.Context, req dto.UserLoginRequest) (dto.UserLoginResponse, error) {
	if m.loginFn != nil {
		return m.loginFn(ctx, req)
	}
	return dto.UserLoginResponse{}, nil
}

func (m *mockUserService) SendVerificationEmail(ctx context.Context, req dto.SendVerificationEmailRequest) error {
	if m.sendVerificationEmailFn != nil {
		return m.sendVerificationEmailFn(ctx, req)
	}
	return nil
}

func (m *mockUserService) VerifyEmail(ctx context.Context, req dto.VerifyEmailRequest) (dto.VerifyEmailResponse, error) {
	if m.verifyEmailFn != nil {
		return m.verifyEmailFn(ctx, req)
	}
	return dto.VerifyEmailResponse{}, nil
}

func (m *mockUserService) ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error {
	if m.forgotPasswordFn != nil {
		return m.forgotPasswordFn(ctx, req)
	}
	return nil
}

func (m *mockUserService) ResetPassword(ctx context.Context, token string, newPassword string) error {
	if m.resetPasswordFn != nil {
		return m.resetPasswordFn(ctx, token, newPassword)
	}
	return nil
}

func (m *mockUserService) GetUserByID(ctx context.Context, userId uuid.UUID) (dto.UserResponse, error) {
	if m.getUserByIDFn != nil {
		return m.getUserByIDFn(ctx, userId)
	}
	return dto.UserResponse{}, nil
}

func (m *mockUserService) UpdateUser(ctx context.Context, userId uuid.UUID, req dto.UserUpdateRequest) (dto.UserResponse, error) {
	if m.updateUserFn != nil {
		return m.updateUserFn(ctx, userId, req)
	}
	return dto.UserResponse{}, nil
}

func setupUserControllerTestContext(method, target string, body []byte) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest(method, target, bytes.NewBuffer(body))
	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}
	c.Request = req

	return c, w
}

func decodeResponseBody(t *testing.T, body []byte) map[string]any {
	t.Helper()

	var got map[string]any
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	return got
}

func TestUserControllerRegisterUserValidationError(t *testing.T) {
	controller := NewUserController(&mockUserService{})
	c, w := setupUserControllerTestContext(http.MethodPost, "/api/auth", []byte(`{}`))

	controller.RegisterUser(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	res := decodeResponseBody(t, w.Body.Bytes())
	if res["success"] != false {
		t.Fatalf("expected success false, got %v", res["success"])
	}
}

func TestUserControllerRegisterUserSuccess(t *testing.T) {
	called := false
	controller := NewUserController(&mockUserService{
		registerUserFn: func(ctx context.Context, req dto.UserRegistrationRequest) (dto.UserResponse, error) {
			called = true
			if req.Email != "john@example.com" {
				t.Fatalf("unexpected email: %s", req.Email)
			}
			return dto.UserResponse{ID: "u1", Email: req.Email, Name: req.Name}, nil
		},
	})

	body := []byte(`{"name":"John","email":"john@example.com","password":"secret","instansi":"ACME","no_telp":"08123"}`)
	c, w := setupUserControllerTestContext(http.MethodPost, "/api/auth", body)

	controller.RegisterUser(c)

	if !called {
		t.Fatal("expected RegisterUser service to be called")
	}
	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUserControllerLoginSuccess(t *testing.T) {
	called := false
	controller := NewUserController(&mockUserService{
		loginFn: func(ctx context.Context, req dto.UserLoginRequest) (dto.UserLoginResponse, error) {
			called = true
			return dto.UserLoginResponse{Token: "token", Role: "user"}, nil
		},
	})

	body := []byte(`{"email":"john@example.com","password":"secret"}`)
	c, w := setupUserControllerTestContext(http.MethodPost, "/api/auth/login", body)

	controller.Login(c)

	if !called {
		t.Fatal("expected Login service to be called")
	}
	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUserControllerSendVerificationEmailSuccess(t *testing.T) {
	called := false
	controller := NewUserController(&mockUserService{
		sendVerificationEmailFn: func(ctx context.Context, req dto.SendVerificationEmailRequest) error {
			called = true
			if req.Email != "john@example.com" {
				t.Fatalf("unexpected email: %s", req.Email)
			}
			return nil
		},
	})

	body := []byte(`{"email":"john@example.com"}`)
	c, w := setupUserControllerTestContext(http.MethodPost, "/api/auth/send-verification-email", body)

	controller.SendVerificationEmail(c)

	if !called {
		t.Fatal("expected SendVerificationEmail service to be called")
	}
	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUserControllerVerifyEmailMissingToken(t *testing.T) {
	controller := NewUserController(&mockUserService{})
	c, w := setupUserControllerTestContext(http.MethodGet, "/api/auth/verify-email", nil)

	controller.VerifyEmail(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestUserControllerVerifyEmailSuccess(t *testing.T) {
	called := false
	controller := NewUserController(&mockUserService{
		verifyEmailFn: func(ctx context.Context, req dto.VerifyEmailRequest) (dto.VerifyEmailResponse, error) {
			called = true
			if req.Token != "abc-token" {
				t.Fatalf("unexpected token: %s", req.Token)
			}
			return dto.VerifyEmailResponse{Email: "john@example.com", IsVerified: true}, nil
		},
	})

	c, w := setupUserControllerTestContext(http.MethodGet, "/api/auth/verify-email?token=abc-token", nil)

	controller.VerifyEmail(c)

	if !called {
		t.Fatal("expected VerifyEmail service to be called")
	}
	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUserControllerForgotPasswordSuccess(t *testing.T) {
	called := false
	controller := NewUserController(&mockUserService{
		forgotPasswordFn: func(ctx context.Context, req dto.ForgotPasswordRequest) error {
			called = true
			return nil
		},
	})

	body := []byte(`{"email":"john@example.com"}`)
	c, w := setupUserControllerTestContext(http.MethodPost, "/api/auth/forgot-password", body)

	controller.ForgotPassword(c)

	if !called {
		t.Fatal("expected ForgotPassword service to be called")
	}
	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUserControllerResetPasswordSuccess(t *testing.T) {
	called := false
	controller := NewUserController(&mockUserService{
		resetPasswordFn: func(ctx context.Context, token string, newPassword string) error {
			called = true
			if token != "reset-token" {
				t.Fatalf("unexpected token: %s", token)
			}
			if newPassword != "new-secret" {
				t.Fatalf("unexpected password")
			}
			return nil
		},
	})

	body := []byte(`{"password":"new-secret"}`)
	c, w := setupUserControllerTestContext(http.MethodPost, "/api/auth/reset-password?token=reset-token", body)

	controller.ResetPassword(c)

	if !called {
		t.Fatal("expected ResetPassword service to be called")
	}
	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUserControllerMeAuthSuccess(t *testing.T) {
	userID := uuid.New()
	called := false
	controller := NewUserController(&mockUserService{
		getUserByIDFn: func(ctx context.Context, id uuid.UUID) (dto.UserResponse, error) {
			called = true
			if id != userID {
				t.Fatalf("unexpected user id: %s", id)
			}
			return dto.UserResponse{ID: id.String(), Email: "john@example.com"}, nil
		},
	})

	c, w := setupUserControllerTestContext(http.MethodGet, "/api/auth/me", nil)
	c.Set("user_id", userID.String())

	controller.MeAuth(c)

	if !called {
		t.Fatal("expected GetUserByID service to be called")
	}
	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUserControllerUpdateUserSuccess(t *testing.T) {
	userID := uuid.New()
	called := false
	controller := NewUserController(&mockUserService{
		updateUserFn: func(ctx context.Context, id uuid.UUID, req dto.UserUpdateRequest) (dto.UserResponse, error) {
			called = true
			if id != userID {
				t.Fatalf("unexpected user id: %s", id)
			}
			if req.Name != "Updated Name" {
				t.Fatalf("unexpected name: %s", req.Name)
			}
			return dto.UserResponse{ID: id.String(), Name: req.Name}, nil
		},
	})

	body := []byte(`{"name":"Updated Name"}`)
	c, w := setupUserControllerTestContext(http.MethodPatch, "/api/auth/update", body)
	c.Set("user_id", userID.String())

	controller.UpdateUser(c)

	if !called {
		t.Fatal("expected UpdateUser service to be called")
	}
	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

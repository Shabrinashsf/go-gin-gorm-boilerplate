package service

import (
	"context"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/dto"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/entity"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/helpers"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/repository"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/utils"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/utils/mailer"
	"github.com/google/uuid"
)

type (
	UserService interface {
		RegisterUser(ctx context.Context, req dto.UserRegistrationRequest) (dto.UserResponse, error)
		Login(ctx context.Context, req dto.UserLoginRequest) (dto.UserLoginResponse, error)
		SendVerificationEmail(ctx context.Context, req dto.SendVerificationEmailRequest) error
		VerifyEmail(ctx context.Context, req dto.VerifyEmailRequest) (dto.VerifyEmailResponse, error)
		ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error
		ResetPassword(ctx context.Context, token string, newPassword string) error
		GetUserByID(ctx context.Context, userId uuid.UUID) (dto.UserResponse, error)
		UpdateUser(ctx context.Context, userId uuid.UUID, req dto.UserUpdateRequest) (dto.UserResponse, error)
	}

	userService struct {
		userRepository repository.UserRepository
		jwtService     JWTService
		mailer         mailer.Mailer
	}
)

func NewUserController(ur repository.UserRepository, jwt JWTService, mailer mailer.Mailer) UserService {
	return &userService{
		userRepository: ur,
		jwtService:     jwt,
		mailer:         mailer,
	}
}

var (
	mu sync.Mutex

	VERIFY_EMAIL_TEMPLATE = "utils/mailer/template/verification_email.html"
	VERIFY_EMAIL_PATH     = "verify-email"
	FORGET_EMAIL_TEMPLATE = "utils/mailer/template/forgot_password_email.html"
	FORGET_EMAIL_PATH     = "reset-password"
)

func (s *userService) RegisterUser(ctx context.Context, req dto.UserRegistrationRequest) (dto.UserResponse, error) {
	mu.Lock()
	defer mu.Unlock()

	_, flag, _ := s.userRepository.GetUserByEmail(ctx, nil, req.Email)
	if flag {
		return dto.UserResponse{}, dto.ErrorEmailAlreadyExists
	}

	user := entity.User{
		Name:       req.Name,
		Email:      req.Email,
		Password:   req.Password,
		Instansi:   req.Instansi,
		NoTelp:     req.NoTelp,
		Role:       entity.RoleUser,
		IsVerified: false,
	}

	expired := time.Now().Add(time.Hour * 24).Format("2006-01-02 15:04:05")
	plainText := user.Email + "_" + expired
	token, err := utils.AESEncrypt(plainText)
	if err != nil {
		return dto.UserResponse{}, err
	}

	verifyLink := os.Getenv("APP_URL") + "/" + VERIFY_EMAIL_PATH + "?token=" + token
	data := map[string]any{
		"Email": user.Email,
		"Link":  verifyLink,
	}

	mail := s.mailer.MakeMail(VERIFY_EMAIL_TEMPLATE, data)
	if mail.Error != nil {
		return dto.UserResponse{}, dto.ErrMakeMail
	}

	if err := mail.SendEmail(user.Email, "Backend Boilerplate - Verification Email").Error; err != nil {
		return dto.UserResponse{}, dto.ErrSendMail
	}

	newUser, err := s.userRepository.RegisterUser(ctx, nil, user)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:         newUser.ID.String(),
		Name:       newUser.Name,
		Email:      newUser.Email,
		Instansi:   newUser.Instansi,
		NoTelp:     newUser.NoTelp,
		Role:       string(newUser.Role),
		IsVerified: newUser.IsVerified,
	}, nil
}

func (s *userService) Login(ctx context.Context, req dto.UserLoginRequest) (dto.UserLoginResponse, error) {
	user, flag, err := s.userRepository.GetUserByEmail(ctx, nil, req.Email)
	if err != nil || !flag {
		return dto.UserLoginResponse{}, dto.ErrUserNotFound
	}

	if !user.IsVerified {
		return dto.UserLoginResponse{}, dto.ErrAccountAlreadyVerified
	}

	checkPassword, err := helpers.CheckPassword(user.Password, []byte(req.Password))
	if err != nil || !checkPassword {
		return dto.UserLoginResponse{}, dto.ErrPasswordNotMatch
	}

	token := s.jwtService.GenerateToken(user.ID.String(), string(user.Role))

	return dto.UserLoginResponse{
		Token: token,
		Role:  string(user.Role),
	}, nil
}

func (s *userService) SendVerificationEmail(ctx context.Context, req dto.SendVerificationEmailRequest) error {
	user, _, err := s.userRepository.GetUserByEmail(ctx, nil, req.Email)
	if err != nil {
		return err
	}

	expired := time.Now().Add(time.Hour * 24).Format("2006-01-02 15:04:05")
	plainText := user.Email + "_" + expired
	token, err := utils.AESEncrypt(plainText)
	if err != nil {
		return err
	}

	verifyLink := os.Getenv("APP_URL") + "/" + VERIFY_EMAIL_PATH + "?token=" + token
	data := map[string]any{
		"Email":  user.Email,
		"Verify": verifyLink,
	}

	mail := s.mailer.MakeMail(VERIFY_EMAIL_TEMPLATE, data)
	if mail.Error != nil {
		return dto.ErrMakeMail
	}

	if err := mail.SendEmail(user.Email, "Backend Boilerplate - Verification Email").Error; err != nil {
		return dto.ErrSendMail
	}

	return nil
}

func (s *userService) VerifyEmail(ctx context.Context, req dto.VerifyEmailRequest) (dto.VerifyEmailResponse, error) {
	decryptedToken, err := utils.AESDecrypt(req.Token)
	if err != nil {
		return dto.VerifyEmailResponse{}, dto.ErrTokenInvalid
	}

	if !strings.Contains(decryptedToken, "_") {
		return dto.VerifyEmailResponse{}, dto.ErrTokenInvalid
	}

	decryptedTokenSplit := strings.Split(decryptedToken, "_")
	email := decryptedTokenSplit[0]
	expired := decryptedTokenSplit[1]

	now := time.Now()
	expiredTime, err := time.Parse("2006-01-02 15:04:05", expired)
	if err != nil {
		return dto.VerifyEmailResponse{}, dto.ErrTokenInvalid
	}

	if expiredTime.Sub(now) < 0 {
		return dto.VerifyEmailResponse{
			Email:      email,
			IsVerified: false,
		}, dto.ErrTokenExpired
	}

	user, _, err := s.userRepository.GetUserByEmail(ctx, nil, email)
	if err != nil {
		return dto.VerifyEmailResponse{}, dto.ErrUserNotFound
	}

	if user.IsVerified {
		return dto.VerifyEmailResponse{}, dto.ErrAccountAlreadyVerified
	}

	userUpdate := entity.User{
		ID:         user.ID,
		IsVerified: true,
	}

	updatedUser, err := s.userRepository.UpdateUser(ctx, nil, userUpdate)
	if err != nil {
		return dto.VerifyEmailResponse{}, dto.ErrUpdateUser
	}

	return dto.VerifyEmailResponse{
		Email:      email,
		IsVerified: updatedUser.IsVerified,
	}, nil
}

func (s *userService) ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error {
	user, _, err := s.userRepository.GetUserByEmail(ctx, nil, req.Email)
	if err != nil {
		return dto.ErrEmailNotFound
	}

	expired := time.Now().Add(time.Hour * 24).Format("2006-01-02 15:04:05")
	plainText := user.Email + "_" + expired
	token, err := utils.AESEncrypt(plainText)
	if err != nil {
		return err
	}

	verifyLink := os.Getenv("APP_URL") + "/" + FORGET_EMAIL_PATH + "?token=" + token
	data := map[string]any{
		"Email":  user.Email,
		"Verify": verifyLink,
	}

	mail := s.mailer.MakeMail(FORGET_EMAIL_TEMPLATE, data)
	if mail.Error != nil {
		return dto.ErrMakeMail
	}

	if err := mail.SendEmail(user.Email, "Schematics 2025 - Reset Password").Error; err != nil {
		return dto.ErrSendMail
	}

	return nil
}

func (s *userService) ResetPassword(ctx context.Context, token string, newPassword string) error {
	decryptedToken, err := utils.AESDecrypt(token)
	if err != nil {
		return dto.ErrTokenInvalid
	}

	tokenParts := strings.Split(decryptedToken, "_")
	if len(tokenParts) < 2 {
		return dto.ErrTokenInvalid
	}

	email := tokenParts[0]
	expirationDate := tokenParts[1]
	expirationTime, err := time.Parse("2006-01-02 15:04:05", expirationDate)

	if err != nil {
		return dto.ErrTokenInvalid
	}

	if time.Now().After(expirationTime) {
		return dto.ErrTokenExpired
	}
	hashedPassword, err := helpers.HashPassword(newPassword)

	if err != nil {
		return dto.ErrHashPasswordFailed
	}

	err = s.userRepository.ResetPassword(ctx, email, hashedPassword)
	if err != nil {
		return dto.ErrUpdateUser
	}

	return nil
}

func (s *userService) GetUserByID(ctx context.Context, userId uuid.UUID) (dto.UserResponse, error) {
	user, err := s.userRepository.GetUserByID(ctx, nil, userId)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		Name:     user.Name,
		Email:    user.Email,
		Instansi: user.Instansi,
		NoTelp:   user.NoTelp,
	}, nil
}

func (s *userService) UpdateUser(ctx context.Context, userId uuid.UUID, req dto.UserUpdateRequest) (dto.UserResponse, error) {
	user, err := s.userRepository.GetUserByID(ctx, nil, userId)
	if err != nil {
		return dto.UserResponse{}, dto.ErrUserNotFound
	}

	userUpdate := entity.User{
		ID:       user.ID,
		Name:     req.Name,
		Email:    user.Email,
		Instansi: req.Instansi,
		NoTelp:   req.NoTelp,
	}

	userUpdate, err = s.userRepository.UpdateUser(ctx, nil, userUpdate)
	if err != nil {
		return dto.UserResponse{}, dto.ErrUpdateUser
	}

	return dto.UserResponse{
		ID:       userUpdate.ID.String(),
		Name:     userUpdate.Name,
		Email:    userUpdate.Email,
		Instansi: userUpdate.Instansi,
		NoTelp:   userUpdate.NoTelp,
	}, nil
}

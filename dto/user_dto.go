package dto

import "errors"

const (
	// Failed
	MESSAGE_FAILED_REGISTER_USER   = "gagal melakukan registrasi user"
	MESSAGE_FAILED_LOGIN_USER      = "gagal melakukan login user"
	MESSAGE_FAILED_VERIFY_EMAIL    = "gagal memverifikasi email"
	MESSAGE_FAILED_FORGET_PASSWORD = "gagal memproses permintaan lupa password"
	MESSAGE_FAILED_RESET_PASSWORD  = "gagal mereset password"
	MESSAGE_FAILED_GET_USER        = "gagal mendapatkan data user"
	MESSAGE_FAILED_UPDATE_USER     = "gagal memperbarui data user"

	// Success
	MESSAGE_SUCCESS_REGISTER_USER           = "berhasil melakukan registrasi user"
	MESSAGE_SUCCESS_LOGIN_USER              = "berhasil melakukan login user"
	MESSAGE_SEND_VERIFICATION_EMAIL_SUCCESS = "berhasil mengirim email verifikasi"
	MESSAGE_SUCCESS_VERIFY_EMAIL            = "berhasil memverifikasi email"
	MESSAGE_SUCCESS_FORGET_PASSWORD         = "berhasil memproses permintaan lupa password"
	MESSAGE_SUCCESS_RESET_PASSWORD          = "berhasil mereset password"
	MESSAGE_SUCCESS_GET_USER                = "berhasil mendapatkan data user"
	MESSAGE_SUCCESS_UPDATE_USER             = "berhasil memperbarui data user"
)

var (
	ErrorEmailAlreadyExists   = errors.New("email sudah terdaftar")
	ErrMakeMail               = errors.New("gagal membuat email")
	ErrSendMail               = errors.New("gagal mengirim email")
	ErrTokenInvalid           = errors.New("token tidak valid atau kadaluarsa")
	ErrTokenExpired           = errors.New("token telah kadaluarsa")
	ErrUserNotFound           = errors.New("user tidak ditemukan")
	ErrAccountAlreadyVerified = errors.New("akun sudah terverifikasi")
	ErrUpdateUser             = errors.New("gagal memperbarui data user")
	ErrPasswordNotMatch       = errors.New("password tidak sesuai")
	ErrEmailNotFound          = errors.New("email tidak ditemukan")
	ErrHashPasswordFailed     = errors.New("gagal melakukan hash password")
	ErrNoChanges              = errors.New("tidak ada perubahan pada data user")
	ErrInvalidCredentials     = errors.New("kredensial tidak valid")
)

type (
	UserRegistrationRequest struct {
		Name     string `json:"name" form:"name" binding:"required"`
		Email    string `json:"email" form:"email" binding:"required,email"`
		Password string `json:"password" form:"password" binding:"required"`
		Instansi string `json:"instansi" form:"instansi" binding:"required"`
		NoTelp   string `json:"no_telp" form:"no_telp" binding:"required"`
	}

	UserResponse struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Email      string `json:"email"`
		Instansi   string `json:"instansi"`
		NoTelp     string `json:"no_telp"`
		Role       string `json:"role"`
		IsVerified bool   `json:"is_verified"`
	}

	UserLoginRequest struct {
		Email    string `json:"email" form:"email" binding:"required,email"`
		Password string `json:"password" form:"password" binding:"required"`
	}

	UserLoginResponse struct {
		Token string `json:"token"`
		Role  string `json:"role"`
	}

	SendVerificationEmailRequest struct {
		Email string `json:"email" form:"email" binding:"required,email"`
	}

	VerifyEmailRequest struct {
		Token string `json:"token" form:"token" binding:"required"`
	}

	VerifyEmailResponse struct {
		Email      string `json:"email"`
		IsVerified bool   `json:"is_verified"`
	}

	ForgotPasswordRequest struct {
		Email string `json:"email" form:"email" binding:"required,email"`
	}

	ResetPasswordRequest struct {
		Password string `json:"password" form:"password" binding:"required"`
	}

	ResetPasswordResponse struct {
		Email string `json:"email"`
	}

	UserUpdateRequest struct {
		Name     string `json:"name" form:"name"`
		Instansi string `json:"instansi" form:"instansi"`
		NoTelp   string `json:"no_telp" form:"no_telp"`
	}
)

package config

import (
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/internal/api/controller"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/internal/api/repository"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/internal/api/service"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/internal/pkg/mailer"
	"github.com/samber/do"
	"gorm.io/gorm"
)

// Provider is the dependency injection container
type Provider struct {
	injector *do.Injector
}

// NewProvider creates a new DI provider with database connection
func NewProvider(db *gorm.DB) *Provider {
	injector := do.New()

	// Register database singleton
	do.Provide(injector, func(i *do.Injector) (*gorm.DB, error) {
		return db, nil
	})

	// =========== SERVICES ===========
	// JWT Service
	do.Provide(injector, func(i *do.Injector) (service.JWTService, error) {
		return service.NewJWTService(), nil
	})

	// Mailer Service
	do.Provide(injector, func(i *do.Injector) (mailer.Mailer, error) {
		return mailer.NewMailer(), nil
	})

	// =========== REPOSITORIES ===========
	// User Repository
	do.Provide(injector, func(i *do.Injector) (repository.UserRepository, error) {
		db := do.MustInvoke[*gorm.DB](i)
		return repository.NewUserController(db), nil
	})

	// Transaction Repository
	do.Provide(injector, func(i *do.Injector) (repository.TransactionRepository, error) {
		db := do.MustInvoke[*gorm.DB](i)
		return repository.NewTransactionRepository(db), nil
	})

	// =========== SERVICES ===========
	// User Service
	do.Provide(injector, func(i *do.Injector) (service.UserService, error) {
		userRepo := do.MustInvoke[repository.UserRepository](i)
		jwtService := do.MustInvoke[service.JWTService](i)
		mailerService := do.MustInvoke[mailer.Mailer](i)
		db := do.MustInvoke[*gorm.DB](i)
		return service.NewUserService(userRepo, jwtService, mailerService, db), nil
	})

	// Transaction Service
	do.Provide(injector, func(i *do.Injector) (service.TransactionService, error) {
		transactionRepo := do.MustInvoke[repository.TransactionRepository](i)
		db := do.MustInvoke[*gorm.DB](i)
		return service.NewTransactionService(transactionRepo, db), nil
	})

	// =========== CONTROLLERS ===========
	// User Controller
	do.Provide(injector, func(i *do.Injector) (controller.UserController, error) {
		userService := do.MustInvoke[service.UserService](i)
		return controller.NewUserController(userService), nil
	})

	// Transaction Controller
	do.Provide(injector, func(i *do.Injector) (controller.TransactionController, error) {
		transactionService := do.MustInvoke[service.TransactionService](i)
		return controller.NewTransactionController(transactionService), nil
	})

	return &Provider{
		injector: injector,
	}
}

// =========== INVOKE METHODS ===========

// InvokeJWTService returns the JWT service instance
func (p *Provider) InvokeJWTService() service.JWTService {
	return do.MustInvoke[service.JWTService](p.injector)
}

// InvokeMailerService returns the Mailer service instance
func (p *Provider) InvokeMailerService() mailer.Mailer {
	return do.MustInvoke[mailer.Mailer](p.injector)
}

// InvokeUserController returns the User controller instance
func (p *Provider) InvokeUserController() controller.UserController {
	return do.MustInvoke[controller.UserController](p.injector)
}

// InvokeTransactionController returns the Transaction controller instance
func (p *Provider) InvokeTransactionController() controller.TransactionController {
	return do.MustInvoke[controller.TransactionController](p.injector)
}

// InvokeUserService returns the User service instance
func (p *Provider) InvokeUserService() service.UserService {
	return do.MustInvoke[service.UserService](p.injector)
}

// InvokeTransactionService returns the Transaction service instance
func (p *Provider) InvokeTransactionService() service.TransactionService {
	return do.MustInvoke[service.TransactionService](p.injector)
}

// InvokeUserRepository returns the User repository instance
func (p *Provider) InvokeUserRepository() repository.UserRepository {
	return do.MustInvoke[repository.UserRepository](p.injector)
}

// InvokeTransactionRepository returns the Transaction repository instance
func (p *Provider) InvokeTransactionRepository() repository.TransactionRepository {
	return do.MustInvoke[repository.TransactionRepository](p.injector)
}

// InvokeDatabase returns the database instance
func (p *Provider) InvokeDatabase() *gorm.DB {
	return do.MustInvoke[*gorm.DB](p.injector)
}

// Shutdown gracefully shuts down the provider and cleans up resources
func (p *Provider) Shutdown() {
	// The do library handles singleton lifecycle automatically
	// Additional cleanup can be added here if needed
}

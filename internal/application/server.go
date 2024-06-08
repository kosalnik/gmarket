package application

import (
	"context"
	"database/sql"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/kosalnik/gmarket/internal/accrual"
	"github.com/kosalnik/gmarket/internal/config"
	"github.com/kosalnik/gmarket/internal/infra/auth"
	"github.com/kosalnik/gmarket/internal/infra/crypt"
	"github.com/kosalnik/gmarket/internal/infra/logger"
	"github.com/kosalnik/gmarket/internal/infra/postgres"
	"github.com/kosalnik/gmarket/pkg/domain"
	"github.com/kosalnik/gmarket/pkg/domain/service"
)

type Application struct {
	cfg *config.Config
	db  *sql.DB

	repo domain.Repository

	passwordHasher service.PasswordHasher
	userService    domain.UserService
	authService    auth.TokenEncoder
	orderService   domain.OrderService
	accrualService service.AccrualService
}

func New(cfg *config.Config) *Application {
	return &Application{cfg: cfg}
}

func (app *Application) Run(ctx context.Context) (err error) {
	app.InitServices(ctx)

	logger.Info("Listen " + app.cfg.Server.Address)

	return http.ListenAndServe(app.cfg.Server.Address, app.GetRoutes(ctx))
}

func (app *Application) InitServices(ctx context.Context) {
	var err error
	if app.db, err = postgres.NewDB(ctx, app.cfg.Database); err != nil {
		panic(err)
	}

	if app.authService, err = auth.NewJwtEncoder(app.cfg.JWT); err != nil {
		panic(err)
	}
	app.passwordHasher = crypt.NewPasswordHasher()
	if app.repo, err = postgres.NewRepository(app.db); err != nil {
		panic(err)
	}
	if app.userService, err = service.NewUserService(app.repo, app.passwordHasher); err != nil {
		panic(err)
	}
	if app.accrualService, err = accrual.NewAccrual(app.cfg.AccrualSystem); err != nil {
		panic(err)
	}
	if app.orderService, err = service.NewOrderService(app.repo, app.accrualService); err != nil {
		panic(err)
	}
}

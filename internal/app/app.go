package app

import (
	dbmigrations "application-service/database/migrations"
	"application-service/internal/config"
	deliveryhttp "application-service/internal/delivery/http"
	"application-service/internal/delivery/http/server"
	"application-service/internal/repository/postgres"
	"application-service/internal/usecase/application"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type App struct {
	container *Container
}

type Container struct {
	Config     *config.Config
	Logger     *zap.Logger
	PG         *pgxpool.Pool
	HTTPServer *http.Server
}

// New constructs the application container and prepares all dependencies.
func New(ctx context.Context, cfg *config.Config, logger *zap.Logger) (*App, error) {
	container, err := BuildContainer(ctx, cfg, logger)
	if err != nil {
		return nil, err
	}

	return &App{container: container}, nil
}

// Run launches background workers and the HTTP server until the context is cancelled.
func (a *App) Run(ctx context.Context) error {
	//a.container.CleanupService.Start(ctx)

	srvErrors := make(chan error, 1)
	go func() {
		a.container.Logger.Info("http: starting", zap.String("addr", a.container.Config.HTTPAddr))
		if err := a.container.HTTPServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			srvErrors <- err
		}
	}()

	select {
	case err := <-srvErrors:
		return a.shutdownWithError(fmt.Errorf("http server: %w", err))
	case <-ctx.Done():
	}

	return a.shutdownWithError(nil)
}

// runMigrations applies all pending up migrations to the database.
func runMigrations(databaseURL, migrationsTable string) error {
	src, err := iofs.New(dbmigrations.FS, ".")
	if err != nil {
		return fmt.Errorf("iofs source: %w", err)
	}

	migrateURL := strings.NewReplacer(
		"postgres://", "pgx5://",
		"postgresql://", "pgx5://",
	).Replace(databaseURL)

	if strings.Contains(migrateURL, "?") {
		migrateURL += "&x-migrations-table=" + migrationsTable
	} else {
		migrateURL += "?x-migrations-table=" + migrationsTable
	}

	m, err := migrate.NewWithSourceInstance("iofs", src, migrateURL)
	if err != nil {
		return fmt.Errorf("migrator: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate up: %w", err)
	}
	return nil
}

// BuildContainer initializes infrastructure, application services, and the HTTP server.
func BuildContainer(ctx context.Context, cfg *config.Config, logger *zap.Logger) (*Container, error) {
	if err := runMigrations(cfg.DatabaseURL, "schema_migrations_admin"); err != nil {
		return nil, fmt.Errorf("migrations: %w", err)
	}

	pgPool, err := newPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	repositories := buildRepositories(pgPool, cfg)

	applicationSvc := application.NewService(
		repositories.applicationRepo,
		logger,
	)

	applicationHandler := deliveryhttp.NewApplicationHandler(applicationSvc, logger)
	integrationHandler := deliveryhttp.NewIntegrationHandler(applicationSvc, logger)
	router := server.NewRouter(cfg, logger, applicationHandler, integrationHandler)

	httpSrv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	return &Container{
		Config: cfg,
		Logger: logger,
		PG:     pgPool,
		//Redis:          redisClient,
		HTTPServer: httpSrv,
		//CleanupService: cleanupSvc,
	}, nil
}

func newPostgresPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	pgPool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pgPool.Ping(pingCtx); err != nil {
		pgPool.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return pgPool, nil
}

type repositories struct {
	applicationRepo *postgres.ApplicationRepo
}

func buildRepositories(pg *pgxpool.Pool, cfg *config.Config) repositories {
	return repositories{
		applicationRepo: postgres.NewApplicationRepo(pg),
	}
}

func (a *App) shutdownWithError(err error) error {
	//a.container.CleanupService.Stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if shutErr := a.container.HTTPServer.Shutdown(shutdownCtx); shutErr != nil {
		a.container.Logger.Warn("http: shutdown", zap.Error(shutErr))
		if err == nil {
			err = shutErr
		}
	}

	a.container.Close()
	return err
}

// Close releases pooled resources.
func (c *Container) Close() {
	//if c.Redis != nil {
	//	_ = c.Redis.Close()
	//}
	if c.PG != nil {
		c.PG.Close()
	}
}

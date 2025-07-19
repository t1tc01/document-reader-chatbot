package database

import (
	"context"
	"database/sql"
	"fmt"

	"document-reader-chatbot/configs"

	_ "github.com/lib/pq"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database wraps the PostgreSQL client with tracing
type Database struct {
	DB           *gorm.DB
	SqlDB        *sql.DB
	AtlasManager *AtlasManager
	tracer       trace.Tracer
}

// Connect establishes a connection to PostgreSQL with proper configuration
func Connect(cfg configs.DatabaseConfig) (*Database, error) {
	tracer := otel.Tracer("database")
	ctx, span := tracer.Start(context.Background(), "database.connect")
	defer span.End()

	span.SetAttributes(
		attribute.String("database.name", cfg.Name),
		attribute.String("database.system", "postgresql"),
	)

	// Configure GORM logger
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Connect to PostgreSQL using GORM
	db, err := gorm.Open(postgres.Open(cfg.URL), gormConfig)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Get underlying sql.DB for connection pool configuration
	sqlDB, err := db.DB()
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.MaxLifetime)

	// Test the connection
	if err := sqlDB.PingContext(ctx); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	return &Database{
		DB:           db,
		SqlDB:        sqlDB,
		AtlasManager: NewAtlasManager(cfg),
		tracer:       tracer,
	}, nil
}

// Disconnect closes the database connection
func (d *Database) Disconnect(ctx context.Context) error {
	_, span := d.tracer.Start(ctx, "database.disconnect")
	defer span.End()

	if err := d.SqlDB.Close(); err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to disconnect from PostgreSQL: %w", err)
	}

	return nil
}

// WithTransaction executes a function within a PostgreSQL transaction
func (d *Database) WithTransaction(ctx context.Context, fn func(*gorm.DB) error) error {
	_, span := d.tracer.Start(ctx, "database.transaction")
	defer span.End()

	tx := d.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		span.RecordError(tx.Error)
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		span.RecordError(err)
		return fmt.Errorf("transaction failed: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

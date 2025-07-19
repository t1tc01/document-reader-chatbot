package database

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"document-reader-chatbot/configs"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// AtlasManager handles Atlas migration operations
type AtlasManager struct {
	config configs.DatabaseConfig
	tracer trace.Tracer
}

// NewAtlasManager creates a new Atlas migration manager
func NewAtlasManager(cfg configs.DatabaseConfig) *AtlasManager {
	return &AtlasManager{
		config: cfg,
		tracer: otel.Tracer("atlas-manager"),
	}
}

// ApplyMigrations applies pending migrations using Atlas
func (am *AtlasManager) ApplyMigrations(ctx context.Context, env string) error {
	ctx, span := am.tracer.Start(ctx, "atlas.apply_migrations")
	defer span.End()

	span.SetAttributes(
		attribute.String("atlas.environment", env),
		attribute.String("database.system", "postgresql"),
	)

	cmd := exec.CommandContext(ctx, "atlas", "migrate", "apply", "--env", env)
	cmd.Dir = am.getProjectRoot()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to apply Atlas migrations: %w", err)
	}

	return nil
}

// GenerateMigration creates a new migration based on schema changes
func (am *AtlasManager) GenerateMigration(ctx context.Context, name string, env string) error {
	ctx, span := am.tracer.Start(ctx, "atlas.generate_migration")
	defer span.End()

	span.SetAttributes(
		attribute.String("atlas.migration_name", name),
		attribute.String("atlas.environment", env),
	)

	cmd := exec.CommandContext(ctx, "atlas", "migrate", "diff", name, "--env", env)
	cmd.Dir = am.getProjectRoot()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to generate Atlas migration: %w", err)
	}

	return nil
}

// ValidateMigrations validates all migration files
func (am *AtlasManager) ValidateMigrations(ctx context.Context) error {
	ctx, span := am.tracer.Start(ctx, "atlas.validate_migrations")
	defer span.End()

	cmd := exec.CommandContext(ctx, "atlas", "migrate", "validate", "--dir", "file://migrations")
	cmd.Dir = am.getProjectRoot()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to validate Atlas migrations: %w", err)
	}

	return nil
}

// GetMigrationStatus returns the current migration status
func (am *AtlasManager) GetMigrationStatus(ctx context.Context, env string) error {
	ctx, span := am.tracer.Start(ctx, "atlas.migration_status")
	defer span.End()

	span.SetAttributes(attribute.String("atlas.environment", env))

	cmd := exec.CommandContext(ctx, "atlas", "migrate", "status", "--env", env)
	cmd.Dir = am.getProjectRoot()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to get Atlas migration status: %w", err)
	}

	return nil
}

// IsAtlasInstalled checks if Atlas CLI is installed
func (am *AtlasManager) IsAtlasInstalled() bool {
	_, err := exec.LookPath("atlas")
	return err == nil
}

// getProjectRoot returns the project root directory
func (am *AtlasManager) getProjectRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "atlas.hcl")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "."
}

package user

import (
	"context"
	"fmt"
	"strings"
	"time"

	"document-reader-chatbot/pkg/database"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

// Repository defines the interface for user data operations
type Repository interface {
	// User
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetUserByEmailOrUsername(ctx context.Context, emailOrUsername string) (*User, error)
	Update(ctx context.Context, id uuid.UUID, update *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, query ListUsersQuery) ([]*User, int64, error)
	UpdateLastLogin(ctx context.Context, userId uuid.UUID, passwordHashed string) error

	//Role
	GetRoleByName(ctx context.Context, name string) (*Role, error)
	CreateUserRole(ctx context.Context, userRole *UserRole) (*UserRole, error)
	CreateSession(ctx context.Context, session *UserSession) (*UserSession, error)
	CreateToken(ctx context.Context, token *UserToken) (*UserToken, error)
}

// repository implements the Repository interface
type repository struct {
	db     *database.Database
	tracer trace.Tracer
}

// CreateSession implements Repository.
func (r *repository) CreateSession(ctx context.Context, session *UserSession) (*UserSession, error) {
	ctx, span := r.tracer.Start(ctx, "user.repository.create_session")
	defer span.End()

	span.SetAttributes(attribute.String("user_session.user_id", session.UserID.String()))

	result := r.db.DB.WithContext(ctx).Create(session)
	if result.Error != nil {
		span.RecordError(result.Error)
		return nil, fmt.Errorf("failed to create session: %w", result.Error)
	}

	return session, nil
}

// CreateToken implements Repository.
func (r *repository) CreateToken(ctx context.Context, token *UserToken) (*UserToken, error) {
	ctx, span := r.tracer.Start(ctx, "user.repository.create_token")
	defer span.End()

	span.SetAttributes(attribute.String("user_token.user_id", token.UserID.String()))

	result := r.db.DB.WithContext(ctx).Create(token)
	if result.Error != nil {
		span.RecordError(result.Error)
		return nil, fmt.Errorf("failed to create token: %w", result.Error)
	}

	return token, nil
}

// GetUserByEmailOrUsername implements Repository.
func (r *repository) GetUserByEmailOrUsername(ctx context.Context, emailOrUsername string) (*User, error) {
	ctx, span := r.tracer.Start(ctx, "user.repository.get_user_by_email_or_username")
	defer span.End()

	span.SetAttributes(attribute.String("user.email_or_username", emailOrUsername))

	var user User
	result := r.db.DB.WithContext(ctx).Where("email = ? OR username = ?", emailOrUsername, emailOrUsername).First(&user)
	if result.Error != nil {
		span.RecordError(result.Error)
		return nil, fmt.Errorf("failed to get user: %w", result.Error)
	}

	return &user, nil
}

// CreateUserRole implements Repository.
func (r *repository) CreateUserRole(ctx context.Context, userRole *UserRole) (*UserRole, error) {
	ctx, span := r.tracer.Start(ctx, "user.repository.create_user_role")
	defer span.End()

	span.SetAttributes(attribute.String("user_role.user_id", userRole.UserID.String()), attribute.String("user_role.role_id", userRole.RoleID.String()))

	result := r.db.DB.WithContext(ctx).Create(userRole)
	if result.Error != nil {
		span.RecordError(result.Error)
		return nil, fmt.Errorf("failed to create user role: %w", result.Error)
	}

	return userRole, nil
}

// NewRepository creates a new user repository
func NewRepository(db *database.Database) Repository {
	return &repository{
		db:     db,
		tracer: otel.Tracer("user-repository"),
	}
}

// Create creates a new user in the database
func (r *repository) Create(ctx context.Context, user *User) error {
	ctx, span := r.tracer.Start(ctx, "user.repository.create")
	defer span.End()

	span.SetAttributes(
		attribute.String("user.email", user.Email),
		attribute.String("user.username", user.Username),
	)

	// Generate UUID if not set
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	result := r.db.DB.WithContext(ctx).Create(user)
	if result.Error != nil {
		span.RecordError(result.Error)
		if strings.Contains(result.Error.Error(), "duplicate key") ||
			strings.Contains(result.Error.Error(), "UNIQUE constraint") {
			return fmt.Errorf("user with email %s already exists", user.Email)
		}
		return fmt.Errorf("failed to create user: %w", result.Error)
	}

	// Create user password
	userPassword := &UserPassword{
		UserID:         user.ID,
		PasswordHashed: user.Password,
	}

	userPasswordResult := r.db.DB.WithContext(ctx).Create(userPassword)
	if userPasswordResult.Error != nil {
		span.RecordError(userPasswordResult.Error)
		return fmt.Errorf("failed to create user password: %w", userPasswordResult.Error)
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	ctx, span := r.tracer.Start(ctx, "user.repository.get_by_id")
	defer span.End()

	span.SetAttributes(attribute.String("user.id", id.String()))

	var user User
	result := r.db.DB.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		span.RecordError(result.Error)
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", result.Error)
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	ctx, span := r.tracer.Start(ctx, "user.repository.get_by_email")
	defer span.End()

	span.SetAttributes(attribute.String("user.email", email))

	var user User
	result := r.db.DB.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		span.RecordError(result.Error)
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", result.Error)
	}

	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *repository) GetByUsername(ctx context.Context, username string) (*User, error) {
	ctx, span := r.tracer.Start(ctx, "user.repository.get_by_username")
	defer span.End()

	span.SetAttributes(attribute.String("user.username", username))

	var user User
	result := r.db.DB.WithContext(ctx).Where("username = ?", username).First(&user)
	if result.Error != nil {
		span.RecordError(result.Error)
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", result.Error)
	}

	return &user, nil
}

// Update updates a user in the database
func (r *repository) Update(ctx context.Context, id uuid.UUID, update *User) error {
	ctx, span := r.tracer.Start(ctx, "user.repository.update")
	defer span.End()

	span.SetAttributes(attribute.String("user.id", id.String()))

	// Only update non-zero fields
	updates := map[string]interface{}{}
	if update.FirstName != "" {
		updates["first_name"] = update.FirstName
	}
	if update.LastName != "" {
		updates["last_name"] = update.LastName
	}
	if update.Password != "" {
		updates["password"] = update.Password
	}

	result := r.db.DB.WithContext(ctx).Model(&User{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		span.RecordError(result.Error)
		return fmt.Errorf("failed to update user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Delete deletes a user from the database (soft delete)
func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span := r.tracer.Start(ctx, "user.repository.delete")
	defer span.End()

	span.SetAttributes(attribute.String("user.id", id.String()))

	result := r.db.DB.WithContext(ctx).Delete(&User{}, id)
	if result.Error != nil {
		span.RecordError(result.Error)
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// List retrieves a list of users with pagination and filtering
func (r *repository) List(ctx context.Context, query ListUsersQuery) ([]*User, int64, error) {
	ctx, span := r.tracer.Start(ctx, "user.repository.list")
	defer span.End()

	span.SetAttributes(
		attribute.Int("query.page", query.Page),
		attribute.Int("query.limit", query.Limit),
		attribute.String("query.search", query.Search),
	)

	db := r.db.DB.WithContext(ctx).Model(&User{})

	// Apply filters
	if query.Search != "" {
		searchPattern := "%" + query.Search + "%"
		db = db.Where("email ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	if query.Role != "" {
		db = db.Where("role = ?", query.Role)
	}

	if query.IsActive != nil {
		db = db.Where("is_active = ?", *query.IsActive)
	}

	// Count total records
	var total int64
	if err := db.Count(&total).Error; err != nil {
		span.RecordError(err)
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Apply pagination
	offset := (query.Page - 1) * query.Limit
	db = db.Offset(offset).Limit(query.Limit)

	// Order by created_at desc
	db = db.Order("created_at DESC")

	var users []*User
	if err := db.Find(&users).Error; err != nil {
		span.RecordError(err)
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, total, nil
}

// UpdateLastLogin updates the user's last login timestamp
func (r *repository) UpdateLastLogin(ctx context.Context, userId uuid.UUID, passwordHashed string) error {
	ctx, span := r.tracer.Start(ctx, "user.repository.update_last_login")
	defer span.End()

	span.SetAttributes(attribute.String("user.id", userId.String()))

	// Create user credential
	userCredential := &UserCredential{
		UserID:         userId,
		PasswordHashed: passwordHashed,
		LastLogin:      time.Now(),
	}

	err := r.db.DB.
		WithContext(ctx).
		Model(&UserCredential{}).
		Create(userCredential).Error
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to create user credential: %w", err)
	}
	return nil
}

// GetRoleByName implements Repository.
func (r *repository) GetRoleByName(ctx context.Context, name string) (*Role, error) {
	ctx, span := r.tracer.Start(ctx, "user.repository.get_role_by_name")
	defer span.End()

	span.SetAttributes(attribute.String("role.name", name))

	var role Role
	result := r.db.DB.WithContext(ctx).Where("name = ?", name).First(&role)
	if result.Error != nil {
		span.RecordError(result.Error)
		return nil, fmt.Errorf("failed to get role: %w", result.Error)
	}

	return &role, nil
}

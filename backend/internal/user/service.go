package user

import (
	"context"
	"fmt"
	"strings"
	"time"

	"document-reader-chatbot/configs"
	"document-reader-chatbot/pkg/utils"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Service defines the interface for user business logic
type Service interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (*User, error)
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	UpdateUser(ctx context.Context, id string, req UpdateUserRequest) (*User, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, query ListUsersQuery) ([]*User, int64, error)
	ChangePassword(ctx context.Context, userID string, req ChangePasswordRequest) error
	GetProfile(ctx context.Context, userID string) (*User, error)
}

// service implements the Service interface
type service struct {
	repository Repository
	jwtConfig  configs.JWTConfig
	tracer     trace.Tracer
}

// NewService creates a new user service
func NewService(repository Repository) Service {
	// Load configuration for JWT
	cfg, err := configs.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load configuration: %v", err))
	}

	return &service{
		repository: repository,
		jwtConfig:  cfg.JWT,
		tracer:     otel.Tracer("user-service"),
	}
}

// CreateUser creates a new user
func (s *service) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
	ctx, span := s.tracer.Start(ctx, "user.service.create_user")
	defer span.End()

	span.SetAttributes(
		attribute.String("user.email", req.Email),
		attribute.String("user.username", req.Username),
		attribute.String("user.first_name", req.FirstName),
	)

	// Validate and normalize email
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	// Validate and normalize username
	username := strings.ToLower(strings.TrimSpace(req.Username))
	if username == "" {
		return nil, fmt.Errorf("username is required")
	}

	// Check if user already exists by email
	existingUser, err := s.repository.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	// Check if user already exists by username
	existingUser, err = s.repository.GetByUsername(ctx, username)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with username %s already exists", username)
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user object
	user := &User{
		Username:  username,
		Email:     email,
		Password:  hashedPassword,
		FirstName: strings.TrimSpace(req.FirstName),
		LastName:  strings.TrimSpace(req.LastName),
	}

	// Save user
	if err := s.repository.Create(ctx, user); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Get role by name
	userRole, err := s.repository.GetRoleByName(ctx, RoleUser)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}

	// Create user role
	role := &UserRole{
		UserID: user.ID,
		RoleID: userRole.ID,
	}

	// Save user role
	if _, err := s.repository.CreateUserRole(ctx, role); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to create user role: %w", err)
	}

	return user, nil
}

// Login authenticates a user and returns a JWT token
func (s *service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	ctx, span := s.tracer.Start(ctx, "user.service.login")
	defer span.End()

	span.SetAttributes(attribute.String("user.email_or_username", req.EmailOrUsername))

	// Normalize input
	emailOrUsername := strings.ToLower(strings.TrimSpace(req.EmailOrUsername))
	var user *User
	var err error

	// Try to get user by email first (if it contains @)
	user, err = s.repository.GetUserByEmailOrUsername(ctx, emailOrUsername)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify password
	if err := utils.VerifyPassword(user.Password, req.Password); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(s.jwtConfig, user.ID.String(), user.Username, user.Email)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Update last login
	if err := s.repository.UpdateLastLogin(ctx, user.ID, user.Password); err != nil {
		// Log error but don't fail the login
		span.RecordError(err)
	}

	//Save token
	_, err = s.repository.CreateToken(ctx, &UserToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(s.jwtConfig.ExpirationTime),
	})
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to save token: %w", err)
	}

	//Save session
	_, err = s.repository.CreateSession(ctx, &UserSession{
		UserID:    user.ID,
		LoginTime: time.Now(),
	})
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to save session: %w", err)
	}
	return &LoginResponse{
		User:  user,
		Token: token,
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *service) GetUserByID(ctx context.Context, id string) (*User, error) {
	ctx, span := s.tracer.Start(ctx, "user.service.get_user_by_id")
	defer span.End()

	span.SetAttributes(attribute.String("user.id", id))

	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format")
	}

	user, err := s.repository.GetByID(ctx, userID)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// UpdateUser updates user information
func (s *service) UpdateUser(ctx context.Context, id string, req UpdateUserRequest) (*User, error) {
	ctx, span := s.tracer.Start(ctx, "user.service.update_user")
	defer span.End()

	span.SetAttributes(attribute.String("user.id", id))

	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format")
	}

	// Get existing user
	existingUser, err := s.repository.GetByID(ctx, userID)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Prepare update data
	updateUser := &User{
		FirstName: existingUser.FirstName,
		LastName:  existingUser.LastName,
	}

	// Apply updates
	if req.FirstName != nil {
		updateUser.FirstName = strings.TrimSpace(*req.FirstName)
	}
	if req.LastName != nil {
		updateUser.LastName = strings.TrimSpace(*req.LastName)
	}

	// Update user
	if err := s.repository.Update(ctx, userID, updateUser); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Return updated user
	return s.repository.GetByID(ctx, userID)
}

// DeleteUser deletes a user
func (s *service) DeleteUser(ctx context.Context, id string) error {
	ctx, span := s.tracer.Start(ctx, "user.service.delete_user")
	defer span.End()

	span.SetAttributes(attribute.String("user.id", id))

	userID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid user ID format")
	}

	if err := s.repository.Delete(ctx, userID); err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ListUsers retrieves a list of users with pagination and filtering
func (s *service) ListUsers(ctx context.Context, query ListUsersQuery) ([]*User, int64, error) {
	ctx, span := s.tracer.Start(ctx, "user.service.list_users")
	defer span.End()

	span.SetAttributes(
		attribute.Int("query.page", query.Page),
		attribute.Int("query.limit", query.Limit),
	)

	// Set defaults
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 10
	}
	if query.Limit > 100 {
		query.Limit = 100
	}

	users, total, err := s.repository.List(ctx, query)
	if err != nil {
		span.RecordError(err)
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, total, nil
}

// ChangePassword changes a user's password
func (s *service) ChangePassword(ctx context.Context, userID string, req ChangePasswordRequest) error {
	ctx, span := s.tracer.Start(ctx, "user.service.change_password")
	defer span.End()

	span.SetAttributes(attribute.String("user.id", userID))

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID format")
	}

	// Get user
	user, err := s.repository.GetByID(ctx, parsedUserID)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Verify current password
	if err := utils.VerifyPassword(user.Password, req.CurrentPassword); err != nil {
		return fmt.Errorf("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password using repository pattern
	updateUser := &User{
		Password: hashedPassword,
	}

	// Update the user with new password
	if err := s.repository.Update(ctx, parsedUserID, updateUser); err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// GetProfile retrieves the user's own profile
func (s *service) GetProfile(ctx context.Context, userID string) (*User, error) {
	ctx, span := s.tracer.Start(ctx, "user.service.get_profile")
	defer span.End()

	span.SetAttributes(attribute.String("user.id", userID))

	return s.GetUserByID(ctx, userID)
}

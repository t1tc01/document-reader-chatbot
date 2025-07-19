package user

import (
	"net/http"

	"document-reader-chatbot/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Controller handles HTTP requests for user operations
type Controller struct {
	service Service
	tracer  trace.Tracer
}

// NewController creates a new user controller
func NewController(service Service) *Controller {
	return &Controller{
		service: service,
		tracer:  otel.Tracer("user-controller"),
	}
}

// CreateUser handles POST /users
func (c *Controller) CreateUser(ctx *gin.Context) {
	spanCtx, span := c.tracer.Start(ctx.Request.Context(), "user.controller.create_user")
	defer span.End()

	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		span.RecordError(err)
		utils.BadRequestResponse(ctx, "Invalid request data", err.Error())
		return
	}

	span.SetAttributes(attribute.String("user.email", req.Email))

	user, err := c.service.CreateUser(spanCtx, req)
	if err != nil {
		span.RecordError(err)
		if err.Error() == "user with email "+req.Email+" already exists" {
			utils.ConflictResponse(ctx, err.Error())
			return
		}
		utils.InternalServerErrorResponse(ctx, "Failed to create user")
		return
	}

	utils.CreatedResponse(ctx, "User created successfully", user)
}

// Login handles POST /auth/login
func (c *Controller) Login(ctx *gin.Context) {
	spanCtx, span := c.tracer.Start(ctx.Request.Context(), "user.controller.login")
	defer span.End()

	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		span.RecordError(err)
		utils.BadRequestResponse(ctx, "Invalid request data", err.Error())
		return
	}

	span.SetAttributes(attribute.String("user.email_or_username", req.EmailOrUsername))

	response, err := c.service.Login(spanCtx, req)
	if err != nil {
		span.RecordError(err)
		utils.UnauthorizedResponse(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Login successful", response)
}

// GetUser handles GET /users/:id
func (c *Controller) GetUser(ctx *gin.Context) {
	spanCtx, span := c.tracer.Start(ctx.Request.Context(), "user.controller.get_user")
	defer span.End()

	userID := ctx.Param("id")
	span.SetAttributes(attribute.String("user.id", userID))

	user, err := c.service.GetUserByID(spanCtx, userID)
	if err != nil {
		span.RecordError(err)
		if err.Error() == "user not found" || err.Error() == "invalid user ID format" {
			utils.NotFoundResponse(ctx, "User not found")
			return
		}
		utils.InternalServerErrorResponse(ctx, "Failed to get user")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "User retrieved successfully", user)
}

// UpdateUser handles PUT /users/:id
func (c *Controller) UpdateUser(ctx *gin.Context) {
	spanCtx, span := c.tracer.Start(ctx.Request.Context(), "user.controller.update_user")
	defer span.End()

	userID := ctx.Param("id")
	span.SetAttributes(attribute.String("user.id", userID))

	var req UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		span.RecordError(err)
		utils.BadRequestResponse(ctx, "Invalid request data", err.Error())
		return
	}

	user, err := c.service.UpdateUser(spanCtx, userID, req)
	if err != nil {
		span.RecordError(err)
		if err.Error() == "user not found" || err.Error() == "invalid user ID format" {
			utils.NotFoundResponse(ctx, "User not found")
			return
		}
		utils.InternalServerErrorResponse(ctx, "Failed to update user")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "User updated successfully", user)
}

// DeleteUser handles DELETE /users/:id
func (c *Controller) DeleteUser(ctx *gin.Context) {
	spanCtx, span := c.tracer.Start(ctx.Request.Context(), "user.controller.delete_user")
	defer span.End()

	userID := ctx.Param("id")
	span.SetAttributes(attribute.String("user.id", userID))

	err := c.service.DeleteUser(spanCtx, userID)
	if err != nil {
		span.RecordError(err)
		if err.Error() == "user not found" || err.Error() == "invalid user ID format" {
			utils.NotFoundResponse(ctx, "User not found")
			return
		}
		utils.InternalServerErrorResponse(ctx, "Failed to delete user")
		return
	}

	utils.NoContentResponse(ctx)
}

// ListUsers handles GET /users
func (c *Controller) ListUsers(ctx *gin.Context) {
	spanCtx, span := c.tracer.Start(ctx.Request.Context(), "user.controller.list_users")
	defer span.End()

	var query ListUsersQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		span.RecordError(err)
		utils.BadRequestResponse(ctx, "Invalid query parameters", err.Error())
		return
	}

	// Set defaults if not provided
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	span.SetAttributes(
		attribute.Int("query.page", query.Page),
		attribute.Int("query.limit", query.Limit),
	)

	users, total, err := c.service.ListUsers(spanCtx, query)
	if err != nil {
		span.RecordError(err)
		utils.InternalServerErrorResponse(ctx, "Failed to list users")
		return
	}

	utils.PaginatedResponse(ctx, users, query.Page, query.Limit, total)
}

// ChangePassword handles POST /users/:id/change-password
func (c *Controller) ChangePassword(ctx *gin.Context) {
	spanCtx, span := c.tracer.Start(ctx.Request.Context(), "user.controller.change_password")
	defer span.End()

	userID := ctx.Param("id")
	span.SetAttributes(attribute.String("user.id", userID))

	// Check if user is trying to change their own password or is admin
	contextUserID, exists := ctx.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(ctx, "Authentication required")
		return
	}

	userRole, exists := ctx.Get("user_role")
	if !exists {
		utils.UnauthorizedResponse(ctx, "Authentication required")
		return
	}

	// Only allow users to change their own password or admin to change any password
	if contextUserID != userID && userRole != "admin" {
		utils.ForbiddenResponse(ctx, "You can only change your own password")
		return
	}

	var req ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		span.RecordError(err)
		utils.BadRequestResponse(ctx, "Invalid request data", err.Error())
		return
	}

	err := c.service.ChangePassword(spanCtx, userID, req)
	if err != nil {
		span.RecordError(err)
		if err.Error() == "current password is incorrect" {
			utils.BadRequestResponse(ctx, err.Error(), "")
			return
		}
		if err.Error() == "user not found" || err.Error() == "invalid user ID format" {
			utils.NotFoundResponse(ctx, "User not found")
			return
		}
		utils.InternalServerErrorResponse(ctx, "Failed to change password")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Password changed successfully", nil)
}

// GetProfile handles GET /profile
func (c *Controller) GetProfile(ctx *gin.Context) {
	spanCtx, span := c.tracer.Start(ctx.Request.Context(), "user.controller.get_profile")
	defer span.End()

	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(ctx, "Authentication required")
		return
	}

	userIDStr := userID.(string)
	span.SetAttributes(attribute.String("user.id", userIDStr))

	user, err := c.service.GetProfile(spanCtx, userIDStr)
	if err != nil {
		span.RecordError(err)
		if err.Error() == "user not found" {
			utils.NotFoundResponse(ctx, "User profile not found")
			return
		}
		utils.InternalServerErrorResponse(ctx, "Failed to get user profile")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Profile retrieved successfully", user)
}

// UpdateProfile handles PUT /profile
func (c *Controller) UpdateProfile(ctx *gin.Context) {
	spanCtx, span := c.tracer.Start(ctx.Request.Context(), "user.controller.update_profile")
	defer span.End()

	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(ctx, "Authentication required")
		return
	}

	userIDStr := userID.(string)
	span.SetAttributes(attribute.String("user.id", userIDStr))

	var req UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		span.RecordError(err)
		utils.BadRequestResponse(ctx, "Invalid request data", err.Error())
		return
	}

	// Don't allow users to change their own active status
	req.IsActive = nil

	user, err := c.service.UpdateUser(spanCtx, userIDStr, req)
	if err != nil {
		span.RecordError(err)
		if err.Error() == "user not found" {
			utils.NotFoundResponse(ctx, "User profile not found")
			return
		}
		utils.InternalServerErrorResponse(ctx, "Failed to update user profile")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Profile updated successfully", user)
}

package user

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required,min=2,max=50"`
	LastName  string `json:"last_name" binding:"required,min=2,max=50"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	FirstName *string `json:"first_name,omitempty" binding:"omitempty,min=2,max=50"`
	LastName  *string `json:"last_name,omitempty" binding:"omitempty,min=2,max=50"`
	IsActive  *bool   `json:"is_active,omitempty"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	EmailOrUsername string `json:"email_or_username" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

// ChangePasswordRequest represents the change password request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// ListUsersQuery represents query parameters for listing users
type ListUsersQuery struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	Limit    int    `form:"limit,default=10" binding:"min=1,max=100"`
	Search   string `form:"search"`
	Role     string `form:"role" binding:"omitempty,oneof=user admin"`
	IsActive *bool  `form:"is_active"`
}

package model

import (
	"context"
	"time"
)

// Interface for the user db
type UserDB interface {
	CreateUser(ctx context.Context, user User) error
	GetUserByUsername(ctx context.Context, tenant, realm, username string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	ListUsers(ctx context.Context, tenant, realm string) ([]User, error)
	ListUsersWithPagination(ctx context.Context, tenant, realm string, offset, limit int) ([]User, error)
	CountUsers(ctx context.Context, tenant, realm string) (int64, error)
	GetUserStats(ctx context.Context, tenant, realm string) (*UserStats, error)
	DeleteUser(ctx context.Context, tenant, realm, username string) error
}

type User struct {

	// Unique UUID for the user
	ID string `json:"id"`

	// Organization Context
	Tenant   string `json:"tenant"`
	Realm    string `json:"realm"`
	Username string `json:"username"`

	// User status
	Status string `json:"status"`

	// Identity Information
	DisplayName string `json:"display_name"`
	GivenName   string `json:"given_name"`
	FamilyName  string `json:"family_name"`

	// Additional contact information
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	EmailVerified bool   `json:"email_verified"`
	PhoneVerified bool   `json:"phone_verified"`

	// Locale
	Locale string `json:"locale"`

	// Authentication credentials
	PasswordCredential string `json:"-"`
	WebAuthnCredential string `json:"-"`
	MFACredential      string `json:"-"`

	PasswordLocked bool `json:"password_locked"`
	WebAuthnLocked bool `json:"webauthn_locked"`
	MFALocked      bool `json:"mfa_locked"`

	FailedLoginAttemptsPassword int `json:"failed_login_attempts_password"`
	FailedLoginAttemptsWebAuthn int `json:"failed_login_attempts_webauthn"`
	FailedLoginAttemptsMFA      int `json:"failed_login_attempts_mfa"`

	// User roles and groups
	Roles  []string `json:"roles"`
	Groups []string `json:"groups"`

	// Extensibility
	Attributes map[string]string `json:"attributes"`

	// Audit
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`

	// Federation
	FederatedIDP *string `json:"federated_idp,omitempty"`
	FederatedID  *string `json:"federated_id,omitempty"`

	// Devices
	TrustedDevices string `json:"trusted_devices,omitempty"`
}

// UserStats represents user statistics
type UserStats struct {
	TotalUsers      int64 `json:"total_users"`
	ActiveUsers     int64 `json:"active_users"`
	InactiveUsers   int64 `json:"inactive_users"`
	LockedUsers     int64 `json:"locked_users"`
	EmailVerified   int64 `json:"email_verified"`
	PhoneVerified   int64 `json:"phone_verified"`
	WebAuthnEnabled int64 `json:"webauthn_enabled"`
	MFAEnabled      int64 `json:"mfa_enabled"`
	FederatedUsers  int64 `json:"federated_users"`
}

package cloud

import "time"

type CloudAccountRef struct {
	Provider string `json:"provider"`
	Subject  string `json:"subject"`
}

type CloudRole string

const (
	CloudRoleOwner    CloudRole = "owner"
	CloudRoleAdmin    CloudRole = "admin"
	CloudRoleAuthor   CloudRole = "author"
	CloudRoleCustomer CloudRole = "customer"
	CloudRoleSupport  CloudRole = "support"
)

type CloudSession struct {
	SessionRef       string            `json:"session_ref"`
	AccountRef       CloudAccountRef   `json:"account_ref"`
	Roles            []CloudRole       `json:"roles"`
	IdentityProvider string            `json:"identity_provider,omitempty"`
	IssuedAt         time.Time         `json:"issued_at"`
	ExpiresAt        time.Time         `json:"expires_at"`
	RevokedAt        *time.Time        `json:"revoked_at,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}

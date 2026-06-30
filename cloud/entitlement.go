package cloud

import "time"

type EntitlementStatus string

const (
	EntitlementStatusCandidate EntitlementStatus = "candidate"
	EntitlementStatusActive    EntitlementStatus = "active"
	EntitlementStatusSuspended EntitlementStatus = "suspended"
	EntitlementStatusRevoked   EntitlementStatus = "revoked"
	EntitlementStatusExpired   EntitlementStatus = "expired"
)

type EntitlementScope string

const (
	EntitlementScopePack         EntitlementScope = "pack"
	EntitlementScopeCapability   EntitlementScope = "capability"
	EntitlementScopeOrganization EntitlementScope = "organization"
)

type CloudEntitlement struct {
	EntitlementRef string            `json:"entitlement_ref"`
	OwnerRef       string            `json:"owner_ref,omitempty"`
	AccountRef     CloudAccountRef   `json:"account_ref"`
	PackRef        string            `json:"pack_ref,omitempty"`
	CapabilityRef  string            `json:"capability_ref,omitempty"`
	Scope          EntitlementScope  `json:"scope"`
	Status         EntitlementStatus `json:"status"`
	OrderRef       string            `json:"order_ref,omitempty"`
	LicenseRef     string            `json:"license_ref,omitempty"`
	StartsAt       *time.Time        `json:"starts_at,omitempty"`
	ExpiresAt      *time.Time        `json:"expires_at,omitempty"`
	IssuedAt       time.Time         `json:"issued_at"`
	UpdatedAt      *time.Time        `json:"updated_at,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

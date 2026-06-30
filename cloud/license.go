package cloud

import "time"

type LicenseStatus string

const (
	LicenseStatusCandidate LicenseStatus = "candidate"
	LicenseStatusActive    LicenseStatus = "active"
	LicenseStatusSuspended LicenseStatus = "suspended"
	LicenseStatusRevoked   LicenseStatus = "revoked"
	LicenseStatusExpired   LicenseStatus = "expired"
)

type LicenseToken struct {
	LicenseRef       string            `json:"license_ref"`
	EntitlementRef   string            `json:"entitlement_ref"`
	AccountRef       CloudAccountRef   `json:"account_ref,omitempty"`
	DeviceBindingRef string            `json:"device_binding_ref,omitempty"`
	TokenDigest      string            `json:"token_digest,omitempty"`
	Status           LicenseStatus     `json:"status"`
	IssuedAt         time.Time         `json:"issued_at"`
	ExpiresAt        *time.Time        `json:"expires_at,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}

type LocalActivationToken struct {
	ActivationRef    string        `json:"activation_ref"`
	LicenseRef       string        `json:"license_ref"`
	LocalDeviceRef   string        `json:"local_device_ref"`
	ActivationDigest string        `json:"activation_digest"`
	Status           LicenseStatus `json:"status"`
	IssuedAt         time.Time     `json:"issued_at"`
	ExpiresAt        *time.Time    `json:"expires_at,omitempty"`
}

type LicenseValidationResult struct {
	LicenseRef      string        `json:"license_ref"`
	EntitlementRef  string        `json:"entitlement_ref,omitempty"`
	Status          LicenseStatus `json:"status"`
	Valid           bool          `json:"valid"`
	Reason          string        `json:"reason,omitempty"`
	CheckedAt       time.Time     `json:"checked_at"`
	RefreshAfterSec int           `json:"refresh_after_sec,omitempty"`
}

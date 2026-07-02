package base

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

const (
	PolicySnapshotIssuerBaseKernel = "base-policy-kernel"
)

// OwnerDirective represents an owner's policy directive
type OwnerDirective struct {
	DirectiveID string `json:"directive_id"`
	Content     string `json:"content"`
}

// PolicySet holds a set of active policies
type PolicySet struct {
	SetID      string           `json:"set_id"`
	Version    string           `json:"version,omitempty"`
	Directives []OwnerDirective `json:"directives"`
}

// GatePolicySet is the replayable policy set emitted by the Base policy kernel.
type GatePolicySet struct {
	PolicySetRef         string    `json:"policy_set_ref"`
	Version              string    `json:"version"`
	OwnerPolicyRef       string    `json:"owner_policy_ref,omitempty"`
	RiskPolicyRef        string    `json:"risk_policy_ref,omitempty"`
	SideEffectPolicyRef  string    `json:"side_effect_policy_ref,omitempty"`
	FormalWritePolicyRef string    `json:"formal_write_policy_ref,omitempty"`
	CapabilityPolicyRef  string    `json:"capability_policy_ref,omitempty"`
	ModuleBoundaryRef    string    `json:"module_boundary_ref,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
}

// PolicyKernel identifies the Base-owned kernel that can mint policy snapshots.
type PolicyKernel struct {
	KernelRef string `json:"kernel_ref"`
	Version   string `json:"version"`
}

// PolicySnapshot represents an immutable snapshot of policies used during gate evaluation
type PolicySnapshot struct {
	SnapshotRef   string         `json:"snapshot_ref,omitempty"`
	SnapshotHash  string         `json:"snapshot_hash"`
	PolicySet     PolicySet      `json:"policy_set"`
	GatePolicySet *GatePolicySet `json:"gate_policy_set,omitempty"`
	KernelRef     string         `json:"kernel_ref,omitempty"`
	Version       string         `json:"version,omitempty"`
	IssuedBy      string         `json:"issued_by,omitempty"`
	CreatedAt     time.Time      `json:"created_at,omitempty"`
}

func NewPolicySnapshot(kernel PolicyKernel, policySet GatePolicySet, createdAt time.Time) (*PolicySnapshot, error) {
	if kernel.KernelRef == "" || kernel.Version == "" {
		return nil, errors.New("policy kernel ref and version are required")
	}
	if err := ValidateGatePolicySet(&policySet); err != nil {
		return nil, err
	}
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}
	raw := fmt.Sprintf("%s|%s|%s|%s|%d", kernel.KernelRef, kernel.Version, policySet.PolicySetRef, policySet.Version, createdAt.UnixNano())
	sum := sha256.Sum256([]byte(raw))
	hash := hex.EncodeToString(sum[:])
	return &PolicySnapshot{
		SnapshotRef:  fmt.Sprintf("policy_snapshot_%s", hash[:16]),
		SnapshotHash: hash,
		PolicySet: PolicySet{
			SetID:   policySet.PolicySetRef,
			Version: policySet.Version,
		},
		GatePolicySet: &policySet,
		KernelRef:     kernel.KernelRef,
		Version:       kernel.Version,
		IssuedBy:      PolicySnapshotIssuerBaseKernel,
		CreatedAt:     createdAt,
	}, nil
}

func ValidatePolicySnapshot(snapshot *PolicySnapshot) error {
	if snapshot == nil {
		return errors.New("policy snapshot is required")
	}
	if snapshot.SnapshotRef == "" || snapshot.SnapshotHash == "" {
		return errors.New("policy snapshot ref and hash are required")
	}
	if snapshot.IssuedBy != PolicySnapshotIssuerBaseKernel {
		return errors.New("policy snapshot must be issued by Base policy kernel")
	}
	if snapshot.KernelRef == "" || snapshot.Version == "" {
		return errors.New("policy snapshot kernel ref and version are required")
	}
	if snapshot.CreatedAt.IsZero() {
		return errors.New("policy snapshot created_at is required")
	}
	return nil
}

func ValidateGatePolicySet(set *GatePolicySet) error {
	if set == nil {
		return errors.New("gate policy set is required")
	}
	if set.PolicySetRef == "" || set.Version == "" {
		return errors.New("policy set ref and version are required")
	}
	if set.CreatedAt.IsZero() {
		return errors.New("policy set created_at is required")
	}
	return nil
}

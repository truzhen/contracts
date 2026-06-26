package registry

import "time"

// RegistrySlice is the ONLY registry projection consumers (13 agent context
// assembly, 08 prompt assembly) are allowed to receive. Agent / Model must
// never read the full registry; the type-level contract enforces that a
// consumer signature accepting *RegistrySlice cannot be handed a full
// registry store (docs/architecture/registry-slice-lifecycle-policy.md).
//
// Lifecycle: request -> resolve -> filter -> rank -> mask -> audit -> return
// -> expire -> rebuild. Every returned slice carries the audit ref that makes
// "why did the model see X" answerable after the fact.
type RegistrySlice struct {
	SliceID        string `json:"slice_id"`
	TransactionRef string `json:"transaction_ref"`
	IntentEventID  string `json:"intent_event_id"`
	ActorRef       string `json:"actor_ref"`
	ScopePolicyRef string `json:"scope_policy_ref"`

	Items       []RegistrySliceItem       `json:"items"`
	AllowedRefs []string                  `json:"allowed_refs"`
	BlockedRefs []RegistrySliceBlockedRef `json:"blocked_refs"`

	AuditRef        string    `json:"audit_ref"`
	ResolverVersion string    `json:"resolver_version"`
	StaleWarning    bool      `json:"stale_warning"`
	CreatedAt       time.Time `json:"created_at"`
	ExpiresAt       time.Time `json:"expires_at"`

	CandidateOnly bool `json:"candidate_only"`
	NonFormal     bool `json:"non_formal"`
}

// RegistrySliceItem is a masked, ranked, scope-filtered registry entry. It
// only ever carries refs and short summaries; raw secrets, raw credentials,
// raw endpoints and sensitive payloads are masked out before the item is
// materialized (mask stage keeps the ref, drops the raw content).
type RegistrySliceItem struct {
	Ref           string `json:"ref"`
	Namespace     string `json:"namespace"`
	Kind          string `json:"kind"`
	Source        string `json:"source"`
	ScopeReason   string `json:"scope_reason"`
	Summary       string `json:"summary"`
	ProvenanceRef string `json:"provenance_ref,omitempty"`
	RankScore     int    `json:"rank_score"`
	Masked        bool   `json:"masked"`
	MaskReason    string `json:"mask_reason,omitempty"`
	CandidateOnly bool   `json:"candidate_only"`
	NonFormal     bool   `json:"non_formal"`
}

// RegistrySliceBlockedRef keeps blocked refs auditable: blocked entries must
// never silently disappear (policy: blocked refs enter audit material).
type RegistrySliceBlockedRef struct {
	Ref    string `json:"ref"`
	Reason string `json:"reason"`
}

// Expired reports whether the slice has passed its TTL relative to now.
// Expired slices must not be consumed; the caller has to rebuild.
func (s *RegistrySlice) Expired(now time.Time) bool {
	if s == nil {
		return true
	}
	return !s.ExpiresAt.IsZero() && now.After(s.ExpiresAt)
}

// ContextRefs returns the ref list a consumer may project into prompt or
// agent context material. Everything here is traceable to AuditRef.
func (s *RegistrySlice) ContextRefs() []string {
	if s == nil {
		return nil
	}
	out := make([]string, 0, len(s.Items))
	for _, item := range s.Items {
		out = append(out, item.Ref)
	}
	return out
}

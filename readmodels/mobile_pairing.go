package readmodels

// MobilePairingBootstrapRequest is an authority-free description supplied by a
// phone before it holds a mobile session. The Host, not the phone, creates all
// binding and pairing refs after the PC owner approves the candidate.
type MobilePairingBootstrapRequest struct {
	DeviceLabel    string `json:"device_label"`
	Platform       string `json:"platform"`
	AppInstanceRef string `json:"app_instance_ref"`
	IdempotencyKey string `json:"idempotency_key"`
}

// MobilePairingBootstrapCandidate is a pre-approval ReadModel/Candidate
// projection. It carries no raw bootstrap proof, bearer, formal identity,
// binding, pairing, OwnerDecision, or Receipt truth. The Host creates formal
// binding and pairing facts only after the PC owner approves this candidate.
type MobilePairingBootstrapCandidate struct {
	OK                    bool   `json:"ok"`
	CandidateCreated      bool   `json:"candidate_created"`
	Duplicate             bool   `json:"duplicate"`
	CandidateOnly         bool   `json:"candidate_only"`
	MobileTruthSource     bool   `json:"mobile_truth_source"`
	Status                string `json:"status"`
	CandidateRef          string `json:"candidate_ref"`
	CandidateKind         string `json:"candidate_kind"`
	DeviceLabel           string `json:"device_label"`
	Platform              string `json:"platform"`
	IdempotencyKey        string `json:"idempotency_key"`
	CreatedAt             string `json:"created_at"`
	ProducesOwnerDecision bool   `json:"produces_owner_decision"`
	CredentialState       string `json:"credential_state"`
}

// MobileSessionIssueIntent is the JSON body for a post-approval session issue.
// The one-time bootstrap proof travels only in a controlled request header and
// is intentionally not represented here.
type MobileSessionIssueIntent struct {
	CandidateRef   string `json:"candidate_ref"`
	IdempotencyKey string `json:"idempotency_key"`
}

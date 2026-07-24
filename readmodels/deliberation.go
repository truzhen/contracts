package readmodels

import (
	"fmt"
)

// DeliberationAutomationMode limits the scope of a server-issued automation
// grant. It is a projection of a Base decision, never a client-supplied
// approval signal.
type DeliberationAutomationMode string

const (
	DeliberationAutomationManual          DeliberationAutomationMode = "manual"
	DeliberationAutomationCurrentTurnAuto DeliberationAutomationMode = "current_turn_auto"
	DeliberationAutomationSessionAuto     DeliberationAutomationMode = "session_auto"
)

// DeliberationSessionStatus is the visible lifecycle of a deliberation
// session. A successful-looking state still requires a real Receipt in the
// owning product before it can be treated as formal.
type DeliberationSessionStatus string

const (
	DeliberationSessionCreated                          DeliberationSessionStatus = "created"
	DeliberationSessionCollecting                       DeliberationSessionStatus = "collecting"
	DeliberationSessionAutoActive                       DeliberationSessionStatus = "auto_active"
	DeliberationSessionReadyPartial                     DeliberationSessionStatus = "ready_partial"
	DeliberationSessionReady                            DeliberationSessionStatus = "ready"
	DeliberationSessionSynthesizing                     DeliberationSessionStatus = "synthesizing"
	DeliberationSessionSynthesisReady                   DeliberationSessionStatus = "synthesis_ready"
	DeliberationSessionBlocked                          DeliberationSessionStatus = "blocked"
	DeliberationSessionRecoveryRequired                 DeliberationSessionStatus = "recovery_required"
	DeliberationSessionSuspendedReauthorizationRequired DeliberationSessionStatus = "suspended_reauthorization_required"
	DeliberationSessionVoided                           DeliberationSessionStatus = "voided"
)

// DeliberationTurnStatus identifies the state of one question turn without
// carrying the question or any provider response body.
type DeliberationTurnStatus string

const (
	DeliberationTurnCreated               DeliberationTurnStatus = "created"
	DeliberationTurnCollecting            DeliberationTurnStatus = "collecting"
	DeliberationTurnDispatching           DeliberationTurnStatus = "dispatching"
	DeliberationTurnReadyPartial          DeliberationTurnStatus = "ready_partial"
	DeliberationTurnReady                 DeliberationTurnStatus = "ready"
	DeliberationTurnSynthesizing          DeliberationTurnStatus = "synthesizing"
	DeliberationTurnSynthesisReady        DeliberationTurnStatus = "synthesis_ready"
	DeliberationTurnBlocked               DeliberationTurnStatus = "blocked"
	DeliberationTurnRecoveryRequired      DeliberationTurnStatus = "recovery_required"
	DeliberationTurnAuthorizationRequired DeliberationTurnStatus = "authorization_required"
)

// DeliberationLaneStatus records an adapter lane's controlled progress.
type DeliberationLaneStatus string

const (
	DeliberationLaneSelected         DeliberationLaneStatus = "selected"
	DeliberationLanePreflight        DeliberationLaneStatus = "preflight"
	DeliberationLaneFilling          DeliberationLaneStatus = "filling"
	DeliberationLaneSubmitting       DeliberationLaneStatus = "submitting"
	DeliberationLaneAwaitingResponse DeliberationLaneStatus = "awaiting_response"
	DeliberationLaneCapturing        DeliberationLaneStatus = "capturing"
	DeliberationLaneImported         DeliberationLaneStatus = "imported"
	DeliberationLaneFailed           DeliberationLaneStatus = "failed"
	DeliberationLaneBlocked          DeliberationLaneStatus = "blocked"
	DeliberationLaneCaptchaRequired  DeliberationLaneStatus = "captcha_required"
	DeliberationLaneLoginRequired    DeliberationLaneStatus = "login_required"
	DeliberationLaneAdapterDrift     DeliberationLaneStatus = "adapter_drift"
	DeliberationLaneTimedOut         DeliberationLaneStatus = "timed_out"
	DeliberationLaneExcluded         DeliberationLaneStatus = "excluded"
	DeliberationLaneRecoveryRequired DeliberationLaneStatus = "recovery_required"
	DeliberationLaneCancelled        DeliberationLaneStatus = "cancelled"
)

// DeliberationProviderGroup fixes the product directory grouping. It is a
// display and selection boundary, not a claim that an adapter is eligible to
// execute.
type DeliberationProviderGroup string

const (
	DeliberationProviderGroupChinaCore               DeliberationProviderGroup = "china_core"
	DeliberationProviderGroupChinaExtended           DeliberationProviderGroup = "china_extended"
	DeliberationProviderGroupInternationalSupplement DeliberationProviderGroup = "international_supplement"
	DeliberationProviderGroupCustom                  DeliberationProviderGroup = "custom"
)

// DeliberationAccessMode describes an adapter's declared access path; it does
// not grant that access.
type DeliberationAccessMode string

const (
	DeliberationAccessOfficialAPI           DeliberationAccessMode = "official_api"
	DeliberationAccessProviderAuthorizedWeb DeliberationAccessMode = "provider_authorized_web"
	DeliberationAccessManualWeb             DeliberationAccessMode = "manual_web"
)

type DeliberationReleaseEligibility string

const (
	DeliberationReleaseReady    DeliberationReleaseEligibility = "ready"
	DeliberationReleaseNotReady DeliberationReleaseEligibility = "not_ready"
	DeliberationReleaseBlocked  DeliberationReleaseEligibility = "blocked"
	DeliberationReleaseStale    DeliberationReleaseEligibility = "stale"
)

type DeliberationRuntimeReadiness string

const (
	DeliberationRuntimeReady           DeliberationRuntimeReadiness = "ready"
	DeliberationRuntimeNotReady        DeliberationRuntimeReadiness = "not_ready"
	DeliberationRuntimeBlocked         DeliberationRuntimeReadiness = "blocked"
	DeliberationRuntimeProviderMissing DeliberationRuntimeReadiness = "provider_missing"
)

type DeliberationFailureCode string

const (
	DeliberationFailureLoginRequired       DeliberationFailureCode = "login_required"
	DeliberationFailureCaptchaRequired     DeliberationFailureCode = "captcha_required"
	DeliberationFailureAdapterDrift        DeliberationFailureCode = "adapter_drift"
	DeliberationFailureTimedOut            DeliberationFailureCode = "timed_out"
	DeliberationFailureCaptureConflict     DeliberationFailureCode = "capture_conflict"
	DeliberationFailureRecoveryRequired    DeliberationFailureCode = "recovery_required"
	DeliberationFailurePolicyStale         DeliberationFailureCode = "policy_stale"
	DeliberationFailureDigestMismatch      DeliberationFailureCode = "digest_mismatch"
	DeliberationFailureReceiptAppendFailed DeliberationFailureCode = "receipt_append_failed"
	DeliberationFailureModelNotReady       DeliberationFailureCode = "model_not_ready"
	DeliberationFailureEmergencyStop       DeliberationFailureCode = "emergency_stop"
	DeliberationFailureOriginMismatch      DeliberationFailureCode = "origin_mismatch"
)

type DeliberationCaptureMode string

const (
	DeliberationCaptureNetworkDOMMatch        DeliberationCaptureMode = "network_dom_match"
	DeliberationCaptureDOMFallback            DeliberationCaptureMode = "dom_fallback"
	DeliberationCaptureConflictReviewRequired DeliberationCaptureMode = "conflict_review_required"
)

// DeliberationAllowedOperation is intentionally closed. A provider adapter
// may perform only the operations carried by a server-issued grant.
type DeliberationAllowedOperation string

const (
	DeliberationOperationFillPrompt       DeliberationAllowedOperation = "fill_prompt"
	DeliberationOperationSubmitPrompt     DeliberationAllowedOperation = "submit_prompt"
	DeliberationOperationObserveResponse  DeliberationAllowedOperation = "observe_response"
	DeliberationOperationCaptureResponse  DeliberationAllowedOperation = "capture_response"
	DeliberationOperationRequestSynthesis DeliberationAllowedOperation = "request_synthesis"
)

type DeliberationGrantStatus string

const (
	DeliberationGrantPrepared                         DeliberationGrantStatus = "prepared"
	DeliberationGrantActive                           DeliberationGrantStatus = "active"
	DeliberationGrantConsumed                         DeliberationGrantStatus = "consumed"
	DeliberationGrantExpired                          DeliberationGrantStatus = "expired"
	DeliberationGrantRevoked                          DeliberationGrantStatus = "revoked"
	DeliberationGrantSuspendedReauthorizationRequired DeliberationGrantStatus = "suspended_reauthorization_required"
)

// DeliberationSessionReadModel is the client-safe projection of a session. It
// retains only governed references, hashes and statuses; question and response
// bodies remain in the owning product's controlled artifact store.
type DeliberationSessionReadModel struct {
	SessionRef                  string                                `json:"session_ref"`
	TransactionRef              string                                `json:"transaction_ref"`
	SourceEventID               string                                `json:"source_event_id"`
	OwnerRef                    string                                `json:"owner_ref"`
	Status                      DeliberationSessionStatus             `json:"status"`
	AutomationMode              DeliberationAutomationMode            `json:"automation_mode"`
	AutomationGrant             *DeliberationAutomationGrantReadModel `json:"automation_grant,omitempty"`
	ActiveTurnRef               string                                `json:"active_turn_ref,omitempty"`
	Turns                       []DeliberationTurnReadModel           `json:"turns"`
	SelectedProviderAdapterRefs []string                              `json:"selected_provider_adapter_refs"`
	OccVersion                  int64                                 `json:"occ_version"`
	CreatedAt                   string                                `json:"created_at"`
	UpdatedAt                   string                                `json:"updated_at"`
	ExpiresAt                   string                                `json:"expires_at,omitempty"`
	CandidateOnly               bool                                  `json:"candidate_only"`
	NonFormal                   bool                                  `json:"non_formal"`
}

// DeliberationTurnReadModel is a client-safe one-turn projection. Question
// content is addressed by artifact ref and SHA-256, never serialized here.
type DeliberationTurnReadModel struct {
	TurnRef               string                              `json:"turn_ref"`
	SessionRef            string                              `json:"session_ref"`
	TransactionRef        string                              `json:"transaction_ref"`
	SourceEventID         string                              `json:"source_event_id"`
	Sequence              int                                 `json:"sequence"`
	QuestionArtifactRef   string                              `json:"question_artifact_ref"`
	QuestionSHA256        string                              `json:"question_sha256"`
	Status                DeliberationTurnStatus              `json:"status"`
	ProviderLanes         []DeliberationProviderLaneReadModel `json:"provider_lanes"`
	SelectedLaneCount     int                                 `json:"selected_lane_count"`
	ImportedLaneCount     int                                 `json:"imported_lane_count"`
	SynthesisCandidateRef string                              `json:"synthesis_candidate_ref,omitempty"`
	AutomationAttemptRef  string                              `json:"automation_attempt_ref,omitempty"`
	OccVersion            int64                               `json:"occ_version"`
	CreatedAt             string                              `json:"created_at"`
	UpdatedAt             string                              `json:"updated_at"`
	CandidateOnly         bool                                `json:"candidate_only"`
	NonFormal             bool                                `json:"non_formal"`
}

// DeliberationProviderLaneReadModel exposes adapter release and runtime
// readiness separately. It carries redacted failure codes and artifact refs,
// never browser DOM, credentials or provider response text.
type DeliberationProviderLaneReadModel struct {
	LaneRef                string                         `json:"lane_ref"`
	SessionRef             string                         `json:"session_ref"`
	TurnRef                string                         `json:"turn_ref"`
	ProviderGroup          DeliberationProviderGroup      `json:"provider_group"`
	ProviderResourceRef    string                         `json:"provider_resource_ref"`
	ProviderAdapterRef     string                         `json:"provider_adapter_ref"`
	AdapterVersion         string                         `json:"adapter_version"`
	AdapterDigest          string                         `json:"adapter_digest"`
	AccessMode             DeliberationAccessMode         `json:"access_mode"`
	ReleaseEligibility     DeliberationReleaseEligibility `json:"release_eligibility"`
	RuntimeReadiness       DeliberationRuntimeReadiness   `json:"runtime_readiness"`
	Status                 DeliberationLaneStatus         `json:"status"`
	BrowserRunRef          string                         `json:"browser_run_ref,omitempty"`
	TabBindingRef          string                         `json:"tab_binding_ref,omitempty"`
	GrantRef               string                         `json:"grant_ref,omitempty"`
	ResponseArtifactRef    string                         `json:"response_artifact_ref,omitempty"`
	ResponseSHA256         string                         `json:"response_sha256,omitempty"`
	ResponseByteCount      int                            `json:"response_byte_count,omitempty"`
	ProviderMessageRefHash string                         `json:"provider_message_ref_hash,omitempty"`
	CaptureMode            DeliberationCaptureMode        `json:"capture_mode,omitempty"`
	FailureCode            DeliberationFailureCode        `json:"failure_code,omitempty"`
	FailureSummaryRedacted string                         `json:"failure_summary_redacted,omitempty"`
	ReceiptRefs            []string                       `json:"receipt_refs,omitempty"`
	StartedAt              string                         `json:"started_at,omitempty"`
	CompletedAt            string                         `json:"completed_at,omitempty"`
	ExpiresAt              string                         `json:"expires_at,omitempty"`
	CreatedAt              string                         `json:"created_at"`
	UpdatedAt              string                         `json:"updated_at"`
	CandidateOnly          bool                           `json:"candidate_only"`
	NonFormal              bool                           `json:"non_formal"`
}

// DeliberationAutomationGrantReadModel projects a Base-issued grant for the
// currently selected automation scope. It never permits more than one dispatch
// per lane and is valid only while its owner product can verify DecisionRef.
type DeliberationAutomationGrantReadModel struct {
	GrantRef             string                         `json:"grant_ref"`
	Mode                 DeliberationAutomationMode     `json:"mode"`
	SessionRef           string                         `json:"session_ref"`
	TurnRef              string                         `json:"turn_ref,omitempty"`
	TransactionRef       string                         `json:"transaction_ref"`
	ProviderAdapterRefs  []string                       `json:"provider_adapter_refs"`
	AllowedOperations    []DeliberationAllowedOperation `json:"allowed_operations"`
	DispatchOnConfirm    bool                           `json:"dispatch_on_confirm"`
	QuestionSHA256       string                         `json:"question_sha256,omitempty"`
	IdleExpiresAt        string                         `json:"idle_expires_at"`
	AbsoluteExpiresAt    string                         `json:"absolute_expires_at"`
	MaxTurns             int                            `json:"max_turns"`
	RemainingTurns       int                            `json:"remaining_turns"`
	MaxDispatchesPerLane int                            `json:"max_dispatches_per_lane"`
	Status               DeliberationGrantStatus        `json:"status"`
	DecisionRef          string                         `json:"decision_ref"`
	PolicySnapshotRef    string                         `json:"policy_snapshot_ref"`
	EvidenceRefs         []string                       `json:"evidence_refs"`
	CandidateOnly        bool                           `json:"candidate_only"`
	NonFormal            bool                           `json:"non_formal"`
	CreatedAt            string                         `json:"created_at"`
	UpdatedAt            string                         `json:"updated_at"`
}

// ValidateDeliberationAutomationGrantReadModel rejects all unverifiable or
// over-broad grant projections. Runtime consumers must additionally verify the
// Base-issued decision and expiry; this shape check does not grant authority.
func ValidateDeliberationAutomationGrantReadModel(grant DeliberationAutomationGrantReadModel) error {
	if grant.GrantRef == "" || grant.SessionRef == "" || grant.TransactionRef == "" ||
		grant.IdleExpiresAt == "" || grant.AbsoluteExpiresAt == "" || grant.DecisionRef == "" ||
		grant.PolicySnapshotRef == "" || grant.CreatedAt == "" || grant.UpdatedAt == "" {
		return fmt.Errorf("deliberation automation grant requires governed references and lifecycle timestamps")
	}
	if !grant.CandidateOnly || !grant.NonFormal {
		return fmt.Errorf("deliberation automation grant projection must remain candidate_only and non_formal")
	}
	if grant.Mode != DeliberationAutomationCurrentTurnAuto && grant.Mode != DeliberationAutomationSessionAuto {
		return fmt.Errorf("deliberation automation grant has an unsupported mode")
	}
	if !ValidDeliberationGrantStatus(grant.Status) {
		return fmt.Errorf("deliberation automation grant has an unknown status")
	}
	if len(grant.ProviderAdapterRefs) == 0 || len(grant.AllowedOperations) == 0 || len(grant.EvidenceRefs) == 0 {
		return fmt.Errorf("deliberation automation grant requires selected adapters, operations and evidence")
	}
	for _, operation := range grant.AllowedOperations {
		if !ValidDeliberationAllowedOperation(operation) {
			return fmt.Errorf("deliberation automation grant has unknown operation %q", operation)
		}
	}
	if grant.MaxTurns < 1 || grant.RemainingTurns < 0 || grant.RemainingTurns > grant.MaxTurns {
		return fmt.Errorf("deliberation automation grant has invalid turn limits")
	}
	if grant.MaxDispatchesPerLane != 1 {
		return fmt.Errorf("deliberation automation grant must permit exactly one dispatch per lane")
	}
	if !grant.DispatchOnConfirm {
		return fmt.Errorf("deliberation automation grant must dispatch immediately on confirmed Base decision")
	}
	if grant.Mode == DeliberationAutomationCurrentTurnAuto {
		if grant.TurnRef == "" || !isSHA256(grant.QuestionSHA256) {
			return fmt.Errorf("current-turn automation grant requires turn_ref and question_sha256")
		}
	}
	return nil
}

// ValidateDeliberationSessionReadModel confirms that a client projection only
// contains governed references and closed lifecycle values. It does not make
// the projection a truth source or authorize any automation.
func ValidateDeliberationSessionReadModel(session DeliberationSessionReadModel) error {
	if session.SessionRef == "" || session.TransactionRef == "" || session.SourceEventID == "" ||
		session.OwnerRef == "" || session.CreatedAt == "" || session.UpdatedAt == "" {
		return fmt.Errorf("deliberation session requires governed references and lifecycle timestamps")
	}
	if !session.CandidateOnly || !session.NonFormal {
		return fmt.Errorf("deliberation session projection must remain candidate_only and non_formal")
	}
	if !ValidDeliberationSessionStatus(session.Status) || !ValidDeliberationAutomationMode(session.AutomationMode) {
		return fmt.Errorf("deliberation session has an unknown status or automation mode")
	}
	for _, adapterRef := range session.SelectedProviderAdapterRefs {
		if adapterRef == "" {
			return fmt.Errorf("deliberation session contains an empty selected adapter ref")
		}
	}
	if session.AutomationGrant != nil {
		if err := ValidateDeliberationAutomationGrantReadModel(*session.AutomationGrant); err != nil {
			return fmt.Errorf("deliberation session automation grant: %w", err)
		}
	}
	for _, turn := range session.Turns {
		if err := ValidateDeliberationTurnReadModel(turn); err != nil {
			return fmt.Errorf("deliberation session turn: %w", err)
		}
		if turn.SessionRef != session.SessionRef || turn.TransactionRef != session.TransactionRef || turn.SourceEventID != session.SourceEventID {
			return fmt.Errorf("deliberation session contains a turn from another session, transaction or source event")
		}
	}
	return nil
}

// ValidateDeliberationTurnReadModel checks references and bounded counters.
// It purposely does not decide whether a turn may dispatch; that remains an
// owner-product Base and Gateway decision.
func ValidateDeliberationTurnReadModel(turn DeliberationTurnReadModel) error {
	if turn.TurnRef == "" || turn.SessionRef == "" || turn.TransactionRef == "" ||
		turn.SourceEventID == "" || turn.QuestionArtifactRef == "" || !isSHA256(turn.QuestionSHA256) ||
		turn.CreatedAt == "" || turn.UpdatedAt == "" {
		return fmt.Errorf("deliberation turn requires governed references, question hash and lifecycle timestamps")
	}
	if turn.Sequence < 1 || turn.SelectedLaneCount < 0 || turn.ImportedLaneCount < 0 || turn.ImportedLaneCount > turn.SelectedLaneCount {
		return fmt.Errorf("deliberation turn has invalid sequence or lane counters")
	}
	if !turn.CandidateOnly || !turn.NonFormal || !ValidDeliberationTurnStatus(turn.Status) {
		return fmt.Errorf("deliberation turn must be non-formal and use a known status")
	}
	for _, lane := range turn.ProviderLanes {
		if err := ValidateDeliberationProviderLaneReadModel(lane); err != nil {
			return fmt.Errorf("deliberation turn provider lane: %w", err)
		}
		if lane.SessionRef != turn.SessionRef || lane.TurnRef != turn.TurnRef {
			return fmt.Errorf("deliberation turn contains a provider lane from another session or turn")
		}
	}
	return nil
}

// ValidateDeliberationProviderLaneReadModel rejects unknown adapter state and
// prevents an imported lane from being projected without its artifact and
// ledger references. Actual Receipt existence must still be verified by os-03.
func ValidateDeliberationProviderLaneReadModel(lane DeliberationProviderLaneReadModel) error {
	if lane.LaneRef == "" || lane.SessionRef == "" || lane.TurnRef == "" || lane.ProviderResourceRef == "" ||
		lane.ProviderAdapterRef == "" || lane.AdapterVersion == "" || !isSHA256(lane.AdapterDigest) ||
		lane.CreatedAt == "" || lane.UpdatedAt == "" {
		return fmt.Errorf("deliberation provider lane requires governed references, adapter digest and lifecycle timestamps")
	}
	if !lane.CandidateOnly || !lane.NonFormal {
		return fmt.Errorf("deliberation provider lane projection must remain candidate_only and non_formal")
	}
	if !ValidDeliberationProviderGroup(lane.ProviderGroup) || !ValidDeliberationAccessMode(lane.AccessMode) ||
		!ValidDeliberationReleaseEligibility(lane.ReleaseEligibility) || !ValidDeliberationRuntimeReadiness(lane.RuntimeReadiness) ||
		!ValidDeliberationLaneStatus(lane.Status) {
		return fmt.Errorf("deliberation provider lane has an unknown group, access mode, eligibility, readiness or status")
	}
	if lane.ResponseByteCount < 0 {
		return fmt.Errorf("deliberation provider lane response byte count must not be negative")
	}
	if lane.ResponseSHA256 != "" && !isSHA256(lane.ResponseSHA256) {
		return fmt.Errorf("deliberation provider lane response sha256 is invalid")
	}
	if lane.ProviderMessageRefHash != "" && !isSHA256(lane.ProviderMessageRefHash) {
		return fmt.Errorf("deliberation provider lane provider message ref hash is invalid")
	}
	if lane.CaptureMode != "" && !ValidDeliberationCaptureMode(lane.CaptureMode) {
		return fmt.Errorf("deliberation provider lane capture mode is unknown")
	}
	if lane.FailureCode != "" && !ValidDeliberationFailureCode(lane.FailureCode) {
		return fmt.Errorf("deliberation provider lane failure code is unknown")
	}
	for _, receiptRef := range lane.ReceiptRefs {
		if receiptRef == "" {
			return fmt.Errorf("deliberation provider lane contains an empty receipt ref")
		}
	}
	if lane.Status == DeliberationLaneImported && (lane.ResponseArtifactRef == "" || !isSHA256(lane.ResponseSHA256) || len(lane.ReceiptRefs) == 0) {
		return fmt.Errorf("imported deliberation provider lane requires response artifact, sha256 and receipt refs")
	}
	return nil
}

func ValidDeliberationAutomationMode(value DeliberationAutomationMode) bool {
	switch value {
	case DeliberationAutomationManual, DeliberationAutomationCurrentTurnAuto, DeliberationAutomationSessionAuto:
		return true
	default:
		return false
	}
}

func ValidDeliberationSessionStatus(value DeliberationSessionStatus) bool {
	switch value {
	case DeliberationSessionCreated, DeliberationSessionCollecting, DeliberationSessionAutoActive, DeliberationSessionReadyPartial, DeliberationSessionReady, DeliberationSessionSynthesizing, DeliberationSessionSynthesisReady, DeliberationSessionBlocked, DeliberationSessionRecoveryRequired, DeliberationSessionSuspendedReauthorizationRequired, DeliberationSessionVoided:
		return true
	default:
		return false
	}
}

func ValidDeliberationTurnStatus(value DeliberationTurnStatus) bool {
	switch value {
	case DeliberationTurnCreated, DeliberationTurnCollecting, DeliberationTurnDispatching, DeliberationTurnReadyPartial, DeliberationTurnReady, DeliberationTurnSynthesizing, DeliberationTurnSynthesisReady, DeliberationTurnBlocked, DeliberationTurnRecoveryRequired, DeliberationTurnAuthorizationRequired:
		return true
	default:
		return false
	}
}

func ValidDeliberationProviderGroup(value DeliberationProviderGroup) bool {
	switch value {
	case DeliberationProviderGroupChinaCore, DeliberationProviderGroupChinaExtended, DeliberationProviderGroupInternationalSupplement, DeliberationProviderGroupCustom:
		return true
	default:
		return false
	}
}

func ValidDeliberationAccessMode(value DeliberationAccessMode) bool {
	switch value {
	case DeliberationAccessOfficialAPI, DeliberationAccessProviderAuthorizedWeb, DeliberationAccessManualWeb:
		return true
	default:
		return false
	}
}

func ValidDeliberationReleaseEligibility(value DeliberationReleaseEligibility) bool {
	switch value {
	case DeliberationReleaseReady, DeliberationReleaseNotReady, DeliberationReleaseBlocked, DeliberationReleaseStale:
		return true
	default:
		return false
	}
}

func ValidDeliberationRuntimeReadiness(value DeliberationRuntimeReadiness) bool {
	switch value {
	case DeliberationRuntimeReady, DeliberationRuntimeNotReady, DeliberationRuntimeBlocked, DeliberationRuntimeProviderMissing:
		return true
	default:
		return false
	}
}

func ValidDeliberationFailureCode(value DeliberationFailureCode) bool {
	switch value {
	case DeliberationFailureLoginRequired, DeliberationFailureCaptchaRequired, DeliberationFailureAdapterDrift, DeliberationFailureTimedOut, DeliberationFailureCaptureConflict, DeliberationFailureRecoveryRequired, DeliberationFailurePolicyStale, DeliberationFailureDigestMismatch, DeliberationFailureReceiptAppendFailed, DeliberationFailureModelNotReady, DeliberationFailureEmergencyStop, DeliberationFailureOriginMismatch:
		return true
	default:
		return false
	}
}

func ValidDeliberationCaptureMode(value DeliberationCaptureMode) bool {
	switch value {
	case DeliberationCaptureNetworkDOMMatch, DeliberationCaptureDOMFallback, DeliberationCaptureConflictReviewRequired:
		return true
	default:
		return false
	}
}

func ValidDeliberationAllowedOperation(value DeliberationAllowedOperation) bool {
	switch value {
	case DeliberationOperationFillPrompt, DeliberationOperationSubmitPrompt, DeliberationOperationObserveResponse, DeliberationOperationCaptureResponse, DeliberationOperationRequestSynthesis:
		return true
	default:
		return false
	}
}

func ValidDeliberationLaneStatus(value DeliberationLaneStatus) bool {
	switch value {
	case DeliberationLaneSelected, DeliberationLanePreflight, DeliberationLaneFilling, DeliberationLaneSubmitting, DeliberationLaneAwaitingResponse, DeliberationLaneCapturing, DeliberationLaneImported, DeliberationLaneFailed, DeliberationLaneBlocked, DeliberationLaneCaptchaRequired, DeliberationLaneLoginRequired, DeliberationLaneAdapterDrift, DeliberationLaneTimedOut, DeliberationLaneExcluded, DeliberationLaneRecoveryRequired, DeliberationLaneCancelled:
		return true
	default:
		return false
	}
}

func ValidDeliberationGrantStatus(value DeliberationGrantStatus) bool {
	switch value {
	case DeliberationGrantPrepared, DeliberationGrantActive, DeliberationGrantConsumed, DeliberationGrantExpired, DeliberationGrantRevoked, DeliberationGrantSuspendedReauthorizationRequired:
		return true
	default:
		return false
	}
}

func isSHA256(value string) bool {
	if len(value) != 64 {
		return false
	}
	for _, character := range value {
		if !(character >= '0' && character <= '9') && !(character >= 'a' && character <= 'f') {
			return false
		}
	}
	return true
}

package monitoring

import "time"

type Severity string

const (
	SeverityDebug    Severity = "debug"
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityError    Severity = "error"
	SeveritySecurity Severity = "security"
)

type MonitoringRun struct {
	RunID       string    `json:"run_id"`
	RootDir     string    `json:"root_dir"`
	RunDir      string    `json:"run_dir"`
	GitBranch   string    `json:"git_branch,omitempty"`
	GitCommit   string    `json:"git_commit,omitempty"`
	Status      string    `json:"status"`
	StartedAt   time.Time `json:"started_at"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
}

type MonitoringEvent struct {
	EventID             string    `json:"event_id"`
	RunID               string    `json:"run_id"`
	Sequence            int64     `json:"sequence"`
	PreviousHash        string    `json:"previous_hash"`
	EventHash           string    `json:"event_hash"`
	SourceKind          string    `json:"source_kind"`
	SourceRef           string    `json:"source_ref,omitempty"`
	ComponentRef        string    `json:"component_ref,omitempty"`
	GatewayRef          string    `json:"gateway_ref,omitempty"`
	ProviderRef         string    `json:"provider_ref,omitempty"`
	TransactionRef      string    `json:"transaction_ref,omitempty"`
	CandidateRef        string    `json:"candidate_ref,omitempty"`
	DecisionRef         string    `json:"decision_ref,omitempty"`
	Severity            Severity  `json:"severity"`
	Status              string    `json:"status"`
	Message             string    `json:"message"`
	RedactedPayloadJSON string    `json:"redacted_payload_json,omitempty"`
	PayloadHash         string    `json:"payload_hash,omitempty"`
	SecurePayloadRef    string    `json:"secure_payload_ref,omitempty"`
	EvidenceRefs        []string  `json:"evidence_refs,omitempty"`
	ReceiptCandidateRef string    `json:"receipt_candidate_ref,omitempty"`
	ReceiptRef          string    `json:"receipt_ref,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
}

type CollectorSnapshot struct {
	SnapshotID          string    `json:"snapshot_id"`
	RunID               string    `json:"run_id"`
	CollectorRef        string    `json:"collector_ref"`
	Status              string    `json:"status"`
	Summary             string    `json:"summary"`
	RedactedPayloadJSON string    `json:"redacted_payload_json,omitempty"`
	PayloadHash         string    `json:"payload_hash,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
}

type RedactionFinding struct {
	FindingID string    `json:"finding_id"`
	RunID     string    `json:"run_id"`
	EventID   string    `json:"event_id,omitempty"`
	Kind      string    `json:"kind"`
	Hash      string    `json:"hash"`
	CreatedAt time.Time `json:"created_at"`
}

type ReceiptLink struct {
	EventID      string    `json:"event_id"`
	ReceiptRef   string    `json:"receipt_ref"`
	CandidateRef string    `json:"candidate_ref"`
	PayloadHash  string    `json:"payload_hash"`
	ReceiptHash  string    `json:"receipt_hash"`
	ReceiptSeq   int64     `json:"receipt_sequence"`
	CreatedAt    time.Time `json:"created_at"`
}

type ExportBundle struct {
	BundleID       string    `json:"bundle_id"`
	RunID          string    `json:"run_id"`
	Path           string    `json:"path"`
	SHA256         string    `json:"sha256"`
	EventChainHead string    `json:"event_chain_head"`
	ReceiptRef     string    `json:"receipt_ref,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type FaultSeverity string

const (
	FaultSeverityInfo     FaultSeverity = "info"
	FaultSeverityWarning  FaultSeverity = "warning"
	FaultSeverityError    FaultSeverity = "error"
	FaultSeveritySecurity FaultSeverity = "security"
)

type FaultSignature struct {
	SignatureID      string        `json:"signature_id"`
	ComponentRef     string        `json:"component_ref"`
	GatewayRef       string        `json:"gateway_ref,omitempty"`
	ProviderRef      string        `json:"provider_ref,omitempty"`
	FaultKind        string        `json:"fault_kind"`
	ErrorCode        string        `json:"error_code"`
	StackHash        string        `json:"stack_hash,omitempty"`
	MessageHash      string        `json:"message_hash"`
	Status           string        `json:"status"`
	Severity         FaultSeverity `json:"severity"`
	FirstEventRef    string        `json:"first_event_ref"`
	LastEventRef     string        `json:"last_event_ref"`
	OccurrenceCount  int           `json:"occurrence_count"`
	RedactionVersion string        `json:"redaction_version"`
}

type FaultIncident struct {
	IncidentID         string         `json:"incident_id"`
	RunID              string         `json:"run_id"`
	Signature          FaultSignature `json:"signature"`
	StartedAt          time.Time      `json:"started_at"`
	EndedAt            time.Time      `json:"ended_at,omitempty"`
	AffectedPageID     string         `json:"affected_page_id,omitempty"`
	AffectedActionRef  string         `json:"affected_action_ref,omitempty"`
	UserVisibleMessage string         `json:"user_visible_message"`
	RecommendedAction  string         `json:"recommended_action"`
	EventRefs          []string       `json:"event_refs"`
	EvidenceRefs       []string       `json:"evidence_refs,omitempty"`
	ReceiptRefs        []string       `json:"receipt_refs,omitempty"`
}

type SupportDiagnosticBundle struct {
	BundleID          string    `json:"bundle_id"`
	RunID             string    `json:"run_id"`
	IncidentIDs       []string  `json:"incident_ids"`
	Path              string    `json:"path"`
	SHA256            string    `json:"sha256"`
	SizeBytes         int64     `json:"size_bytes"`
	EventChainHead    string    `json:"event_chain_head"`
	RedactionVersion  string    `json:"redaction_version"`
	RedactionPassed   bool      `json:"redaction_passed"`
	ForbiddenFindings []string  `json:"forbidden_findings,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}

type SupportUploadCandidate struct {
	CandidateRef        string    `json:"candidate_ref"`
	BundleID            string    `json:"bundle_id"`
	RunID               string    `json:"run_id"`
	ConsentRef          string    `json:"consent_ref"`
	UploadMode          string    `json:"upload_mode"`
	CloudEndpointRef    string    `json:"cloud_endpoint_ref"`
	TruzhenUserIDHash   string    `json:"truzhen_user_id_hash,omitempty"`
	Anonymous           bool      `json:"anonymous"`
	CandidateOnly       bool      `json:"candidate_only"`
	NonFormal           bool      `json:"non_formal"`
	ReceiptCandidateRef string    `json:"receipt_candidate_ref"`
	EvidenceRefs        []string  `json:"evidence_refs"`
	CreatedAt           time.Time `json:"created_at"`
}

type SupportUploadAck struct {
	SupportBundleRef     string             `json:"support_bundle_ref"`
	ServerReceiptRef     string             `json:"server_receipt_ref"`
	FaultSignatureID     string             `json:"fault_signature_id"`
	FeedbackStatus       string             `json:"feedback_status"`
	SymbolicationStatus  string             `json:"symbolication_status,omitempty"`
	TopSymbolicatedFrame *SymbolicatedFrame `json:"top_symbolicated_frame,omitempty"`
	KnownIssueRef        string             `json:"known_issue_ref,omitempty"`
	Message              string             `json:"message"`
	ReceivedAt           time.Time          `json:"received_at"`
}

type SymbolicatedFrame struct {
	SourceFile   string `json:"source_file"`
	Line         int    `json:"line"`
	Column       int    `json:"column"`
	FunctionName string `json:"function_name"`
	Confidence   string `json:"confidence"`
}

type BuildSymbolManifest struct {
	BuildID         string              `json:"build_id"`
	ReleaseID       string              `json:"release_id"`
	GitCommit       string              `json:"git_commit"`
	BuildChannel    string              `json:"build_channel"`
	Platform        string              `json:"platform"`
	FrontendChunks  []FrontendChunkMap  `json:"frontend_chunks"`
	NativeArtifacts []NativeSymbolEntry `json:"native_artifacts"`
	CreatedAt       time.Time           `json:"created_at"`
}

type FrontendChunkMap struct {
	ChunkName       string `json:"chunk_name"`
	ChunkHash       string `json:"chunk_hash"`
	OutputFile      string `json:"output_file"`
	SourceMapRef    string `json:"source_map_ref"`
	SourceMapSHA256 string `json:"source_map_sha256"`
}

type NativeSymbolEntry struct {
	NativeSymbolID    string `json:"native_symbol_id"`
	TargetTriple      string `json:"target_triple"`
	BinaryName        string `json:"binary_name"`
	BinarySHA256      string `json:"binary_sha256"`
	SymbolArtifactRef string `json:"symbol_artifact_ref"`
	SymbolSHA256      string `json:"symbol_sha256"`
}

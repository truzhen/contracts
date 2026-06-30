package cloud

import "time"

type CloudDeployTarget string

const (
	CloudDeployTargetLocal   CloudDeployTarget = "local"
	CloudDeployTargetStaging CloudDeployTarget = "staging"
	CloudDeployTargetProd    CloudDeployTarget = "production"
)

type CloudReleaseSmokeState string

const (
	CloudReleaseSmokePending CloudReleaseSmokeState = "pending"
	CloudReleaseSmokePassed  CloudReleaseSmokeState = "passed"
	CloudReleaseSmokeFailed  CloudReleaseSmokeState = "failed"
	CloudReleaseSmokeBlocked CloudReleaseSmokeState = "blocked"
)

type CloudReleaseCandidate struct {
	CandidateRef     string            `json:"candidate_ref"`
	CommitSHA        string            `json:"commit_sha"`
	ArtifactDigest   string            `json:"artifact_digest"`
	MigrationVersion string            `json:"migration_version,omitempty"`
	ConfigVersion    string            `json:"config_version,omitempty"`
	DeployTarget     CloudDeployTarget `json:"deploy_target"`
	CreatedAt        time.Time         `json:"created_at"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}

type CloudReleaseReceipt struct {
	ReleaseRef     string                 `json:"release_ref"`
	CandidateRef   string                 `json:"candidate_ref"`
	DeployTarget   CloudDeployTarget      `json:"deploy_target"`
	ArtifactDigest string                 `json:"artifact_digest"`
	SmokeState     CloudReleaseSmokeState `json:"smoke_state"`
	PreflightRef   string                 `json:"preflight_ref,omitempty"`
	RollbackRef    string                 `json:"rollback_ref,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
}

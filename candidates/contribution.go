package candidates

type ContributionStatus string

const (
	ContributionStatusProposed      ContributionStatus = "proposed"
	ContributionStatusOwnerApproved ContributionStatus = "owner_approved"
	ContributionStatusSubmitted     ContributionStatus = "submitted"
	ContributionStatusAuthorAdopted ContributionStatus = "author_adopted"
	ContributionStatusRejected      ContributionStatus = "rejected"
)

type MigrationRiskColor string

const (
	MigrationRiskGreen  MigrationRiskColor = "green"
	MigrationRiskYellow MigrationRiskColor = "yellow"
	MigrationRiskOrange MigrationRiskColor = "orange"
	MigrationRiskRed    MigrationRiskColor = "red"
)

type MigrationOwnerAction string

const (
	MigrationOwnerConfirm MigrationOwnerAction = "owner_confirm"
	MigrationBlocked      MigrationOwnerAction = "blocked"
)

type MigrationStatus string

const (
	MigrationStatusProposed                      MigrationStatus = "proposed"
	MigrationStatusOwnerRejected                 MigrationStatus = "owner_rejected"
	MigrationStatusOwnerApprovedPendingMigration MigrationStatus = "owner_approved_pending_migration"
)

// PackUsageContributionCandidate 是项目复盘经 Owner 授权、脱敏后拟贡献给 Pack 作者的候选。
// 它只定义候选形状；真实提交、采纳和回执落点由基座与云端实现。
type PackUsageContributionCandidate struct {
	CandidateID            string             `json:"candidate_id"`
	SourceTransactionRef   string             `json:"source_transaction_ref"`
	SourcePackRef          string             `json:"source_pack_ref"`
	SourcePackVersionRef   string             `json:"source_pack_version_ref"`
	SanitizationReceiptRef string             `json:"sanitization_receipt_ref"`
	OwnerDecisionRef       string             `json:"owner_decision_ref"`
	SummaryRedacted        string             `json:"summary_redacted"`
	Status                 ContributionStatus `json:"status"`
	CandidateOnly          bool               `json:"candidate_only"`
	NonFormal              bool               `json:"non_formal"`
	CreatedAt              string             `json:"created_at"`
}

// PackVersionMigrationCandidate 是“运行中项目从 pinned 版本迁到新版本”的候选。
// 本契约只允许生成与展示；真实切 pin 必须另经 Owner + Base Gate。
type PackVersionMigrationCandidate struct {
	CandidateID         string               `json:"candidate_id"`
	TransactionRef      string               `json:"transaction_ref"`
	PackRef             string               `json:"pack_ref"`
	FromPackVersionRef  string               `json:"from_pack_version_ref"`
	ToPackVersionRef    string               `json:"to_pack_version_ref"`
	DiffSummary         string               `json:"diff_summary,omitempty"`
	RiskColor           MigrationRiskColor   `json:"risk_color"`
	RequiredOwnerAction MigrationOwnerAction `json:"required_owner_action"`
	Status              MigrationStatus      `json:"status"`
	CandidateOnly       bool                 `json:"candidate_only"`
	NonFormal           bool                 `json:"non_formal"`
	CreatedAt           string               `json:"created_at"`
}

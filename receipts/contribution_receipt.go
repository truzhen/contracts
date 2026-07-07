package receipts

type ContributionReceiptKind string

const (
	ContributionReceiptKindOwnerAuthorization ContributionReceiptKind = "owner_authorization"
	ContributionReceiptKindSanitization       ContributionReceiptKind = "sanitization"
	ContributionReceiptKindSubmission         ContributionReceiptKind = "submission"
	ContributionReceiptKindAuthorAdoption     ContributionReceiptKind = "author_adoption"
)

// ContributionReceipt 记录贡献链正式动作的可反查回执引用。
// 本类型不实现账本 append；正式回执仍由 03 ledger 生成。
type ContributionReceipt struct {
	ReceiptRef               string                  `json:"receipt_ref"`
	Kind                     ContributionReceiptKind `json:"kind"`
	ContributionCandidateRef string                  `json:"contribution_candidate_ref"`
	SourcePackVersionRef     string                  `json:"source_pack_version_ref"`
	ResultPackVersionRef     string                  `json:"result_pack_version_ref,omitempty"`
	CreatedAt                string                  `json:"created_at"`
}

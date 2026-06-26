package candidates

// Candidate defines the common interface for all candidate objects.
// A-1 Type Isolation Redline.
type Candidate interface {
	IsCandidate() bool
	IsFormal() bool
}

// CandidateEnvelope 是所有候选对象的统一包装。
// 根据纪律，Candidate 默认 candidate_only=true、non_formal=true。
type CandidateEnvelope struct {
	CandidateRef   string      `json:"candidate_ref,omitempty"`
	TransactionRef string      `json:"transaction_ref,omitempty"`
	SourceEventID  string      `json:"source_event_id,omitempty"`
	ReceiptRef     string      `json:"receipt_ref,omitempty"`
	CandidateOnly  bool        `json:"candidate_only"`
	NonFormal      bool        `json:"non_formal"`
	Payload        interface{} `json:"payload"`
}

func (c *CandidateEnvelope) IsCandidate() bool {
	return c.CandidateOnly
}

func (c *CandidateEnvelope) IsFormal() bool {
	return !c.NonFormal
}

// NewCandidateEnvelope ensures default flags are set.
func NewCandidateEnvelope(payload interface{}) *CandidateEnvelope {
	return &CandidateEnvelope{
		CandidateOnly: true,
		NonFormal:     true,
		Payload:       payload,
	}
}

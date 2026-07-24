package candidates

import "fmt"

// DeliberationConfidence is a bounded confidence label for an AI-produced
// synthesis item. It is not an authority or a formal decision.
type DeliberationConfidence string

const (
	DeliberationConfidenceLow    DeliberationConfidence = "low"
	DeliberationConfidenceMedium DeliberationConfidence = "medium"
	DeliberationConfidenceHigh   DeliberationConfidence = "high"
)

// DeliberationSynthesisItem is a concise candidate claim linked to the
// selected material set. It must never contain a copied provider response.
type DeliberationSynthesisItem struct {
	ItemRef      string                 `json:"item_ref"`
	Summary      string                 `json:"summary"`
	Confidence   DeliberationConfidence `json:"confidence"`
	MaterialRefs []string               `json:"material_refs"`
}

// DeliberationMissingLane explains a lane omitted from a synthesis candidate
// without implying that it succeeded.
type DeliberationMissingLane struct {
	LaneRef    string `json:"lane_ref"`
	ReasonCode string `json:"reason_code"`
	ReceiptRef string `json:"receipt_ref,omitempty"`
}

// DeliberationSynthesisCandidate is an AI-produced, non-formal output for one
// deliberation turn. Formal actions remain owned by the consuming product and
// require a separately verified Base decision and Receipt.
type DeliberationSynthesisCandidate struct {
	CandidateRef         string                      `json:"candidate_ref"`
	SessionRef           string                      `json:"session_ref"`
	TurnRef              string                      `json:"turn_ref"`
	TransactionRef       string                      `json:"transaction_ref"`
	SourceEventID        string                      `json:"source_event_id"`
	ModelRunRef          string                      `json:"model_run_ref"`
	SelectedLaneRefs     []string                    `json:"selected_lane_refs"`
	MaterialRefs         []string                    `json:"material_refs"`
	CoreConclusions      []DeliberationSynthesisItem `json:"core_conclusions"`
	Consensus            []DeliberationSynthesisItem `json:"consensus"`
	Disagreements        []DeliberationSynthesisItem `json:"disagreements"`
	StrongestObjections  []DeliberationSynthesisItem `json:"strongest_objections"`
	Unknowns             []DeliberationSynthesisItem `json:"unknowns"`
	RecommendedNextSteps []DeliberationSynthesisItem `json:"recommended_next_steps"`
	MissingLanes         []DeliberationMissingLane   `json:"missing_lanes"`
	SingleSource         bool                        `json:"single_source"`
	SourceCount          int                         `json:"source_count"`
	ReceiptRef           string                      `json:"receipt_ref"`
	CandidateOnly        bool                        `json:"candidate_only"`
	NonFormal            bool                        `json:"non_formal"`
	CreatedAt            string                      `json:"created_at"`
}

func (c DeliberationSynthesisCandidate) IsCandidate() bool { return c.CandidateOnly }

func (c DeliberationSynthesisCandidate) IsFormal() bool { return !c.NonFormal }

// ValidateDeliberationSynthesisCandidate keeps synthesis output in the
// candidate domain, requires a true Receipt reference, and rejects claims that
// cannot be traced to the candidate's selected material set.
func ValidateDeliberationSynthesisCandidate(candidate DeliberationSynthesisCandidate) error {
	if candidate.CandidateRef == "" || candidate.SessionRef == "" || candidate.TurnRef == "" ||
		candidate.TransactionRef == "" || candidate.SourceEventID == "" || candidate.ModelRunRef == "" ||
		candidate.ReceiptRef == "" || candidate.CreatedAt == "" {
		return fmt.Errorf("deliberation synthesis candidate requires governed references and timestamp")
	}
	if !candidate.CandidateOnly || !candidate.NonFormal {
		return fmt.Errorf("deliberation synthesis candidate must remain candidate_only and non_formal")
	}
	if len(candidate.SelectedLaneRefs) == 0 || len(candidate.MaterialRefs) == 0 || len(candidate.CoreConclusions) == 0 {
		return fmt.Errorf("deliberation synthesis candidate requires selected lanes, material refs and core conclusions")
	}
	if candidate.SourceCount < 1 || candidate.SingleSource != (candidate.SourceCount == 1) {
		return fmt.Errorf("deliberation synthesis candidate has inconsistent source count")
	}
	available := make(map[string]struct{}, len(candidate.MaterialRefs))
	for _, ref := range candidate.MaterialRefs {
		if ref == "" {
			return fmt.Errorf("deliberation synthesis candidate contains an empty material ref")
		}
		available[ref] = struct{}{}
	}
	for _, group := range [][]DeliberationSynthesisItem{
		candidate.CoreConclusions,
		candidate.Consensus,
		candidate.Disagreements,
		candidate.StrongestObjections,
		candidate.Unknowns,
		candidate.RecommendedNextSteps,
	} {
		for _, item := range group {
			if err := validateDeliberationSynthesisItem(item, available); err != nil {
				return err
			}
		}
	}
	for _, lane := range candidate.MissingLanes {
		if lane.LaneRef == "" || lane.ReasonCode == "" {
			return fmt.Errorf("deliberation synthesis candidate has incomplete missing-lane explanation")
		}
	}
	return nil
}

func validateDeliberationSynthesisItem(item DeliberationSynthesisItem, available map[string]struct{}) error {
	if item.ItemRef == "" || item.Summary == "" || !ValidDeliberationConfidence(item.Confidence) || len(item.MaterialRefs) == 0 {
		return fmt.Errorf("deliberation synthesis item is incomplete or has unknown confidence")
	}
	for _, materialRef := range item.MaterialRefs {
		if _, found := available[materialRef]; !found {
			return fmt.Errorf("deliberation synthesis item references material outside its candidate")
		}
	}
	return nil
}

func ValidDeliberationConfidence(value DeliberationConfidence) bool {
	switch value {
	case DeliberationConfidenceLow, DeliberationConfidenceMedium, DeliberationConfidenceHigh:
		return true
	default:
		return false
	}
}

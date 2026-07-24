package candidates

import "testing"

func validDeliberationSynthesisCandidate() DeliberationSynthesisCandidate {
	return DeliberationSynthesisCandidate{
		CandidateRef:     "deliberation_synthesis_candidate://owner/1",
		SessionRef:       "deliberation_session://owner/1",
		TurnRef:          "deliberation_turn://owner/1/1",
		TransactionRef:   "transaction://owner/1",
		SourceEventID:    "intent_event://owner/1",
		ModelRunRef:      "model_run://owner/1",
		SelectedLaneRefs: []string{"deliberation_lane://owner/1"},
		MaterialRefs:     []string{"artifact://owner/1"},
		CoreConclusions: []DeliberationSynthesisItem{{
			ItemRef:      "deliberation_synthesis_item://owner/1/conclusion/1",
			Summary:      "候选结论",
			Confidence:   DeliberationConfidenceMedium,
			MaterialRefs: []string{"artifact://owner/1"},
		}},
		SingleSource:  true,
		SourceCount:   1,
		ReceiptRef:    "receipt://owner/1",
		CandidateOnly: true,
		NonFormal:     true,
		CreatedAt:     "2026-07-24T00:00:00Z",
	}
}

func TestValidateDeliberationSynthesisCandidate(t *testing.T) {
	if err := ValidateDeliberationSynthesisCandidate(validDeliberationSynthesisCandidate()); err != nil {
		t.Fatalf("valid synthesis candidate rejected: %v", err)
	}

	tests := []struct {
		name   string
		mutate func(*DeliberationSynthesisCandidate)
	}{
		{
			name: "formal candidate",
			mutate: func(candidate *DeliberationSynthesisCandidate) {
				candidate.CandidateOnly = false
			},
		},
		{
			name: "unmapped conclusion",
			mutate: func(candidate *DeliberationSynthesisCandidate) {
				candidate.CoreConclusions[0].MaterialRefs = nil
			},
		},
		{
			name: "foreign material",
			mutate: func(candidate *DeliberationSynthesisCandidate) {
				candidate.CoreConclusions[0].MaterialRefs = []string{"artifact://other-turn/1"}
			},
		},
		{
			name: "unknown confidence",
			mutate: func(candidate *DeliberationSynthesisCandidate) {
				candidate.CoreConclusions[0].Confidence = "certain"
			},
		},
		{
			name: "single source mismatch",
			mutate: func(candidate *DeliberationSynthesisCandidate) {
				candidate.SingleSource = false
			},
		},
		{
			name: "receipt absent",
			mutate: func(candidate *DeliberationSynthesisCandidate) {
				candidate.ReceiptRef = ""
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			candidate := validDeliberationSynthesisCandidate()
			test.mutate(&candidate)
			if err := ValidateDeliberationSynthesisCandidate(candidate); err == nil {
				t.Fatal("invalid synthesis candidate was accepted")
			}
		})
	}
}

func TestDeliberationSynthesisCandidateImplementsCandidate(t *testing.T) {
	candidate := validDeliberationSynthesisCandidate()
	if !candidate.IsCandidate() || candidate.IsFormal() {
		t.Fatal("synthesis must remain candidate-only and non-formal")
	}
}

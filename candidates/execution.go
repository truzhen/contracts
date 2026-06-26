package candidates

type ExecutionIntentCandidate struct {
	CandidateEnvelope
	Command string `json:"command"`
}

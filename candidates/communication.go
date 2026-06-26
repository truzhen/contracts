package candidates

type CommunicationDraftCandidate struct {
	CandidateEnvelope
	Message string `json:"message"`
}

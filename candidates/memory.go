package candidates

type MemoryWriteCandidate struct {
	CandidateEnvelope
	Content string `json:"content"`
}

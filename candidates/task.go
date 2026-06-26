package candidates

type TaskCandidate struct {
	CandidateEnvelope
	TaskName string `json:"task_name"`
}

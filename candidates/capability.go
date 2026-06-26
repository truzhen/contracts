package candidates

import "github.com/truzhen/contracts/registry"

type CapabilityInvocationCandidate struct {
	CandidateEnvelope
	RegistryRef registry.RegistryRef `json:"registry_ref"`
	SkillRef    *registry.SkillRef   `json:"skill_ref,omitempty"`
	Parameters  interface{}          `json:"parameters"`
}

func NewCapabilityInvocationCandidate(regRef registry.RegistryRef, skillRef *registry.SkillRef, params interface{}) *CapabilityInvocationCandidate {
	return &CapabilityInvocationCandidate{
		CandidateEnvelope: *NewCandidateEnvelope(nil),
		RegistryRef:       regRef,
		SkillRef:          skillRef,
		Parameters:        params,
	}
}

package base

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestOwnerDelegationGrantLegacyJSONRoundTripDoesNotDefaultExecution(t *testing.T) {
	raw := []byte(`{
		"grant_id":"grant:legacy",
		"owner_decision_ref":"owner_decision:001",
		"delegate_ref":"agent://secretary_chief",
		"scope":{
			"task_types":["stage"],
			"risk_ceiling":"medium",
			"quota":{"per_day":3}
		},
		"expires_at":"2026-07-11T00:00:00Z",
		"revocable":true,
		"status":"active"
	}`)

	var grant OwnerDelegationGrant
	if err := json.Unmarshal(raw, &grant); err != nil {
		t.Fatalf("Unmarshal legacy grant: %v", err)
	}
	if grant.Scope.ExecutionScope != nil {
		t.Fatalf("legacy grant must not default execution_scope: %+v", grant.Scope.ExecutionScope)
	}
	if err := ValidateOwnerDelegationGrant(&grant); err != nil {
		t.Fatalf("ValidateOwnerDelegationGrant legacy: %v", err)
	}

	out, err := json.Marshal(grant)
	if err != nil {
		t.Fatalf("Marshal legacy grant: %v", err)
	}
	if strings.Contains(string(out), "execution_scope") || strings.Contains(string(out), "execution") {
		t.Fatalf("legacy grant JSON must omit execution fields, got %s", out)
	}
}

func TestDelegationSubjectOmitsExecutionWhenNil(t *testing.T) {
	subject := DelegationSubject{
		CandidateRef:   "candidate:001",
		TransactionRef: "transaction:001",
		TaskType:       "stage",
		RiskLevel:      RiskLow,
	}

	out, err := json.Marshal(subject)
	if err != nil {
		t.Fatalf("Marshal delegation subject: %v", err)
	}
	if strings.Contains(string(out), "execution") {
		t.Fatalf("delegation subject JSON must omit nil execution, got %s", out)
	}
}

func TestDelegationExecutionSubjectMatchesScope(t *testing.T) {
	scope := validDelegationScopeWithExecution()
	subject := validDelegationExecutionSubject()

	if err := ValidateDelegationScope(scope); err != nil {
		t.Fatalf("ValidateDelegationScope: %v", err)
	}
	if err := ValidateDelegationExecutionSubject(subject); err != nil {
		t.Fatalf("ValidateDelegationExecutionSubject: %v", err)
	}
	if err := DelegationExecutionWithinScope(scope, subject); err != nil {
		t.Fatalf("DelegationExecutionWithinScope: %v", err)
	}
}

func TestDelegationExecutionSubjectRejectsEveryOutOfScopeDimension(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*DelegationExecutionSubject)
	}{
		{name: "capability", mutate: func(s *DelegationExecutionSubject) { s.CapabilityRef = "capability:other" }},
		{name: "workroot", mutate: func(s *DelegationExecutionSubject) { s.WorkrootRef = "workroot:other" }},
		{name: "provider", mutate: func(s *DelegationExecutionSubject) { s.ProviderRef = "provider:other" }},
		{name: "sandbox", mutate: func(s *DelegationExecutionSubject) { s.SandboxProfileRef = "sandbox:other" }},
		{name: "network policy", mutate: func(s *DelegationExecutionSubject) { s.NetworkPolicy = ExecutionNetworkPolicyGatedBridge }},
		{name: "max runs", mutate: func(s *DelegationExecutionSubject) { s.ConsumedRuns = 3 }},
		{name: "max duration", mutate: func(s *DelegationExecutionSubject) { s.ConsumedDurationSeconds = 301 }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := validDelegationScopeWithExecution()
			subject := validDelegationExecutionSubject()
			tt.mutate(subject)
			if err := DelegationExecutionWithinScope(scope, subject); err == nil {
				t.Fatalf("expected out-of-scope %s rejection", tt.name)
			}
		})
	}
}

func TestDelegationExecutionScopeRejectsInvalidShape(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*DelegationExecutionScope)
	}{
		{name: "empty capability", mutate: func(s *DelegationExecutionScope) { s.CapabilityRefs = []string{"capability:writer", ""} }},
		{name: "duplicate capability", mutate: func(s *DelegationExecutionScope) {
			s.CapabilityRefs = []string{"capability:writer", "capability:writer"}
		}},
		{name: "empty workroot", mutate: func(s *DelegationExecutionScope) { s.WorkrootRef = "" }},
		{name: "empty provider", mutate: func(s *DelegationExecutionScope) { s.ProviderRefs = []string{"provider:local-exec", ""} }},
		{name: "duplicate provider", mutate: func(s *DelegationExecutionScope) {
			s.ProviderRefs = []string{"provider:local-exec", "provider:local-exec"}
		}},
		{name: "empty sandbox", mutate: func(s *DelegationExecutionScope) { s.SandboxProfileRef = "" }},
		{name: "gated bridge ceiling", mutate: func(s *DelegationExecutionScope) { s.NetworkPolicyCeiling = ExecutionNetworkPolicyGatedBridge }},
		{name: "zero max runs", mutate: func(s *DelegationExecutionScope) { s.MaxRuns = 0 }},
		{name: "zero max duration", mutate: func(s *DelegationExecutionScope) { s.MaxDurationSeconds = 0 }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := validDelegationScopeWithExecution()
			tt.mutate(scope.ExecutionScope)
			if err := ValidateDelegationScope(scope); err == nil {
				t.Fatalf("expected invalid execution scope rejection for %s", tt.name)
			}
		})
	}
}

func TestDelegationExecutionSubjectShapeAndLegacyGrantAuthorization(t *testing.T) {
	t.Run("subject with gated bridge is valid fact but cannot be authorized by ceiling", func(t *testing.T) {
		scope := validDelegationScopeWithExecution()
		subject := validDelegationExecutionSubject()
		subject.NetworkPolicy = ExecutionNetworkPolicyGatedBridge

		if err := ValidateDelegationExecutionSubject(subject); err != nil {
			t.Fatalf("gated_bridge subject is a valid server-derived fact: %v", err)
		}
		if err := DelegationExecutionWithinScope(scope, subject); err == nil {
			t.Fatal("expected gated_bridge subject authorization rejection")
		}
	})

	t.Run("empty subject fields are rejected", func(t *testing.T) {
		subject := validDelegationExecutionSubject()
		subject.ProviderRef = ""
		if err := ValidateDelegationExecutionSubject(subject); err == nil {
			t.Fatal("expected empty subject provider_ref rejection")
		}
	})

	t.Run("legacy grant cannot authorize execution subject", func(t *testing.T) {
		scope := &DelegationScope{
			TaskTypes:   []string{"stage"},
			RiskCeiling: RiskMedium,
			Quota:       DelegationQuota{PerDay: 3},
		}
		if err := DelegationExecutionWithinScope(scope, validDelegationExecutionSubject()); err == nil {
			t.Fatal("expected legacy scope to reject execution subject")
		}
	})

	t.Run("execution scope cannot authorize nil execution subject", func(t *testing.T) {
		if err := DelegationExecutionWithinScope(validDelegationScopeWithExecution(), nil); err == nil {
			t.Fatal("expected nil execution subject rejection")
		}
	})
}

func validDelegationScopeWithExecution() *DelegationScope {
	return &DelegationScope{
		TaskTypes:   []string{"stage"},
		RiskCeiling: RiskMedium,
		Quota:       DelegationQuota{PerDay: 3},
		ExecutionScope: &DelegationExecutionScope{
			CapabilityRefs:       []string{"capability:writer", "capability:test"},
			WorkrootRef:          "workroot:repo-001",
			ProviderRefs:         []string{"provider:local-exec", "provider:codex"},
			SandboxProfileRef:    "sandbox:workspace-write",
			NetworkPolicyCeiling: ExecutionNetworkPolicyEgressModelOnly,
			MaxRuns:              2,
			MaxDurationSeconds:   300,
		},
	}
}

func validDelegationExecutionSubject() *DelegationExecutionSubject {
	return &DelegationExecutionSubject{
		CapabilityRef:           "capability:writer",
		WorkrootRef:             "workroot:repo-001",
		ProviderRef:             "provider:local-exec",
		SandboxProfileRef:       "sandbox:workspace-write",
		NetworkPolicy:           ExecutionNetworkPolicyEgressModelOnly,
		ConsumedRuns:            2,
		ConsumedDurationSeconds: 300,
	}
}

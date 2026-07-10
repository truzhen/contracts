package base

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
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
		{name: "negative max runs", mutate: func(s *DelegationExecutionScope) { s.MaxRuns = -1 }},
		{name: "zero max duration", mutate: func(s *DelegationExecutionScope) { s.MaxDurationSeconds = 0 }},
		{name: "negative max duration", mutate: func(s *DelegationExecutionScope) { s.MaxDurationSeconds = -1 }},
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

	t.Run("invalid subject fields are rejected", func(t *testing.T) {
		tests := []struct {
			name   string
			mutate func(*DelegationExecutionSubject)
		}{
			{name: "empty capability", mutate: func(s *DelegationExecutionSubject) { s.CapabilityRef = "" }},
			{name: "empty workroot", mutate: func(s *DelegationExecutionSubject) { s.WorkrootRef = "" }},
			{name: "empty provider", mutate: func(s *DelegationExecutionSubject) { s.ProviderRef = "" }},
			{name: "empty sandbox", mutate: func(s *DelegationExecutionSubject) { s.SandboxProfileRef = "" }},
			{name: "unknown network policy", mutate: func(s *DelegationExecutionSubject) { s.NetworkPolicy = ExecutionNetworkPolicy("unknown") }},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				subject := validDelegationExecutionSubject()
				tt.mutate(subject)
				if err := ValidateDelegationExecutionSubject(subject); err == nil {
					t.Fatalf("expected invalid execution subject rejection for %s", tt.name)
				}
			})
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

func TestDelegationGrantWithinScopeChecksParentBeforeExecution(t *testing.T) {
	t.Run("valid complete grant and subject", func(t *testing.T) {
		if err := DelegationGrantWithinScope(validOwnerDelegationGrantWithExecution(), validDelegationSubjectWithExecution(), delegationEvaluationTime()); err != nil {
			t.Fatalf("DelegationGrantWithinScope: %v", err)
		}
	})

	t.Run("parent rejection is evaluated before execution", func(t *testing.T) {
		subject := validDelegationSubjectWithExecution()
		subject.TaskType = "other"
		subject.Execution.CapabilityRef = "capability:other"
		err := DelegationGrantWithinScope(validOwnerDelegationGrantWithExecution(), subject, delegationEvaluationTime())
		if err == nil || !strings.Contains(err.Error(), "task_type") {
			t.Fatalf("expected parent task_type rejection before execution rejection, got %v", err)
		}
	})

	tests := []struct {
		name   string
		mutate func(*OwnerDelegationGrant, *DelegationSubject)
	}{
		{name: "task type", mutate: func(_ *OwnerDelegationGrant, s *DelegationSubject) { s.TaskType = "other" }},
		{name: "risk ceiling", mutate: func(g *OwnerDelegationGrant, s *DelegationSubject) {
			g.Scope.RiskCeiling = RiskLow
			s.RiskLevel = RiskMedium
		}},
		{name: "transaction", mutate: func(_ *OwnerDelegationGrant, s *DelegationSubject) { s.TransactionRef = "transaction:other" }},
		{name: "pack", mutate: func(_ *OwnerDelegationGrant, s *DelegationSubject) { s.PackRef = "pack:other" }},
		{name: "amount", mutate: func(_ *OwnerDelegationGrant, s *DelegationSubject) { s.AmountCents = 1001 }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grant := validOwnerDelegationGrantWithExecution()
			subject := validDelegationSubjectWithExecution()
			tt.mutate(grant, subject)
			if err := DelegationGrantWithinScope(grant, subject, delegationEvaluationTime()); err == nil {
				t.Fatalf("expected parent %s boundary rejection", tt.name)
			}
		})
	}
}

func TestDelegationGrantWithinScopePreservesRiskHardFloor(t *testing.T) {
	for _, risk := range []RiskClass{RiskHigh, RiskCritical} {
		t.Run(string(risk)+" subject", func(t *testing.T) {
			grant := validOwnerDelegationGrantWithExecution()
			subject := validDelegationSubjectWithExecution()
			subject.RiskLevel = risk
			if err := DelegationGrantWithinScope(grant, subject, delegationEvaluationTime()); err == nil {
				t.Fatalf("expected %s subject hard-floor rejection", risk)
			}
		})

		t.Run(string(risk)+" grant ceiling", func(t *testing.T) {
			grant := validOwnerDelegationGrantWithExecution()
			grant.Scope.RiskCeiling = risk
			if err := DelegationGrantWithinScope(grant, validDelegationSubjectWithExecution(), delegationEvaluationTime()); err == nil {
				t.Fatalf("expected %s grant ceiling hard-floor rejection", risk)
			}
		})
	}
}

func TestDelegationGrantWithinScopeRejectsInvalidParentSubject(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*DelegationSubject)
	}{
		{name: "empty candidate ref", mutate: func(s *DelegationSubject) { s.CandidateRef = "" }},
		{name: "empty transaction ref", mutate: func(s *DelegationSubject) { s.TransactionRef = "" }},
		{name: "empty task type", mutate: func(s *DelegationSubject) { s.TaskType = "" }},
		{name: "unknown risk", mutate: func(s *DelegationSubject) { s.RiskLevel = RiskClass("unknown") }},
		{name: "negative amount", mutate: func(s *DelegationSubject) { s.AmountCents = -1 }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subject := validDelegationSubjectWithExecution()
			tt.mutate(subject)
			if err := DelegationGrantWithinScope(validOwnerDelegationGrantWithExecution(), subject, delegationEvaluationTime()); err == nil {
				t.Fatalf("expected invalid parent subject rejection for %s", tt.name)
			}
		})
	}
}

func TestDelegationExecutionNetworkPolicyUsesPermissionOrdering(t *testing.T) {
	tests := []struct {
		name    string
		ceiling ExecutionNetworkPolicy
		policy  ExecutionNetworkPolicy
		wantErr bool
	}{
		{name: "none accepts none", ceiling: ExecutionNetworkPolicyNone, policy: ExecutionNetworkPolicyNone},
		{name: "none rejects model egress", ceiling: ExecutionNetworkPolicyNone, policy: ExecutionNetworkPolicyEgressModelOnly, wantErr: true},
		{name: "model egress accepts none", ceiling: ExecutionNetworkPolicyEgressModelOnly, policy: ExecutionNetworkPolicyNone},
		{name: "model egress accepts model egress", ceiling: ExecutionNetworkPolicyEgressModelOnly, policy: ExecutionNetworkPolicyEgressModelOnly},
		{name: "model egress rejects gated bridge", ceiling: ExecutionNetworkPolicyEgressModelOnly, policy: ExecutionNetworkPolicyGatedBridge, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := validDelegationScopeWithExecution()
			scope.ExecutionScope.NetworkPolicyCeiling = tt.ceiling
			subject := validDelegationExecutionSubject()
			subject.NetworkPolicy = tt.policy
			err := DelegationExecutionWithinScope(scope, subject)
			if tt.wantErr && err == nil {
				t.Fatal("expected network policy rejection")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected network policy authorization: %v", err)
			}
		})
	}
}

func TestDelegationExecutionSubjectRequiresPositiveReservedCumulativeBudget(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*DelegationExecutionSubject)
	}{
		{name: "zero runs", mutate: func(s *DelegationExecutionSubject) { s.ConsumedRuns = 0 }},
		{name: "negative runs", mutate: func(s *DelegationExecutionSubject) { s.ConsumedRuns = -1 }},
		{name: "runs over scope", mutate: func(s *DelegationExecutionSubject) { s.ConsumedRuns = 3 }},
		{name: "zero duration", mutate: func(s *DelegationExecutionSubject) { s.ConsumedDurationSeconds = 0 }},
		{name: "negative duration", mutate: func(s *DelegationExecutionSubject) { s.ConsumedDurationSeconds = -1 }},
		{name: "duration over scope", mutate: func(s *DelegationExecutionSubject) { s.ConsumedDurationSeconds = 301 }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subject := validDelegationExecutionSubject()
			tt.mutate(subject)
			if err := DelegationExecutionWithinScope(validDelegationScopeWithExecution(), subject); err == nil {
				t.Fatalf("expected reserved cumulative budget rejection for %s", tt.name)
			}
		})
	}
}

func TestDelegationGrantWithinScopeLegacyGrantCannotAuthorizeExecution(t *testing.T) {
	grant := validOwnerDelegationGrantWithExecution()
	grant.Scope.ExecutionScope = nil
	if err := DelegationGrantWithinScope(grant, validDelegationSubjectWithExecution(), delegationEvaluationTime()); err == nil {
		t.Fatal("expected legacy grant to reject complete execution subject")
	}
}

func TestDelegationGrantWithinScopeRequiresActiveUnexpiredGrant(t *testing.T) {
	t.Run("active revocable grant with future expiry passes", func(t *testing.T) {
		grant := validOwnerDelegationGrantWithExecution()
		grant.ExpiresAt = delegationEvaluationTime().Add(time.Second)
		if err := DelegationGrantWithinScope(grant, validDelegationSubjectWithExecution(), delegationEvaluationTime()); err != nil {
			t.Fatalf("expected active future grant authorization: %v", err)
		}
	})

	for _, status := range []DelegationGrantStatus{
		DelegationGrantRevoked,
		DelegationGrantExpired,
		DelegationGrantSuspendedByEmergencyStop,
	} {
		t.Run(string(status), func(t *testing.T) {
			grant := validOwnerDelegationGrantWithExecution()
			grant.Status = status
			if err := DelegationGrantWithinScope(grant, validDelegationSubjectWithExecution(), delegationEvaluationTime()); err == nil {
				t.Fatalf("expected %s grant rejection", status)
			}
		})
	}

	for _, tc := range []struct {
		name      string
		expiresAt time.Time
	}{
		{name: "expires at evaluation time", expiresAt: delegationEvaluationTime()},
		{name: "expires before evaluation time", expiresAt: delegationEvaluationTime().Add(-time.Nanosecond)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			grant := validOwnerDelegationGrantWithExecution()
			grant.ExpiresAt = tc.expiresAt
			if err := DelegationGrantWithinScope(grant, validDelegationSubjectWithExecution(), delegationEvaluationTime()); err == nil {
				t.Fatal("expected expired-at-evaluation grant rejection")
			}
		})
	}

	t.Run("zero evaluation time", func(t *testing.T) {
		if err := DelegationGrantWithinScope(validOwnerDelegationGrantWithExecution(), validDelegationSubjectWithExecution(), time.Time{}); err == nil {
			t.Fatal("expected zero evaluation time rejection")
		}
	})
}

func TestDelegationGrantWithinScopeRejectsAuthorizationHardDenies(t *testing.T) {
	for _, sideEffect := range AuthorizationHardDenies() {
		t.Run(string(sideEffect), func(t *testing.T) {
			subject := delegationSubjectWithServerFacts(t, sideEffect, 1, "2026-07-10")
			if err := DelegationGrantWithinScope(validOwnerDelegationGrantWithExecution(), subject, delegationEvaluationTime()); err == nil {
				t.Fatalf("expected hard-denied side effect %q rejection", sideEffect)
			}
		})
	}

	t.Run("unknown side effect", func(t *testing.T) {
		subject := delegationSubjectWithServerFacts(t, SideEffectClass("custom"), 1, "2026-07-10")
		if err := DelegationGrantWithinScope(validOwnerDelegationGrantWithExecution(), subject, delegationEvaluationTime()); err == nil {
			t.Fatal("expected unknown side effect rejection")
		}
	})

	t.Run("external send", func(t *testing.T) {
		subject := delegationSubjectWithServerFacts(t, SideEffectExternalSend, 1, "2026-07-10")
		if err := DelegationGrantWithinScope(validOwnerDelegationGrantWithExecution(), subject, delegationEvaluationTime()); err == nil {
			t.Fatal("expected external_send rejection")
		}
	})
}

func TestDelegationGrantWithinScopeEnforcesDailyDecisionQuota(t *testing.T) {
	t.Run("valid reserved decision passes", func(t *testing.T) {
		subject := delegationSubjectWithServerFacts(t, SideEffectLocalFileWrite, 3, "2026-07-10")
		if err := DelegationGrantWithinScope(validOwnerDelegationGrantWithExecution(), subject, delegationEvaluationTime()); err != nil {
			t.Fatalf("expected decision within daily quota: %v", err)
		}
	})

	for _, tc := range []struct {
		name      string
		consumed  int
		quotaDate string
	}{
		{name: "missing reservation", consumed: 0, quotaDate: "2026-07-10"},
		{name: "quota exceeded", consumed: 4, quotaDate: "2026-07-10"},
		{name: "wrong quota date", consumed: 1, quotaDate: "2026-07-09"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			subject := delegationSubjectWithServerFacts(t, SideEffectLocalFileWrite, tc.consumed, tc.quotaDate)
			if err := DelegationGrantWithinScope(validOwnerDelegationGrantWithExecution(), subject, delegationEvaluationTime()); err == nil {
				t.Fatalf("expected quota fact rejection: consumed=%d date=%q", tc.consumed, tc.quotaDate)
			}
		})
	}
}

func TestValidateOwnerDelegationGrantPreservesLegacyRevocableFalse(t *testing.T) {
	t.Run("revocable true", func(t *testing.T) {
		if err := ValidateOwnerDelegationGrant(validOwnerDelegationGrantWithExecution()); err != nil {
			t.Fatalf("expected revocable grant validation: %v", err)
		}
	})

	t.Run("revocable false", func(t *testing.T) {
		grant := validOwnerDelegationGrantWithExecution()
		grant.Revocable = false
		if err := ValidateOwnerDelegationGrant(grant); err != nil {
			t.Fatalf("legacy structural validation must preserve revocable=false behavior: %v", err)
		}
	})
}

func TestDelegationGrantWithinScopeFailsClosedOnGrantHardInvariants(t *testing.T) {
	t.Run("nil grant", func(t *testing.T) {
		if err := DelegationGrantWithinScope(nil, validDelegationSubjectWithExecution(), delegationEvaluationTime()); err == nil {
			t.Fatal("expected nil grant rejection")
		}
	})

	tests := []struct {
		name   string
		mutate func(*OwnerDelegationGrant)
	}{
		{name: "empty grant id", mutate: func(g *OwnerDelegationGrant) { g.GrantID = "" }},
		{name: "empty owner decision ref", mutate: func(g *OwnerDelegationGrant) { g.OwnerDecisionRef = "" }},
		{name: "empty delegate ref", mutate: func(g *OwnerDelegationGrant) { g.DelegateRef = "" }},
		{name: "empty receipt ref", mutate: func(g *OwnerDelegationGrant) { g.ReceiptRef = "" }},
		{name: "zero expiry", mutate: func(g *OwnerDelegationGrant) { g.ExpiresAt = time.Time{} }},
		{name: "unknown status", mutate: func(g *OwnerDelegationGrant) { g.Status = DelegationGrantStatus("unknown") }},
		{name: "irrevocable", mutate: func(g *OwnerDelegationGrant) { g.Revocable = false }},
		{name: "zero daily quota", mutate: func(g *OwnerDelegationGrant) { g.Scope.Quota.PerDay = 0 }},
		{name: "negative daily quota", mutate: func(g *OwnerDelegationGrant) { g.Scope.Quota.PerDay = -1 }},
		{name: "negative amount limit", mutate: func(g *OwnerDelegationGrant) { g.Scope.AmountLimitCents = -1 }},
		{name: "zero execution max runs", mutate: func(g *OwnerDelegationGrant) { g.Scope.ExecutionScope.MaxRuns = 0 }},
		{name: "negative execution max runs", mutate: func(g *OwnerDelegationGrant) { g.Scope.ExecutionScope.MaxRuns = -1 }},
		{name: "zero execution max duration", mutate: func(g *OwnerDelegationGrant) { g.Scope.ExecutionScope.MaxDurationSeconds = 0 }},
		{name: "negative execution max duration", mutate: func(g *OwnerDelegationGrant) { g.Scope.ExecutionScope.MaxDurationSeconds = -1 }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grant := validOwnerDelegationGrantWithExecution()
			tt.mutate(grant)
			if err := DelegationGrantWithinScope(grant, validDelegationSubjectWithExecution(), delegationEvaluationTime()); err == nil {
				t.Fatalf("expected invalid grant rejection for %s", tt.name)
			}
		})
	}
}

func validDelegationScopeWithExecution() *DelegationScope {
	return &DelegationScope{
		TaskTypes:        []string{"stage"},
		RiskCeiling:      RiskMedium,
		TransactionRefs:  []string{"transaction:001"},
		PackRefs:         []string{"pack:001"},
		Quota:            DelegationQuota{PerDay: 3},
		AmountLimitCents: 1000,
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

func validOwnerDelegationGrantWithExecution() *OwnerDelegationGrant {
	return &OwnerDelegationGrant{
		GrantID:          "grant:execution-001",
		OwnerDecisionRef: "owner_decision:001",
		DelegateRef:      "agent://secretary_chief",
		Scope:            *validDelegationScopeWithExecution(),
		ExpiresAt:        time.Date(2026, 7, 11, 0, 0, 0, 0, time.UTC),
		Revocable:        true,
		Status:           DelegationGrantActive,
		ReceiptRef:       "receipt:grant-001",
	}
}

func delegationEvaluationTime() time.Time {
	return time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)
}

func validDelegationSubjectWithExecution() *DelegationSubject {
	return &DelegationSubject{
		CandidateRef:           "candidate:001",
		TransactionRef:         "transaction:001",
		TaskType:               "stage",
		RiskLevel:              RiskLow,
		PackRef:                "pack:001",
		AmountCents:            1000,
		SideEffectClass:        SideEffectLocalFileWrite,
		QuotaDate:              "2026-07-10",
		ConsumedDecisionsToday: 1,
		Execution:              validDelegationExecutionSubject(),
	}
}

func delegationSubjectWithServerFacts(t *testing.T, sideEffect SideEffectClass, consumed int, quotaDate string) *DelegationSubject {
	t.Helper()
	raw, err := json.Marshal(validDelegationSubjectWithExecution())
	if err != nil {
		t.Fatalf("marshal delegation subject: %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatalf("decode delegation subject payload: %v", err)
	}
	payload["side_effect_class"] = sideEffect
	payload["consumed_decisions_today"] = consumed
	payload["quota_date"] = quotaDate
	raw, err = json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal delegation subject payload: %v", err)
	}
	var subject DelegationSubject
	if err := json.Unmarshal(raw, &subject); err != nil {
		t.Fatalf("unmarshal delegation subject with server facts: %v", err)
	}
	return &subject
}

package base

import (
	"strings"
	"testing"
	"time"
)

var fixedBaseTime = time.Date(2026, 7, 2, 10, 30, 0, 123, time.UTC)

func TestNewPolicySnapshotDerivesStableBaseIssuedRef(t *testing.T) {
	kernel := PolicyKernel{KernelRef: "policy_kernel:base", Version: "v1"}
	policySet := GatePolicySet{
		PolicySetRef: "policy_set:owner",
		Version:      "2026-07-02",
		CreatedAt:    fixedBaseTime,
	}

	first, err := NewPolicySnapshot(kernel, policySet, fixedBaseTime)
	if err != nil {
		t.Fatalf("NewPolicySnapshot first: %v", err)
	}
	second, err := NewPolicySnapshot(kernel, policySet, fixedBaseTime)
	if err != nil {
		t.Fatalf("NewPolicySnapshot second: %v", err)
	}

	if first.SnapshotRef == "" || !strings.HasPrefix(first.SnapshotRef, "policy_snapshot_") {
		t.Fatalf("unexpected snapshot ref: %q", first.SnapshotRef)
	}
	if first.SnapshotRef != second.SnapshotRef || first.SnapshotHash != second.SnapshotHash {
		t.Fatalf("snapshot ref/hash must be stable for explicit inputs: first=%+v second=%+v", first, second)
	}
	if first.IssuedBy != PolicySnapshotIssuerBaseKernel {
		t.Fatalf("snapshot must be Base issued, got %q", first.IssuedBy)
	}
	if err := ValidatePolicySnapshot(first); err != nil {
		t.Fatalf("ValidatePolicySnapshot: %v", err)
	}
}

func TestGateIntakeDecisionReceiptAndGrantKeepStableRefs(t *testing.T) {
	req := validGateRequest()
	snapshot := validPolicySnapshot(t)

	intake, err := NewCandidateIntakeEnvelope(req, snapshot.SnapshotRef, fixedBaseTime)
	if err != nil {
		t.Fatalf("NewCandidateIntakeEnvelope: %v", err)
	}
	if intake.IntakeRef == "" || !strings.HasPrefix(intake.IntakeRef, "candidate_intake_") {
		t.Fatalf("unexpected intake ref: %q", intake.IntakeRef)
	}
	if intake.CandidateRef != req.CandidateRef {
		t.Fatalf("request candidate override must win: got %q want %q", intake.CandidateRef, req.CandidateRef)
	}
	if got, want := intake.EvidenceRefs, []string{"evidence:owner-click", "intent_event:001"}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("evidence refs must dedupe while preserving first occurrence: got %v want %v", got, want)
	}

	decision, err := NewGateDecision(GateDecisionAllow, "owner approved", req, snapshot, nil, fixedBaseTime)
	if err != nil {
		t.Fatalf("NewGateDecision: %v", err)
	}
	if decision.DecisionRef == "" || !strings.HasPrefix(decision.DecisionRef, "gate_decision_") {
		t.Fatalf("unexpected decision ref: %q", decision.DecisionRef)
	}
	if len(decision.GateTrace) != 1 || decision.GateTrace[0].PolicySnapshotHash != snapshot.SnapshotHash {
		t.Fatalf("default trace must bind policy snapshot hash: %+v", decision.GateTrace)
	}

	receipt, err := NewGateReceiptCandidate(decision, req)
	if err != nil {
		t.Fatalf("NewGateReceiptCandidate: %v", err)
	}
	if receipt.ReceiptID == "" || !strings.HasPrefix(receipt.ReceiptID, "gate_receipt_candidate_") {
		t.Fatalf("unexpected receipt ref: %q", receipt.ReceiptID)
	}
	if receipt.PolicySnapshot != snapshot.SnapshotHash {
		t.Fatalf("receipt must bind decision policy hash: got %q want %q", receipt.PolicySnapshot, snapshot.SnapshotHash)
	}

	decision.ReceiptCandidate = receipt
	grant, err := NewBaseFormalizationGrant(decision, "FormalTask", "task.write", fixedBaseTime.Add(time.Hour))
	if err != nil {
		t.Fatalf("NewBaseFormalizationGrant: %v", err)
	}
	if grant.GrantRef != "base_grant_"+decision.DecisionRef {
		t.Fatalf("grant ref must derive from decision ref: got %q", grant.GrantRef)
	}
	if grant.ReceiptCandidateRef != receipt.ReceiptID {
		t.Fatalf("grant must bind receipt candidate: got %q want %q", grant.ReceiptCandidateRef, receipt.ReceiptID)
	}
}

func TestValidationHelpersRejectUnsafeOrIncompleteInputs(t *testing.T) {
	t.Run("candidate must stay candidate only and non formal", func(t *testing.T) {
		env := NewGateCandidateEnvelope("candidate:001", "transaction:001", "intent_event:001", RiskLow, SideEffectReadOnly)
		env.NonFormal = false
		if err := ValidateGateCandidateEnvelope(env); err == nil {
			t.Fatal("expected non_formal violation")
		}
	})

	t.Run("owner decision must target the candidate", func(t *testing.T) {
		decision := &OwnerDecision{
			TargetCandidateRef: "candidate:other",
			ActorContext:       validActor(),
			PolicySnapshotRef:  "policy_snapshot:001",
			DecidedAt:          fixedBaseTime,
		}
		if err := ValidateOwnerDecisionForCandidate(decision, "candidate:001"); err == nil {
			t.Fatal("expected owner decision target mismatch")
		}
	})

	t.Run("adapter request rejects raw secret markers", func(t *testing.T) {
		req := validBaseAdapterRequest()
		req.EvidenceRefs = append(req.EvidenceRefs, "evidence:raw_secret_leak")
		if err := ValidateBaseAdapterRequest(req, AdapterKindExecution); err == nil {
			t.Fatal("expected secret-ish adapter request rejection")
		}
	})

	t.Run("object directory boundary fixes truth source", func(t *testing.T) {
		ref := &ObjectDirectoryBoundaryRef{
			ObjectDirectoryEntryRef: "object:001",
			ObjectTruthSourceModule: "12-frontend-shell",
			TransactionRef:          "transaction:001",
		}
		if err := ValidateObjectDirectoryBoundaryRef(ref); err == nil {
			t.Fatal("expected object truth source rejection")
		}
	})
}

func TestStableRefIsDeterministicAndPrefixScoped(t *testing.T) {
	first := stableRef("gate_decision", "candidate:001", "transaction:001")
	second := stableRef("gate_decision", "candidate:001", "transaction:001")
	otherPrefix := stableRef("gate_receipt_candidate", "candidate:001", "transaction:001")

	if first != second {
		t.Fatalf("stableRef must be deterministic: %q != %q", first, second)
	}
	if !strings.HasPrefix(first, "gate_decision_") {
		t.Fatalf("stableRef must include prefix: %q", first)
	}
	if first == otherPrefix {
		t.Fatalf("stableRef prefixes must stay scoped: %q", first)
	}
}

func validPolicySnapshot(t *testing.T) *PolicySnapshot {
	t.Helper()
	snapshot, err := NewPolicySnapshot(
		PolicyKernel{KernelRef: "policy_kernel:base", Version: "v1"},
		GatePolicySet{PolicySetRef: "policy_set:owner", Version: "2026-07-02", CreatedAt: fixedBaseTime},
		fixedBaseTime,
	)
	if err != nil {
		t.Fatalf("validPolicySnapshot: %v", err)
	}
	return snapshot
}

func validGateRequest() *GateRequest {
	candidate := NewGateCandidateEnvelope("candidate:001", "transaction:001", "intent_event:001", RiskHigh, SideEffectExternalSend)
	return &GateRequest{
		Candidate:          candidate,
		ActorContext:       validActor(),
		TransactionRef:     "transaction:override",
		CandidateRef:       "candidate:override",
		IntentEventRef:     "intent_event:001",
		EvidenceRefs:       []string{"evidence:owner-click", "intent_event:001"},
		ReceiptRequirement: "receipt_required",
		PolicySnapshotRef:  "policy_snapshot:001",
	}
}

func validActor() *ActorContext {
	return &ActorContext{
		OwnerProfile:      "owner:001",
		LocalIdentityRef:  "local_identity:001",
		DeviceIdentityRef: "device:001",
	}
}

func validBaseAdapterRequest() BaseGateAdapterRequest {
	return BaseGateAdapterRequest{
		AdapterRequestRef:   "adapter_request:001",
		AdapterKind:         AdapterKindExecution,
		CandidateRef:        "candidate:001",
		TransactionRef:      "transaction:001",
		DecisionRef:         "gate_decision:001",
		ReceiptCandidateRef: "gate_receipt_candidate:001",
		PolicySnapshotHash:  "policy_hash:001",
		EvidenceRefs:        []string{"evidence:001"},
		ActorContext:        validActor(),
		RequestedAt:         fixedBaseTime,
	}
}

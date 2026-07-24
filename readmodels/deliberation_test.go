package readmodels

import (
	"reflect"
	"strings"
	"testing"
)

func validDeliberationGrantReadModel() DeliberationAutomationGrantReadModel {
	return DeliberationAutomationGrantReadModel{
		GrantRef:             "delegation_grant://deliberation/current-turn/1",
		Mode:                 DeliberationAutomationCurrentTurnAuto,
		SessionRef:           "deliberation_session://owner/1",
		TurnRef:              "deliberation_turn://owner/1/1",
		TransactionRef:       "transaction://owner/1",
		ProviderAdapterRefs:  []string{"provider_adapter://fixture/v1"},
		AllowedOperations:    []DeliberationAllowedOperation{DeliberationOperationFillPrompt, DeliberationOperationSubmitPrompt, DeliberationOperationObserveResponse, DeliberationOperationCaptureResponse, DeliberationOperationRequestSynthesis},
		DispatchOnConfirm:    true,
		QuestionSHA256:       strings.Repeat("a", 64),
		IdleExpiresAt:        "2026-07-24T00:10:00Z",
		AbsoluteExpiresAt:    "2026-07-24T02:00:00Z",
		MaxTurns:             1,
		RemainingTurns:       1,
		MaxDispatchesPerLane: 1,
		Status:               DeliberationGrantPrepared,
		DecisionRef:          "decision://base/1",
		PolicySnapshotRef:    "policy_snapshot://base/1",
		EvidenceRefs:         []string{"evidence://base/1"},
		CandidateOnly:        true,
		NonFormal:            true,
		CreatedAt:            "2026-07-24T00:00:00Z",
		UpdatedAt:            "2026-07-24T00:00:00Z",
	}
}

func TestValidateDeliberationAutomationGrantReadModel(t *testing.T) {
	if err := ValidateDeliberationAutomationGrantReadModel(validDeliberationGrantReadModel()); err != nil {
		t.Fatalf("valid current-turn automation grant rejected: %v", err)
	}

	tests := []struct {
		name   string
		mutate func(*DeliberationAutomationGrantReadModel)
	}{
		{
			name: "current turn missing turn ref",
			mutate: func(grant *DeliberationAutomationGrantReadModel) {
				grant.TurnRef = ""
			},
		},
		{
			name: "current turn missing question hash",
			mutate: func(grant *DeliberationAutomationGrantReadModel) {
				grant.QuestionSHA256 = ""
			},
		},
		{
			name: "unknown operation",
			mutate: func(grant *DeliberationAutomationGrantReadModel) {
				grant.AllowedOperations = []DeliberationAllowedOperation{"open_any_page"}
			},
		},
		{
			name: "multiple dispatches per lane",
			mutate: func(grant *DeliberationAutomationGrantReadModel) {
				grant.MaxDispatchesPerLane = 2
			},
		},
		{
			name: "confirmation that does not dispatch",
			mutate: func(grant *DeliberationAutomationGrantReadModel) {
				grant.DispatchOnConfirm = false
			},
		},
		{
			name: "formal grant projection",
			mutate: func(grant *DeliberationAutomationGrantReadModel) {
				grant.CandidateOnly = false
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			grant := validDeliberationGrantReadModel()
			test.mutate(&grant)
			if err := ValidateDeliberationAutomationGrantReadModel(grant); err == nil {
				t.Fatal("invalid automation grant was accepted")
			}
		})
	}
}

func TestDeliberationAutomationGrantSessionModeAllowsNoTurnHash(t *testing.T) {
	grant := validDeliberationGrantReadModel()
	grant.Mode = DeliberationAutomationSessionAuto
	grant.TurnRef = ""
	grant.QuestionSHA256 = ""
	grant.MaxTurns = 20
	grant.RemainingTurns = 20
	if err := ValidateDeliberationAutomationGrantReadModel(grant); err != nil {
		t.Fatalf("valid session automation grant rejected: %v", err)
	}
}

func TestDeliberationReadModelsDoNotContainRawConversationOrAuthorityInputs(t *testing.T) {
	for _, value := range []any{
		DeliberationSessionReadModel{},
		DeliberationTurnReadModel{},
		DeliberationProviderLaneReadModel{},
		DeliberationAutomationGrantReadModel{},
	} {
		typeOf := reflect.TypeOf(value)
		for index := 0; index < typeOf.NumField(); index++ {
			field := typeOf.Field(index)
			jsonName := strings.Split(field.Tag.Get("json"), ",")[0]
			switch jsonName {
			case "question", "prompt", "response", "response_text", "approved", "owner_action_evidence_ref", "run_id", "nonce":
				t.Fatalf("%s must not expose raw conversation or client authority input %q", typeOf.Name(), jsonName)
			}
		}
	}
}

func TestDeliberationReadModelsCarryRequiredTraceabilityFields(t *testing.T) {
	assertJSONFields := func(value any, names ...string) {
		t.Helper()
		fields := map[string]bool{}
		typeOf := reflect.TypeOf(value)
		for index := 0; index < typeOf.NumField(); index++ {
			name := strings.Split(typeOf.Field(index).Tag.Get("json"), ",")[0]
			fields[name] = true
		}
		for _, name := range names {
			if !fields[name] {
				t.Fatalf("%s must carry %q", typeOf.Name(), name)
			}
		}
	}
	assertJSONFields(DeliberationSessionReadModel{}, "source_event_id", "turns", "selected_provider_adapter_refs")
	assertJSONFields(DeliberationTurnReadModel{}, "source_event_id", "question_artifact_ref", "question_sha256", "provider_lanes")
	assertJSONFields(DeliberationProviderLaneReadModel{}, "adapter_digest", "release_eligibility", "runtime_readiness", "response_byte_count", "provider_message_ref_hash", "receipt_refs", "started_at", "completed_at", "expires_at")
	assertJSONFields(DeliberationAutomationGrantReadModel{}, "dispatch_on_confirm", "max_dispatches_per_lane", "decision_ref")
}

func TestDeliberationEnumValidationFailsClosed(t *testing.T) {
	if ValidDeliberationAutomationMode(DeliberationAutomationMode("ambient_auto")) {
		t.Fatal("unknown automation mode must be rejected")
	}
	if ValidDeliberationAllowedOperation(DeliberationAllowedOperation("navigate_anywhere")) {
		t.Fatal("unknown operation must be rejected")
	}
	if ValidDeliberationLaneStatus(DeliberationLaneStatus("success_without_receipt")) {
		t.Fatal("unknown lane status must be rejected")
	}
	if ValidDeliberationSessionStatus(DeliberationSessionStatus("silently_formal")) {
		t.Fatal("unknown session status must be rejected")
	}
	if ValidDeliberationTurnStatus(DeliberationTurnStatus("replayed_success")) {
		t.Fatal("unknown turn status must be rejected")
	}
	if ValidDeliberationProviderGroup(DeliberationProviderGroup("unreviewed_group")) {
		t.Fatal("unknown provider group must be rejected")
	}
	if ValidDeliberationAccessMode(DeliberationAccessMode("browser_agent")) {
		t.Fatal("unknown access mode must be rejected")
	}
	if ValidDeliberationReleaseEligibility(DeliberationReleaseEligibility("assumed_ready")) {
		t.Fatal("unknown release eligibility must be rejected")
	}
	if ValidDeliberationRuntimeReadiness(DeliberationRuntimeReadiness("optimistic")) {
		t.Fatal("unknown runtime readiness must be rejected")
	}
	if ValidDeliberationFailureCode(DeliberationFailureCode("ignored")) {
		t.Fatal("unknown failure code must be rejected")
	}
	if ValidDeliberationCaptureMode(DeliberationCaptureMode("copy_everything")) {
		t.Fatal("unknown capture mode must be rejected")
	}
}

func validDeliberationProviderLaneReadModel() DeliberationProviderLaneReadModel {
	return DeliberationProviderLaneReadModel{
		LaneRef:             "deliberation_lane://owner/1",
		SessionRef:          "deliberation_session://owner/1",
		TurnRef:             "deliberation_turn://owner/1/1",
		ProviderGroup:       DeliberationProviderGroupChinaCore,
		ProviderResourceRef: "provider_resource://owner/1",
		ProviderAdapterRef:  "provider_adapter://owner/1",
		AdapterVersion:      "1.0.0",
		AdapterDigest:       strings.Repeat("a", 64),
		AccessMode:          DeliberationAccessProviderAuthorizedWeb,
		ReleaseEligibility:  DeliberationReleaseReady,
		RuntimeReadiness:    DeliberationRuntimeReady,
		Status:              DeliberationLaneImported,
		ResponseArtifactRef: "artifact://owner/1",
		ResponseSHA256:      strings.Repeat("b", 64),
		ResponseByteCount:   12,
		ReceiptRefs:         []string{"receipt://owner/1"},
		CreatedAt:           "2026-07-24T00:00:00Z",
		UpdatedAt:           "2026-07-24T00:00:00Z",
		CandidateOnly:       true,
		NonFormal:           true,
	}
}

func TestValidateDeliberationProviderLaneReadModelFailsClosed(t *testing.T) {
	lane := validDeliberationProviderLaneReadModel()
	if err := ValidateDeliberationProviderLaneReadModel(lane); err != nil {
		t.Fatalf("valid imported lane rejected: %v", err)
	}
	lane.AdapterDigest = "not-a-sha"
	if err := ValidateDeliberationProviderLaneReadModel(lane); err == nil {
		t.Fatal("lane with invalid adapter digest was accepted")
	}
}

func TestValidateDeliberationTurnAndSessionReadModelsRejectCrossScopeData(t *testing.T) {
	turn := DeliberationTurnReadModel{
		TurnRef:             "deliberation_turn://owner/1/1",
		SessionRef:          "deliberation_session://owner/1",
		TransactionRef:      "transaction://owner/1",
		SourceEventID:       "intent_event://owner/1",
		Sequence:            1,
		QuestionArtifactRef: "artifact://owner/question/1",
		QuestionSHA256:      strings.Repeat("c", 64),
		Status:              DeliberationTurnReady,
		ProviderLanes:       []DeliberationProviderLaneReadModel{validDeliberationProviderLaneReadModel()},
		SelectedLaneCount:   1,
		ImportedLaneCount:   1,
		OccVersion:          1,
		CreatedAt:           "2026-07-24T00:00:00Z",
		UpdatedAt:           "2026-07-24T00:00:00Z",
		CandidateOnly:       true,
		NonFormal:           true,
	}
	if err := ValidateDeliberationTurnReadModel(turn); err != nil {
		t.Fatalf("valid turn rejected: %v", err)
	}
	session := DeliberationSessionReadModel{
		SessionRef:                  turn.SessionRef,
		TransactionRef:              turn.TransactionRef,
		SourceEventID:               turn.SourceEventID,
		OwnerRef:                    "owner://owner/1",
		Status:                      DeliberationSessionReady,
		AutomationMode:              DeliberationAutomationManual,
		Turns:                       []DeliberationTurnReadModel{turn},
		SelectedProviderAdapterRefs: []string{"provider_adapter://owner/1"},
		OccVersion:                  1,
		CreatedAt:                   "2026-07-24T00:00:00Z",
		UpdatedAt:                   "2026-07-24T00:00:00Z",
		CandidateOnly:               true,
		NonFormal:                   true,
	}
	if err := ValidateDeliberationSessionReadModel(session); err != nil {
		t.Fatalf("valid session rejected: %v", err)
	}
	session.Turns[0].SourceEventID = "intent_event://other/1"
	if err := ValidateDeliberationSessionReadModel(session); err == nil {
		t.Fatal("session accepted a turn from another source event")
	}
}

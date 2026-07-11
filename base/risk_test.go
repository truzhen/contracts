package base_test

// #11：GateCandidateEnvelope.DeclaredRiskTypes additive 守卫（同 DeclaredImpacts 契约：
// 声明非授权、缺席零影响、omitempty）。

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/truzhen/contracts/base"
	"github.com/truzhen/contracts/spines"
)

func TestGateEnvelopeDeclaredRiskTypesRoundtripAndOmitempty(t *testing.T) {
	bare, err := json.Marshal(base.GateCandidateEnvelope{EnvelopeID: "e1"})
	if err != nil {
		t.Fatalf("marshal bare: %v", err)
	}
	if strings.Contains(string(bare), "declared_risk_types") {
		t.Fatalf("未声明时 declared_risk_types 必须 omitempty：%s", bare)
	}

	env := base.GateCandidateEnvelope{
		EnvelopeID: "e1",
		DeclaredRiskTypes: []spines.DeclaredRiskType{{
			RiskTypeID:        "final_notice",
			TriggerActionType: "scene_flow_controlled_execute",
			EscalationPath:    spines.RiskEscalationOwnerGate,
		}},
	}
	raw, err := json.Marshal(env)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var back base.GateCandidateEnvelope
	if err := json.Unmarshal(raw, &back); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(back.DeclaredRiskTypes) != 1 || back.DeclaredRiskTypes[0].EscalationPath != spines.RiskEscalationOwnerGate {
		t.Fatalf("roundtrip 失真: %+v", back.DeclaredRiskTypes)
	}
}

package market_test

// 统一决策表 #11：PackRiskType 五件套守卫——schema↔Go 枚举同步 + roundtrip/omitempty。
// KWeaver RiskType 反例（规格-only 无消费者）的防线之一：形状在这里被测试钉死，
// 消费者（Base orchestrator risk-type floor）在 truzhenos 同轮接线。

import (
	"encoding/json"
	"strings"
	"testing"

	contracts "github.com/truzhen/contracts"
	"github.com/truzhen/contracts/market"
	"github.com/truzhen/contracts/spines"
)

func TestPackRiskTypeEscalationEnumMatchesSchema(t *testing.T) {
	var schema struct {
		Properties struct {
			RiskTypes struct {
				Items struct {
					Properties struct {
						EscalationPath struct {
							Enum []string `json:"enum"`
						} `json:"escalation_path"`
					} `json:"properties"`
					Required []string `json:"required"`
				} `json:"items"`
			} `json:"risk_types"`
		} `json:"properties"`
	}
	if err := json.Unmarshal(contracts.PackManifestSchemaJSON, &schema); err != nil {
		t.Fatalf("parse embedded pack-manifest schema: %v", err)
	}
	want := []string{string(spines.RiskEscalationNone), string(spines.RiskEscalationOwnerGate)}
	got := schema.Properties.RiskTypes.Items.Properties.EscalationPath.Enum
	if len(got) != len(want) {
		t.Fatalf("escalation_path enum count %d != spines constants %d: %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("enum[%d]: schema %q != spines %q", i, got[i], want[i])
		}
	}
	req := schema.Properties.RiskTypes.Items.Required
	for _, must := range []string{"risk_type_id", "definition", "trigger_action_types", "escalation_path"} {
		found := false
		for _, r := range req {
			if r == must {
				found = true
			}
		}
		if !found {
			t.Fatalf("schema items.required 缺 %s: %v", must, req)
		}
	}
}

func TestPackManifestRiskTypesRoundtripAndOmitempty(t *testing.T) {
	bare, err := json.Marshal(market.PackManifest{PackID: "p", Name: "n", Version: "1", Kind: "scene_pack", MinTruzhenVersion: "3"})
	if err != nil {
		t.Fatalf("marshal bare: %v", err)
	}
	if strings.Contains(string(bare), "risk_types") {
		t.Fatalf("未声明时 risk_types 必须 omitempty：%s", bare)
	}

	m := market.PackManifest{
		PackID: "p", Name: "n", Version: "1", Kind: "scene_pack", MinTruzhenVersion: "3",
		RiskTypes: []market.PackRiskType{{
			RiskTypeID:          "final_notice",
			Definition:          "对外正式处罚告知",
			TriggerActionTypes:  []string{"scene_flow_controlled_execute"},
			EvidenceRequirement: "现场证据链齐全",
			EscalationPath:      string(spines.RiskEscalationOwnerGate),
			Fallback:            "退回人工复核",
		}},
	}
	raw, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var back market.PackManifest
	if err := json.Unmarshal(raw, &back); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(back.RiskTypes) != 1 || back.RiskTypes[0].RiskTypeID != "final_notice" ||
		back.RiskTypes[0].EscalationPath != "owner_gate" ||
		len(back.RiskTypes[0].TriggerActionTypes) != 1 {
		t.Fatalf("roundtrip 失真: %+v", back.RiskTypes)
	}
	for _, tag := range []string{"\"risk_type_id\":", "\"definition\":", "\"trigger_action_types\":", "\"evidence_requirement\":", "\"escalation_path\":", "\"fallback\":"} {
		if !strings.Contains(string(raw), tag) {
			t.Fatalf("五件套 json tag 缺 %s: %s", tag, raw)
		}
	}
}

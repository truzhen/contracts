package base_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/truzhen/contracts/base"
	"github.com/truzhen/contracts/spines"
)

func TestGateCandidateEnvelopeDeclaredImpactsGoldenValues(t *testing.T) {
	env := base.GateCandidateEnvelope{
		EnvelopeID:      "envelope://impact-001",
		CandidateOnly:   true,
		NonFormal:       true,
		RiskClass:       base.RiskMedium,
		SideEffectClass: base.SideEffectLocalDraft,
		Payload:         map[string]string{"kind": "demo"},
		DeclaredImpacts: []spines.DeclaredImpact{
			{
				ObjectType:     "transaction_object",
				Operation:      spines.ImpactOperationDelete,
				AffectedFields: []string{"archive_flag"},
				Description:    "归档旧项目对象",
			},
		},
	}
	b, err := json.Marshal(env)
	if err != nil {
		t.Fatalf("marshal envelope: %v", err)
	}
	s := string(b)
	for _, want := range []string{
		`"declared_impacts":[{`,
		`"object_type":"transaction_object"`,
		`"operation":"delete"`,
		`"affected_fields":["archive_flag"]`,
	} {
		if !strings.Contains(s, want) {
			t.Fatalf("missing %s in %s", want, s)
		}
	}

	// omitempty：未声明影响的候选序列化后不得出现该键（保守档=缺省行为不变）。
	legacy := base.GateCandidateEnvelope{EnvelopeID: "envelope://legacy", CandidateOnly: true, NonFormal: true}
	lb, err := json.Marshal(legacy)
	if err != nil {
		t.Fatalf("marshal legacy envelope: %v", err)
	}
	if strings.Contains(string(lb), "declared_impacts") {
		t.Fatalf("legacy envelope must omit declared_impacts: %s", string(lb))
	}
}

// DeclaredImpact 的 object_ref 可选（执行前未必知道对象）；ActualEdit 的必填约束
// 由 receipts 侧 schema 测试守（object_ref 必填），两端语义不对称是设计而非疏漏。
func TestDeclaredImpactObjectRefOptional(t *testing.T) {
	d := spines.DeclaredImpact{ObjectType: "customer", Operation: spines.ImpactOperationCreate}
	b, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("marshal declared impact: %v", err)
	}
	if strings.Contains(string(b), "object_ref") {
		t.Fatalf("empty object_ref must be omitted: %s", string(b))
	}
}

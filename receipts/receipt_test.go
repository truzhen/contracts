package receipts_test

import (
	"encoding/json"
	"strings"
	"testing"

	contracts "github.com/truzhen/contracts"
	"github.com/truzhen/contracts/receipts"
	"github.com/truzhen/contracts/spines"
)

func TestReceiptEnvelopeActualEditsGoldenValues(t *testing.T) {
	env := receipts.ReceiptEnvelope{
		ReceiptID:    "receipt://impact-001",
		Sequence:     7,
		PreviousHash: "prev",
		PayloadHash:  "payload",
		Hash:         "hash",
		ActualEdits: []spines.ActualEdit{
			{
				ObjectType:     "transaction_object",
				Operation:      spines.ImpactOperationModify,
				ObjectRef:      "transaction://t-001",
				AffectedFields: []string{"status"},
			},
		},
	}
	b, err := json.Marshal(env)
	if err != nil {
		t.Fatalf("marshal receipt: %v", err)
	}
	s := string(b)
	for _, want := range []string{
		`"actual_edits":[{`,
		`"object_type":"transaction_object"`,
		`"operation":"modify"`,
		`"object_ref":"transaction://t-001"`,
		`"affected_fields":["status"]`,
	} {
		if !strings.Contains(s, want) {
			t.Fatalf("missing %s in %s", want, s)
		}
	}

	// omitempty：无 actual_edits 的历史回执序列化后不得出现该键。
	legacy := receipts.ReceiptEnvelope{ReceiptID: "receipt://legacy", Sequence: 1, PreviousHash: "p", PayloadHash: "pl", Hash: "h"}
	lb, err := json.Marshal(legacy)
	if err != nil {
		t.Fatalf("marshal legacy receipt: %v", err)
	}
	if strings.Contains(string(lb), "actual_edits") {
		t.Fatalf("legacy receipt must omit actual_edits: %s", string(lb))
	}
}

// 三操作词表 = 对象编辑域封顶（send/execute 归 SideEffectClass，Owner O-3 裁定不扩）。
func TestImpactOperationCoversObjectEditDomainOnly(t *testing.T) {
	ops := []spines.ImpactOperation{
		spines.ImpactOperationCreate,
		spines.ImpactOperationModify,
		spines.ImpactOperationDelete,
	}
	want := []string{"create", "modify", "delete"}
	for i, op := range ops {
		if string(op) != want[i] {
			t.Fatalf("op %d: got %q want %q", i, op, want[i])
		}
	}
}

// embed schema 的 actual_edits.items.operation.enum 必须与 Go 常量逐值一致
// （go-schema 一致性门 v1 只查顶层 kind/required，不查嵌套 enum，本测试补位）。
func TestReceiptActualEditsEnumMatchesSchema(t *testing.T) {
	var schema struct {
		Properties struct {
			ActualEdits struct {
				Items struct {
					Required   []string `json:"required"`
					Properties struct {
						Operation struct {
							Enum []string `json:"enum"`
						} `json:"operation"`
					} `json:"properties"`
				} `json:"items"`
			} `json:"actual_edits"`
		} `json:"properties"`
	}
	if err := json.Unmarshal(contracts.ReceiptEnvelopeSchemaJSON, &schema); err != nil {
		t.Fatalf("parse embedded receipt schema: %v", err)
	}
	wantEnum := []string{
		string(spines.ImpactOperationCreate),
		string(spines.ImpactOperationModify),
		string(spines.ImpactOperationDelete),
	}
	got := schema.Properties.ActualEdits.Items.Properties.Operation.Enum
	if len(got) != len(wantEnum) {
		t.Fatalf("schema enum %v != Go %v", got, wantEnum)
	}
	for i := range wantEnum {
		if got[i] != wantEnum[i] {
			t.Fatalf("enum[%d]: schema %q != Go %q", i, got[i], wantEnum[i])
		}
	}
	// ActualEdit 三必填：事实必须指到对象（object_ref 必填区别于 DeclaredImpact）。
	wantReq := []string{"object_type", "operation", "object_ref"}
	gotReq := schema.Properties.ActualEdits.Items.Required
	if len(gotReq) != len(wantReq) {
		t.Fatalf("items.required %v != %v", gotReq, wantReq)
	}
	for i := range wantReq {
		if gotReq[i] != wantReq[i] {
			t.Fatalf("required[%d]: %q != %q", i, gotReq[i], wantReq[i])
		}
	}
}

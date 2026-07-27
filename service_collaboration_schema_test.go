package contracts_test

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	contracts "github.com/truzhen/contracts"
	"github.com/truzhen/contracts/readmodels"
)

func TestServiceCollaborationSchemaIsClosedAndComplete(t *testing.T) {
	var document struct {
		OneOf []json.RawMessage `json:"oneOf"`
		Defs  map[string]struct {
			AdditionalProperties *bool                      `json:"additionalProperties"`
			Properties           map[string]json.RawMessage `json:"properties"`
			Required             []string                   `json:"required"`
		} `json:"$defs"`
	}
	if err := json.Unmarshal(contracts.ServiceCollaborationSchemaJSON, &document); err != nil {
		t.Fatal(err)
	}
	if len(document.OneOf) != 5 {
		t.Fatalf("schema must expose five consumable roots, got %d", len(document.OneOf))
	}
	for _, name := range []string{"service_collaboration_invite", "service_agreement_disclosure", "counterparty_decision_envelope", "daily_delivery_disclosure", "participant_scope_readmodel"} {
		shape, ok := document.Defs[name]
		if !ok || shape.AdditionalProperties == nil || *shape.AdditionalProperties || len(shape.Properties) == 0 || len(shape.Required) == 0 {
			t.Fatalf("missing or open shape %s", name)
		}
	}
	for name, value := range map[string]any{
		"service_collaboration_invite":   readmodels.ServiceCollaborationInvite{},
		"service_agreement_disclosure":   readmodels.ServiceAgreementDisclosure{},
		"counterparty_decision_envelope": readmodels.CounterpartyDecisionEnvelope{},
		"daily_delivery_disclosure":      readmodels.DailyDeliveryDisclosure{},
		"participant_scope_readmodel":    readmodels.ParticipantScopeReadModel{},
	} {
		shape := document.Defs[name]
		typ := reflect.TypeOf(value)
		for index := 0; index < typ.NumField(); index++ {
			field := strings.Split(typ.Field(index).Tag.Get("json"), ",")[0]
			if _, ok := shape.Properties[field]; !ok {
				t.Fatalf("%s missing Go JSON field %s", name, field)
			}
			found := false
			for _, required := range shape.Required {
				if required == field {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("%s must require Go JSON field %s", name, field)
			}
		}
		if len(shape.Properties) != typ.NumField() {
			t.Fatalf("%s contains schema-only properties", name)
		}
	}
	refs := map[string]bool{}
	for _, raw := range document.OneOf {
		var item struct {
			Ref string `json:"$ref"`
		}
		if err := json.Unmarshal(raw, &item); err != nil {
			t.Fatal(err)
		}
		refs[item.Ref] = true
	}
	for _, name := range []string{"service_collaboration_invite", "service_agreement_disclosure", "counterparty_decision_envelope", "daily_delivery_disclosure", "participant_scope_readmodel"} {
		if !refs["#/$defs/"+name] {
			t.Fatalf("oneOf missing %s", name)
		}
	}
}

func TestServiceCollaborationSchemaAdversarialConstraints(t *testing.T) {
	var document struct {
		Defs map[string]struct {
			Properties map[string]json.RawMessage `json:"properties"`
		} `json:"$defs"`
	}
	if err := json.Unmarshal(contracts.ServiceCollaborationSchemaJSON, &document); err != nil {
		t.Fatal(err)
	}
	assertConst := func(def, field string, want bool) {
		var rule struct {
			Const bool `json:"const"`
		}
		if err := json.Unmarshal(document.Defs[def].Properties[field], &rule); err != nil || rule.Const != want {
			t.Fatalf("%s.%s const mismatch", def, field)
		}
	}
	assertConst("service_collaboration_invite", "candidate_only", true)
	assertConst("service_collaboration_invite", "produces_owner_decision", false)
	assertConst("counterparty_decision_envelope", "candidate_only", true)
	assertConst("counterparty_decision_envelope", "produces_owner_decision", false)
	assertConst("service_agreement_disclosure", "read_model_only", true)
	assertConst("service_agreement_disclosure", "formal_agreement", false)
	assertConst("daily_delivery_disclosure", "read_model_only", true)
	assertConst("participant_scope_readmodel", "read_model_only", true)
	assertConst("participant_scope_readmodel", "can_produce_owner_decision", false)
	var decision struct {
		Enum []string `json:"enum"`
	}
	_ = json.Unmarshal(document.Defs["counterparty_decision_envelope"].Properties["decision"], &decision)
	if strings.Join(decision.Enum, ",") != "accepted,revision_requested,rejected" {
		t.Fatal("counterparty decision enum must exclude Owner semantics")
	}
	for _, pair := range [][2]string{{"daily_delivery_disclosure", "item_refs"}, {"participant_scope_readmodel", "allowed_transaction_refs"}} {
		var rule struct {
			Min    int  `json:"minItems"`
			Max    int  `json:"maxItems"`
			Unique bool `json:"uniqueItems"`
		}
		_ = json.Unmarshal(document.Defs[pair[0]].Properties[pair[1]], &rule)
		if !rule.Unique || rule.Min != 5 && pair[0] == "daily_delivery_disclosure" || pair[0] == "daily_delivery_disclosure" && rule.Max != 5 || pair[0] == "participant_scope_readmodel" && (rule.Min != 1 || rule.Max != 1) {
			t.Fatalf("%s.%s cardinality/uniqueness mismatch", pair[0], pair[1])
		}
	}
}

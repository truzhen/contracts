package candidates

import (
	"encoding/json"
	"strings"
	"testing"
)

// task_kind is an additive field on the cross-repo TaskCandidate contract: it
// declares the 07 TaskType (immediate|stage|scheduled) at the authoring/contract
// layer so consumers (packs / cloud / client) read the task kind from the
// contract, not only the truzhenos-internal domain.
func TestTaskCandidate_TaskKindAndRefsSerialize(t *testing.T) {
	tc := TaskCandidate{TaskName: "现场检查准备", TaskKind: "stage", ParentTaskRef: "formal_task://s1", PhaseRef: "phase://inspect"}
	b, err := json.Marshal(tc)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(b)
	for _, want := range []string{`"task_kind":"stage"`, `"parent_task_ref":"formal_task://s1"`, `"phase_ref":"phase://inspect"`} {
		if !strings.Contains(s, want) {
			t.Fatalf("expected %s in %s", want, s)
		}
	}
}

// Additive + backward compatible: an omitted task_kind (older packs) is absent
// via omitempty, so existing payloads/consumers are unaffected (default
// immediate is applied by the consumer).
func TestTaskCandidate_AdditiveFieldsOmitEmpty(t *testing.T) {
	tc := TaskCandidate{TaskName: "初筛"}
	b, _ := json.Marshal(tc)
	s := string(b)
	for _, forbidden := range []string{"task_kind", "parent_task_ref", "phase_ref"} {
		if strings.Contains(s, forbidden) {
			t.Fatalf("empty %s must be omitted (omitempty), got %s", forbidden, s)
		}
	}
}

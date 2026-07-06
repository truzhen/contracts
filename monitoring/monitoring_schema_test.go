package monitoring

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	contracts "github.com/truzhen/contracts"
)

func TestMonitoringEventErrorCodeJSONShape(t *testing.T) {
	event := MonitoringEvent{
		EventID:    "event://system-monitoring/run-1/evt-1",
		RunID:      "run-1",
		Sequence:   1,
		SourceKind: "monitor.http_fault",
		Severity:   SeverityError,
		Status:     "failed",
		ErrorCode:  "TZ-OS-HTTP-002",
		CreatedAt:  time.Date(2026, 7, 6, 12, 0, 0, 0, time.UTC),
	}
	raw, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("marshal MonitoringEvent: %v", err)
	}
	if !strings.Contains(string(raw), `"error_code":"TZ-OS-HTTP-002"`) {
		t.Fatalf("MonitoringEvent must expose additive error_code field: %s", raw)
	}

	event.ErrorCode = ""
	raw, err = json.Marshal(event)
	if err != nil {
		t.Fatalf("marshal MonitoringEvent without error code: %v", err)
	}
	if strings.Contains(string(raw), "error_code") {
		t.Fatalf("empty error_code must be omitted for older event compatibility: %s", raw)
	}
}

func TestMonitoringSchemasEmbeddedAndDeclareErrorCodePattern(t *testing.T) {
	for name, raw := range map[string][]byte{
		"monitoring-event": contracts.MonitoringEventSchemaJSON,
		"fault-incident":   contracts.FaultIncidentSchemaJSON,
	} {
		var doc map[string]any
		if err := json.Unmarshal(raw, &doc); err != nil {
			t.Fatalf("%s schema must be valid JSON: %v", name, err)
		}
		if !strings.Contains(string(raw), `"error_code"`) {
			t.Fatalf("%s schema must mention error_code: %s", name, raw)
		}
		if !strings.Contains(string(raw), `^TZ-[A-Z]{2,8}-[A-Z0-9]{2,10}-\\d{3}$`) {
			t.Fatalf("%s schema must constrain stable error_code pattern: %s", name, raw)
		}
	}
}

package events

import "github.com/truzhen/contracts/spines"

type ModuleEvent struct {
	EventID        string      `json:"event_id"`
	TransactionRef string      `json:"transaction_ref,omitempty"`
	SourceEventID  string      `json:"source_event_id,omitempty"`
	EventType      string      `json:"event_type"`
	Payload        interface{} `json:"payload"`
}

type IntentEvent = spines.IntentEvent

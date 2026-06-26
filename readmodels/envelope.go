package readmodels

// ReadModelEnvelope 用于前端展示，根据纪律，它绝不可作为真相源。
type ReadModelEnvelope struct {
	ViewID string      `json:"view_id"`
	Data   interface{} `json:"data"`
}

package events

type InconsistentEvent struct {
	ID   int64
	Type string
}

const (
	InconsistentEventTypeTargetMissing = "target_missing"
	InconsistentEventTypeNEQ           = "neq"
	InconsistentEventTypeBaseMissing   = "base_missing"
)

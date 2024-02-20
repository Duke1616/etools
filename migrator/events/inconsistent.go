package events

type InconsistentEvent struct {
	ID        int64
	Type      string
	Direction string
}

const (
	// InconsistentEventTypeTargetMissing target数据未命中
	InconsistentEventTypeTargetMissing = "target_missing"
	// InconsistentEventTypeNEQ 数据不相等
	InconsistentEventTypeNEQ = "neq"
	// InconsistentEventTypeBaseMissing base数据未命中
	InconsistentEventTypeBaseMissing = "base_missing"
)

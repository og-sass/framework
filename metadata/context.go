package metadata

type ContextTimeout int

const (
	TimeoutShort  ContextTimeout = 5
	TimeoutNormal ContextTimeout = 10
	TimeoutLong   ContextTimeout = 30
)

package constants

type ContextKeys string

const (
	TraceID ContextKeys = "x-trace-id"

	SpanID ContextKeys = "x-span-id"
)

func (ck ContextKeys) String() string {
	return string(ck)
}

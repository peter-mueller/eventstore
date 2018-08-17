package eventstore

type (
	Event struct {
		Data  []byte
		Index uint
	}

	EventStore interface {
		Push(data []byte) error
		StreamAll() (<-chan Event, error)
	}
)

func (e Event) IsBefore(other Event) bool {
	return e.Index < other.Index
}

func (e Event) IsAfter(other Event) bool {
	return e.Index > other.Index
}

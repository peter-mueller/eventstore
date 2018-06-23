package eventstore

type (
	Event struct {
		Data []byte
		Index uint
	}

	EventStore interface {
		Push(data []byte) error
		StreamAll() (<-chan Event, error)
		StreamFromIndex(index uint) (<-chan Event, error)
	}

	Information struct {
		StartIndex uint
		Name       string
	}
)

func (i Information) HasIndex(index uint) bool {
	return i.StartIndex <= index
}

func (e Event) IsBefore(other Event) bool {
	return e.Index < other.Index
}

func (e Event) IsAfter(other Event) bool {
	return e.Index > other.Index
}
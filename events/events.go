package events

type Event interface {
	IsEvent()
}

type EventTest1 struct {
	I int
}

func (*EventTest1) IsEvent() {}

type EventTest2 struct {
	J int
}

func (*EventTest2) IsEvent() {}

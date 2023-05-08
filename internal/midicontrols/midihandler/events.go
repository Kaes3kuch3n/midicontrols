package midihandler

type EventType int

const (
	Note EventType = iota
	CC
	Prog
)

type Event struct {
	Type  EventType `json:"type"`
	Value uint8     `json:"value"`
}

package domain

const (
	TimeCreatedEvent = "times.TimeCreated"
)

type TimeCreated struct {
	Time string
}

// Key implements registry.Registerable
func (TimeCreated) Key() string { return TimeCreatedEvent }

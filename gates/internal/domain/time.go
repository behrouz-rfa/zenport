package domain

type Time struct {
	Time string
}

func NewTime(time string) *Time {
	return &Time{Time: time}
}

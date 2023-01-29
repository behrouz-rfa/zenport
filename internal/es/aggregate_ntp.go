package es

import (
	"context"

	"zenport/internal/ddd"
)

type EventSourcedAggregate interface {
	ddd.IDer
	AggregateName() string
	ddd.Eventer
	Versioner
	EventApplier
	EventCommitter
}

type AggregateNtpMiddleware func(store AggregateNtp) AggregateNtp

type AggregateNtp interface {
	Save(ctx context.Context, aggregate EventSourcedAggregate) error
}

func AggregateNtpWithMiddleware(store AggregateNtp, mws ...AggregateNtpMiddleware) AggregateNtp {
	//	var s AggregateNtp
	s := store
	// middleware are applied in reverse; this makes the first middleware
	// in the slice the outermost i.e. first to enter, last to exit
	// given: store, A, B, C
	// result: A(B(C(store)))
	for i := len(mws) - 1; i >= 0; i-- {
		s = mws[i](s)
	}
	return s
}

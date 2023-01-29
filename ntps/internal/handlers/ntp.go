package handlers

import (
	"zenport/internal/ddd"
	"zenport/ntps/internal/domain"
)

func RegisterNtpHandlers(ntpHandlers ddd.EventHandler[ddd.AggregateEvent], domainSubscriber ddd.EventSubscriber[ddd.AggregateEvent]) {
	domainSubscriber.Subscribe(ntpHandlers,
		domain.TimeCreatedEvent,
	)
}

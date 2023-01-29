package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/stackus/errors"

	"zenport/internal/es"
	"zenport/internal/registry"
)

type EventStore struct {
	tableName string
	db        *sql.DB
	registry  registry.Registry
}

var _ es.AggregateNtp = (*EventStore)(nil)

func NewEventStore(tableName string, db *sql.DB, registry registry.Registry) EventStore {
	return EventStore{
		tableName: tableName,
		db:        db,
		registry:  registry,
	}
}

func (s EventStore) Save(ctx context.Context, aggregate es.EventSourcedAggregate) (err error) {
	const query = `INSERT INTO %s (stream_id, stream_name, stream_version, event_id, event_name, event_data, occurred_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	var tx *sql.Tx
	tx, err = s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		p := recover()
		switch {
		case p != nil:
			_ = tx.Rollback()
			panic(p)
		case err != nil:
			rErr := tx.Rollback()
			if rErr != nil {
				err = errors.Wrap(err, rErr.Error())
			}
		default:
			err = tx.Commit()
		}
	}()

	aggregateID := aggregate.ID()
	aggregateName := aggregate.AggregateName()

	for _, event := range aggregate.Events() {
		var payloadData []byte

		payloadData, err = s.registry.Serialize(event.EventName(), event.Payload())
		if err != nil {
			return err
		}
		if _, err = tx.ExecContext(
			ctx, s.table(query),
			aggregateID, aggregateName, event.AggregateVersion(), event.ID(), event.EventName(), payloadData, event.OccurredAt(),
		); err != nil {
			return err
		}
	}

	return nil
}

func (s EventStore) table(query string) string {
	return fmt.Sprintf(query, s.tableName)
}

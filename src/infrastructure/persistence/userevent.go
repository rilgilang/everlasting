package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	errDomain "everlasting/src/domain/error"
	"everlasting/src/domain/event"
	"everlasting/src/domain/sharedkernel/identity"
	"everlasting/src/domain/sharedkernel/unitofwork"
	"everlasting/src/infrastructure/pkg/logger"
	"time"

	"github.com/jmoiron/sqlx"
)

type UserEventPersistence struct {
	db     *sqlx.DB
	logger *logger.AppLogger
}

func NewUserEventPersistence(db *sqlx.DB, logger *logger.AppLogger) (result *UserEventPersistence) {
	return &UserEventPersistence{db, logger}
}

// GetOneByID retrieves an event by its ID
func (p *UserEventPersistence) GetOneByUserID(ctx context.Context, id identity.ID) (*event.UserEvent, error) {
	query := `SELECT 
		id, user_id, event_id,created_at, updated_at 
	FROM 
		user_events 
	WHERE user_id = $1`

	stmt, err := p.generateStatement(ctx, query)
	if err != nil {
		p.logger.Error(ctx, "postgres:user_events_get_one_by_id", err.Error())
		return nil, err
	}
	defer stmt.Close()

	var userEvent event.UserEvent
	var eventIdsJSON []byte
	err = stmt.QueryRowContext(ctx, id).Scan(
		&userEvent.ID, &userEvent.UserId, &eventIdsJSON, &userEvent.CreatedAt, &userEvent.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errDomain.ErrEventNotFound
		}
		p.logger.Error(ctx, "postgres:user_events_get_one_by_id", err.Error())
		return nil, err
	}

	if len(eventIdsJSON) > 0 {
		if err := json.Unmarshal(eventIdsJSON, &userEvent.EventId); err != nil {
			p.logger.Error(ctx, "postgres:user_events_json_unmarshal", err.Error())
			return nil, err
		}
	}

	return &userEvent, nil
}

func (p *UserEventPersistence) UpdateByUserId(ctx context.Context, eventIds []string, userId string) error {
	now := time.Now().UTC()
	query := `UPDATE
		user_events
	SET
		event_id=$2
		updated_at=$3
	WHERE
		user_id=$1`

	stmt, err := p.generateStatement(ctx, query)
	if err != nil {
		p.logger.Error(ctx, "postgres:user_events_update", err.Error())
		return err
	}
	defer stmt.Close()

	vals := []interface{}{
		userId,
		eventIds,
		now,
	}

	_, err = stmt.ExecContext(ctx, vals...)
	if err != nil {
		if err == sql.ErrNoRows {
			return errDomain.ErrEventNotFound
		}
		p.logger.Error(ctx, "postgres:user_events_update", err.Error())
		return err
	}

	return nil
}

func (p *UserEventPersistence) CreateUserEvents(ctx context.Context, eventIds []string, userId string) error {
	now := time.Now().UTC()
	query := `INSERT INTO
		user_events ( user_id, event_id, created_at, updated_at ) VALUES ($1, $2, $3, $4)`

	stmt, err := p.generateStatement(ctx, query)
	if err != nil {
		p.logger.Error(ctx, "postgres:user_events_update", err.Error())
		return err
	}
	defer stmt.Close()

	vals := []interface{}{
		userId,
		eventIds,
		now,
		now,
	}

	_, err = stmt.ExecContext(ctx, vals...)
	if err != nil {
		if err == sql.ErrNoRows {
			return errDomain.ErrEventNotFound
		}
		p.logger.Error(ctx, "postgres:user_events_create", err.Error())
		return err
	}

	return nil
}

// generateStatement prepares a statement, handling transaction context if present
func (p *UserEventPersistence) generateStatement(ctx context.Context, sql string) (*sqlx.Stmt, error) {
	stmt, err := p.db.Preparex(sql)

	// Check if operation is in transaction context
	tx, ok := ctx.Value(unitofwork.TransactionContextKey).(*sqlx.Tx)
	if ok && tx != nil {
		stmt, err = tx.Preparex(sql)
	}

	return stmt, err
}
